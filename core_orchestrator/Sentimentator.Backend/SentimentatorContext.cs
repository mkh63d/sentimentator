using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend;

public class SentimentatorContext
    : IdentityDbContext<User, IdentityRole<Guid>, Guid>
{
    private const string PostgresConnStringEnvVar = "ORCHESTRATOR_POSTGRES_CONN_STRING";
    private const string PostgresDefaultConnString = "Host=localhost;Database=orchestrator";

    private static readonly string RabbitMqUri = Environment.GetEnvironmentVariable(PostgresConnStringEnvVar)
                                                 ?? PostgresDefaultConnString;

    public DbSet<Analysis> Analyses { get; set; }

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        base.OnModelCreating(modelBuilder);
    }

    protected override void OnConfiguring(DbContextOptionsBuilder optionsBuilder)
    {
        optionsBuilder.UseNpgsql(RabbitMqUri);
    }
}