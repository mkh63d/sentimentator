using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Sentimentator.Backend.Service.Analysis;
using Sentimentator.Backend.Service.RabbitMq.Consumer.Model;
using static System.Text.Json.JsonSerializer;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class MetadataRabbitConsumer(
    ILogger<MetadataRabbitConsumer> logger,
    IServiceScopeFactory serviceScopeFactory,
    IRabbitConsumer rabbitConsumer) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        const string queueName = "video.metadata";
        const string routingKey = "video.metadata.completed";

        await rabbitConsumer.Receive(queueName, routingKey, HandleMessage);
    }

    private async Task<bool> HandleMessage(string message)
    {
        MetadataComplete? data;
        try
        {
            data = Deserialize<MetadataComplete>(message);
        }
        catch (JsonException e)
        {
            logger.LogError("Exception occurred while deserializing message; LineNumber: {Line}; Path: {Path}",
                e.LineNumber, e.Path);
            return false;
        }

        if (data is null)
        {
            logger.LogError("Failed to deserialize message");
            return false;
        }

        using var scope = serviceScopeFactory.CreateScope();
        var service = scope.ServiceProvider.GetRequiredService<AnalysisService>();

        return await service.Update(data.AnalysisId, data.Platform, data.Title, data.ThumbnailUrl, data.Duration, DateTime.UtcNow);
    }
}