using Microsoft.Extensions.Hosting;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class HostedRabbitConsumer(IRabbitConsumer rabbitConsumer) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        await rabbitConsumer.Receive("video.analysis.complete");
    }
}