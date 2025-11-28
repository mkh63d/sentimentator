using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

public class MetadataComplete
{
    [JsonPropertyName("analysis_id")] public Guid AnalysisId { get; set; }
    [JsonPropertyName("platform")] public required string Platform { get; set; }
    [JsonPropertyName("title")] public required string Title { get; set; }
    [JsonPropertyName("thumbnail_url")] public required string ThumbnailUrl { get; set; }
}