using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer.Model;

public class DownloadRequest
{
    [JsonPropertyName("video_id")] public Guid VideoId { get; set; }
    [JsonPropertyName("url")] public required Uri Url { get; set; }
}