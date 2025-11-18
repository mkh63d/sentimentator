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

public class RabbitConsumer(ILogger<RabbitConsumer> logger, IServiceScopeFactory serviceScopeFactory)
    : RabbitBase(logger), IRabbitConsumer
{
    public async Task Receive(string queue, string routingKey)
    {
        await Channel.ExchangeDeclareAsync(ExchangeName, ExchangeType.Topic, durable: true, autoDelete: false);
        await Channel.QueueDeclareAsync(queue, durable: true, exclusive: false, autoDelete: false);
        await Channel.QueueBindAsync(queue, ExchangeName, routingKey);
        
        var consumer = new AsyncEventingBasicConsumer(Channel);
        consumer.ReceivedAsync += async (_, ea) =>
        {
            var body = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(body);
            Logger.LogInformation("Received message from queue {Queue} with routing key {RoutingKey} at {Time} with message \"{Message}\"", queue, routingKey,
                DateTime.UtcNow.ToLongTimeString(), message);

            var result = await HandleMessage(message);
            if (result)
            {
                Logger.LogInformation("Message processed successfully, acknowledging...");
                await Channel.BasicAckAsync(ea.DeliveryTag, false);
            }
            else
            {
                Logger.LogInformation("Message processing failed, not acknowledging.");
                await Channel.BasicNackAsync(ea.DeliveryTag, false, false);
            }
        };

        await Channel.BasicConsumeAsync(queue, false, consumer);
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
            Logger.LogError("Exception occurred while deserializing message; LineNumber: {Line}; Path: {Path}", e.LineNumber, e.Path);
            return false; 
        }
        if (data is null)
        {
            Logger.LogError("Failed to deserialize message");
            return false;
        }
        
        using var scope = serviceScopeFactory.CreateScope();
        var service = scope.ServiceProvider.GetRequiredService<AnalysisService>();
        
        return await service.Update(data.AnalysisId, data.Transcript, Serialize(data.Sentiment), DateTime.UtcNow);
    }
}