using Microsoft.Extensions.Hosting;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class HostedRabbitConsumer(IRabbitConsumer rabbitConsumer) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        const string queueName = "video.analysis";
        const string routingKey = "video.analysis.completed";

        await rabbitConsumer.Receive(queueName, routingKey);
    }
}