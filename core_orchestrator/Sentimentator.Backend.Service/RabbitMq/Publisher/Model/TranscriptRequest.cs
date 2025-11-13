using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Publisher.Model;

public class TranscriptRequest
{
    [JsonPropertyName("analysis_id")] public Guid AnalysisId { get; set; }
    [JsonPropertyName("url")] public required Uri Url { get; set; }
    [JsonPropertyName("timestamp")] public required DateTime Timestamp { get; set; }
}