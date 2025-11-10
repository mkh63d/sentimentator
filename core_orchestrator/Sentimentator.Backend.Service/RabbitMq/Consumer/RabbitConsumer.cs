using System.Text;
using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;
using Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

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
        _channel.DisposeAsync();
        _connection.DisposeAsync();
    }

    public async Task Receive(string queue)
    {
        await _channel.QueueDeclareAsync(queue, false, false, false);

        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, ea) =>
        {
            var body = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(body);
            _logger.LogInformation("Received message from queue {Queue} at {Time} with message \"{Message}\"", queue,
                DateTime.UtcNow.ToLongTimeString(), message);

            await HandleMessage(message);
            await _channel.BasicAckAsync(ea.DeliveryTag, false);
        };

        await _channel.BasicConsumeAsync(queue, false, consumer);
    }

    private async Task HandleMessage(string message)
    {
        var data = JsonSerializer.Deserialize<AnalysisComplete>(message);
        if (data is null)
        {
            _logger.LogError("Failed to deserialize message: {Message}", message);
            return;
        }

        using var scope = _serviceScopeFactory.CreateScope();
        var service = scope.ServiceProvider.GetRequiredService<AnalysisService>();

        await service.Update(data.VideoId, data.Transcript, JsonSerializer.Serialize(data.Data));
    }
}