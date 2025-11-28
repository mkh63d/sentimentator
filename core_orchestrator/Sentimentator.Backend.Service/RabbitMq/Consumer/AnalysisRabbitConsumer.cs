using System.Text.Json;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.Extensions.Hosting;
using Microsoft.Extensions.Logging;
using Sentimentator.Backend.Service.Analysis;
using Sentimentator.Backend.Service.RabbitMq.Consumer.Model;
using static System.Text.Json.JsonSerializer;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class AnalysisRabbitConsumer(ILogger<AnalysisRabbitConsumer> logger, IServiceScopeFactory serviceScopeFactory, IRabbitConsumer rabbitConsumer) : BackgroundService
{
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        const string queueName = "video.analysis";
        const string routingKey = "video.analysis.completed";

        await rabbitConsumer.Receive(queueName, routingKey, HandleMessage);
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
            logger.LogError("Exception occurred while deserializing message; LineNumber: {Line}; Path: {Path}", e.LineNumber, e.Path);
            return false; 
        }
        if (data is null)
        {
            logger.LogError("Failed to deserialize message");
            return false;
        }
        
        using var scope = serviceScopeFactory.CreateScope();
        var service = scope.ServiceProvider.GetRequiredService<AnalysisService>();
        
        return await service.Update(data.AnalysisId, data.Transcript, Serialize(data.Sentiment), DateTime.UtcNow);
    }
}