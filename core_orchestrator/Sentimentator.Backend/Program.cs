using Sentimentator.Backend;
using Sentimentator.Backend.Model;

var context = new SentimentatorContext();

context.Database.EnsureCreated();

context.Analyses.Add(new Analysis
{
    Id = Guid.NewGuid(),
    CreatedAt = DateTime.UtcNow,
    Url = new Uri("https://www.example.com/video1")
});

context.SaveChanges();