using System.Text;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq;

public class RabbitPublisher(ILogger<RabbitPublisher> logger)
{
    readonly ConnectionFactory _factory = new() { Uri = new Uri("amqp://guest:guest@localhost:5672/") }; // TODO: change to environment variable
    
    public async Task Send(string queue, string message)
    {
        await using var connection = await _factory.CreateConnectionAsync();
        await using var channel = await connection.CreateChannelAsync();
        
        await channel.QueueDeclareAsync(queue: queue, durable: false, exclusive: false, autoDelete: false,
            arguments: null);
        
        var body = Encoding.UTF8.GetBytes(message);
        
        await channel.BasicPublishAsync(exchange: string.Empty, routingKey: queue, body: body);
        logger.LogInformation("Message published to queue {} at {}", queue, DateTime.UtcNow.ToLongTimeString());
    }
}