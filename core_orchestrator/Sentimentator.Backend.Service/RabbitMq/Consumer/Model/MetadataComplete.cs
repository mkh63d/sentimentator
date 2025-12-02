using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

public class MetadataComplete
{
    [JsonPropertyName("analysis_id")] public Guid AnalysisId { get; set; }
    [JsonPropertyName("platform")] public string? Platform { get; set; }
    [JsonPropertyName("title")] public string? Title { get; set; }
    [JsonPropertyName("thumbnail_url")] public string? ThumbnailUrl { get; set; }
    [JsonPropertyName("duration")] public int? Duration { get; set; }
}