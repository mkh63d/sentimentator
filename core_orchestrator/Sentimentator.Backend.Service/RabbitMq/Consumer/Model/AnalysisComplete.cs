using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

public class AnalysisComplete
{
    [JsonPropertyName("video_id")] public Guid VideoId { get; set; }
    [JsonPropertyName("transcript")] public required string Transcript { get; set; }
    [JsonPropertyName("data")] public required Dictionary<string, string> Data { get; set; }
}