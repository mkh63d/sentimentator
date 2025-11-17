using System.Text;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq.Publisher;

public class RabbitPublisher(ILogger<RabbitPublisher> logger) : RabbitBase(logger), IRabbitPublisher
{
    public async Task Send(string queue, string routingKey, string message)
    {
        await Channel.ExchangeDeclareAsync(ExchangeName, ExchangeType.Topic, durable: true, autoDelete: false);
        await Channel.QueueDeclareAsync(queue, durable: true, exclusive: false, autoDelete: false);
        await Channel.QueueBindAsync(queue, ExchangeName, routingKey);

        var body = Encoding.UTF8.GetBytes(message);

        await Channel.BasicPublishAsync(ExchangeName, routingKey, body);
        logger.LogInformation("Message published to exchange {Exchange} with routing key {RoutingKey} at {Time} with message \"{Message}\"",
            ExchangeName, routingKey, DateTime.UtcNow.ToLongTimeString(), message);
    }
}
