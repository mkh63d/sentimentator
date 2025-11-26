using System.Text.Json.Serialization;

namespace Sentimentator.Backend.Model;

public class Analysis
{
    [JsonPropertyName("id")] public Guid Id { get; init; }
    [JsonPropertyName("url")] public required Uri Url { get; init; }
    [JsonPropertyName("created_at")] public DateTime CreatedAt { get; init; }
    [JsonPropertyName("updated_at")] public DateTime UpdatedAt { get; set; }
    [JsonPropertyName("transcription")] public string? Transcription { get; set; }
    [JsonPropertyName("summary")] public string? Summary { get; set; }
    [JsonIgnore] [JsonPropertyName("users")] public List<User> Users { get; } = [];
}