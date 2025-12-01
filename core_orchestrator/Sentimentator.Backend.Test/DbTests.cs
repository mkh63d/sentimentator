using Microsoft.EntityFrameworkCore;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Test;

public class DbTests : IDisposable
{
    private readonly SentimentatorContext _context;

    public DbTests()
    {
        Environment.SetEnvironmentVariable("ORCHESTRATOR_POSTGRES_CONN_STRING",
            "Server=localhost;Port=5432;Database=orchestrator;User Id=postgres;Password=postgres;");
        _context = new SentimentatorContext();
    }

    public void Dispose()
    {
        _context.ChangeTracker.Entries().ToList()
            .ForEach(e => e.State = EntityState.Detached);
        _context.Analyses.RemoveRange(_context.Analyses.Where(x => x.Url == new Uri("http://example.com")));
        _context.SaveChanges();
        _context.Dispose();
    }

    [Fact]
    public void ShouldCompleteWhenAddingCorrectAnalysis()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = new Uri("http://example.com"),
            CreatedAt = DateTime.UtcNow
        };
        _context.Analyses.Add(analysis);
        var exception = Record.Exception(() => _context.SaveChanges());
        Assert.Null(exception);
    }

    [Fact]
    public void ShouldThrowWhenAddingAnalysisWithInvalidUrl()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = null!,
            CreatedAt = DateTime.UtcNow
        };
        _context.Analyses.Add(analysis);

        void Act()
        {
            _context.SaveChanges();
        }

        Assert.Throws<DbUpdateException>(Act);
    }

    [Fact]
    public void ShouldCompleteWhenRetrievingExistingAnalysis()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = new Uri("http://example.com"),
            CreatedAt = DateTime.UtcNow
        };
        _context.Analyses.Add(analysis);
        _context.SaveChanges();

        var retrievedAnalysis = _context.Set<Analysis>().Find(analysis.Id);
        Assert.NotNull(retrievedAnalysis);
        Assert.Equal(analysis.Id, retrievedAnalysis.Id);
    }

    [Fact]
    public void ShouldReturnNullWhenRetrievingNonExistingAnalysis()
    {
        var retrievedAnalysis = _context.Set<Analysis>().Find(Guid.NewGuid());
        Assert.Null(retrievedAnalysis);
    }

    [Fact]
    public void ShouldCompleteWhenUpdatingExistingAnalysisWithCorrectInfo()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = new Uri("http://example.com"),
            CreatedAt = DateTime.UtcNow
        };
        _context.Analyses.Add(analysis);
        _context.SaveChanges();
        analysis.Title = "Updated Title";
        _context.Analyses.Update(analysis);
        var exception = Record.Exception(() => _context.SaveChanges());
        Assert.Null(exception);
    }

    [Fact]
    public void ShouldCompleteWhenDeletingExistingAnalysis()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = new Uri("http://example.com"),
            CreatedAt = DateTime.UtcNow
        };
        _context.Analyses.Add(analysis);
        _context.SaveChanges();
        _context.Analyses.Remove(analysis);
        var exception = Record.Exception(() => _context.SaveChanges());
        Assert.Null(exception);
    }

    [Fact]
    public void ShouldThrowWhenDeletingNonExistingAnalysis()
    {
        var analysis = new Analysis
        {
            Id = Guid.NewGuid(),
            Url = new Uri("http://example.com"),
            CreatedAt = DateTime.UtcNow
        };
        Assert.Throws<DbUpdateConcurrencyException>(() =>
        {
            _context.Analyses.Remove(analysis);
            _context.SaveChanges();
        });
    }
}