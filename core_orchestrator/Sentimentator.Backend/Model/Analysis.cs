namespace Sentimentator.Backend.Model;

public class Analysis
{
    public Guid Id { get; init; }
    public required Uri Url { get; init; }
    public string? Transcription { get; set; }
    public string? Summary { get; set; }
}