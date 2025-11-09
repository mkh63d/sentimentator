using Sentimentator.Backend;
using Sentimentator.Backend.Model;

var context = new SentimentatorContext();

context.Database.EnsureCreated();

context.Analyses.Add(new Analysis
{
    Id = Guid.NewGuid(),
    Summary = "test summary",
    Transcription = "test transcription",
    Url = "test url"
});

context.SaveChanges();