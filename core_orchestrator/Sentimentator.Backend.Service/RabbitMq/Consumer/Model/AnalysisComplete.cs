using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

public class AnalysisComplete
{
    [JsonPropertyName("analysis_id")] public Guid AnalysisId { get; set; }
    [JsonPropertyName("transcript")] public required string Transcript { get; set; }
    [JsonPropertyName("sentiment")] public required Dictionary<string, string> Sentiment { get; set; }
}