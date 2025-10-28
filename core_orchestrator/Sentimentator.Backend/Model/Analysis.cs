namespace Sentimentator.Backend.Model;

public class Analysis
{
    public Guid Id { get; set; }
    public required string Url { get; set; }
    public string? Transcription { get; set; }
    public string? Summary  { get; set; }
}