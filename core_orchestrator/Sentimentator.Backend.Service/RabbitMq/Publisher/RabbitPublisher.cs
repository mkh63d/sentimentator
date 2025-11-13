using System.Text;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq.Publisher;

public class RabbitPublisher(ILogger<RabbitPublisher> logger) : RabbitBase
{
    private const string ExchangeName = "sentimentator.topic";

    public async Task Send(string routingKey, string message)
    {
        await using var connection = await Factory.CreateConnectionAsync();
        await using var channel = await connection.CreateChannelAsync();

        await channel.ExchangeDeclareAsync(ExchangeName, "topic", durable: true, autoDelete: false);

        var body = Encoding.UTF8.GetBytes(message);

        await channel.BasicPublishAsync(ExchangeName, routingKey, body);
        logger.LogInformation("Message published to exchange {Exchange} with routing key {RoutingKey} at {Time} with message \"{Message}\"",
            ExchangeName, routingKey, DateTime.UtcNow.ToLongTimeString(), message);
    }
}
