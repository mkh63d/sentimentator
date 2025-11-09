using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend;

public class SentimentatorContext
    : IdentityDbContext<User, IdentityRole<Guid>, Guid>
{
    public DbSet<Analysis> Analyses { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
    }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        optionsBuilder.UseNpgsql("Host=localhost;Database=orchestrator");
        // TODO: change to environment variable
    }
}