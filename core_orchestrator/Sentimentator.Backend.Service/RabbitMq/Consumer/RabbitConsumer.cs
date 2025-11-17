using System.Text;
using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using Sentimentator.Backend.Service.Analysis;
using Sentimentator.Backend.Service.RabbitMq.Consumer.Model;
using static System.Text.Json.JsonSerializer;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class RabbitConsumer : RabbitBase, IRabbitConsumer, IDisposable
{
    private readonly IChannel _channel;
    private readonly IConnection _connection;
    private readonly ILogger<RabbitConsumer> _logger;
    private readonly IServiceScopeFactory _serviceScopeFactory;

    public RabbitConsumer(ILogger<RabbitConsumer> logger, IServiceScopeFactory serviceScopeFactory)
    {
        _logger = logger;
        _connection = Factory.CreateConnectionAsync().Result;
        _channel = _connection.CreateChannelAsync().Result;
        _serviceScopeFactory = serviceScopeFactory;
    }

    public void Dispose()
    {
        _channel.CloseAsync().Wait();
        _connection.CloseAsync().Wait();
        _channel.DisposeAsync().AsTask().Wait();
        _connection.DisposeAsync().AsTask().Wait();
    }

    public async Task Receive(string queue, string routingKey)
    {
        await _channel.ExchangeDeclareAsync(ExchangeName, ExchangeType.Topic, durable: true, autoDelete: false);
        await _channel.QueueDeclareAsync(queue, durable: true, autoDelete: false);
        await _channel.QueueBindAsync(queue, ExchangeName, routingKey);
        
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (_, ea) =>
        {
            var body = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(body);
            _logger.LogInformation("Received message from queue {Queue} with routing key {RoutingKey} at {Time} with message \"{Message}\"", queue, routingKey,
                DateTime.UtcNow.ToLongTimeString(), message);

            var result = await HandleMessage(message);
            if (result)
            {
                _logger.LogInformation("Message processed successfully, acknowledging...");
                await _channel.BasicAckAsync(ea.DeliveryTag, false);
            }
            else
            {
                _logger.LogInformation("Message processing failed, not acknowledging.");
                await _channel.BasicNackAsync(ea.DeliveryTag, false, false);
            }
        };

        await _channel.BasicConsumeAsync(queue, false, consumer);
    }

    private async Task<bool> HandleMessage(string message)
    {
        AnalysisComplete? data;
        try
        {
            data = Deserialize<AnalysisComplete>(message);
        }
        catch (JsonException e)
        {
            _logger.LogError("Exception occurred while deserializing message; LineNumber: {Line}; Path: {Path}", e.LineNumber, e.Path);
            return false; 
        }
        if (data is null)
        {
            _logger.LogError("Failed to deserialize message");
            return false;
        }
        
        using var scope = _serviceScopeFactory.CreateScope();
        var service = scope.ServiceProvider.GetRequiredService<AnalysisService>();
        
        return await service.Update(data.AnalysisId, data.Transcript, Serialize(data.Sentiment));
    }
}