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

        await channel.QueueDeclareAsync(queue, false, false, false,
            null);

        var body = Encoding.UTF8.GetBytes(message);

        await channel.BasicPublishAsync(string.Empty, queue, body);
        logger.LogInformation("Message published to queue {} at {} with message \"{}\"", queue,
            DateTime.UtcNow.ToLongTimeString(), message);
    }
}