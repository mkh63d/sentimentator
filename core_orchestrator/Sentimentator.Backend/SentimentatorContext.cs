using Microsoft.EntityFrameworkCore;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend;

public class SentimentatorContext : DbContext
{
    public DbSet<Analysis> Analyses { get; set; }
}