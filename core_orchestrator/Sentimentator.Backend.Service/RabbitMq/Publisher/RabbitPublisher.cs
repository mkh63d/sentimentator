using System.Text;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq.Publisher;

public class RabbitPublisher(ILogger<RabbitPublisher> logger) : RabbitBase
{
    public async Task Send(string queue, string message)
    {
        await using var connection = await Factory.CreateConnectionAsync();
        await using var channel = await connection.CreateChannelAsync();
        
        await channel.QueueDeclareAsync(queue: queue, durable: false, exclusive: false, autoDelete: false,
            arguments: null);
        
        var body = Encoding.UTF8.GetBytes(message);
        
        await channel.BasicPublishAsync(exchange: string.Empty, routingKey: queue, body: body);
        logger.LogInformation("Message published to queue {} at {} with message \"{}\"", queue, DateTime.UtcNow.ToLongTimeString(), message);
    }
}