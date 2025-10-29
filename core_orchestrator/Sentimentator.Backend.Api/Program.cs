using System.Web;
using Microsoft.AspNetCore.Identity;
using Microsoft.OpenApi.Models;
using Scalar.AspNetCore;
using Sentimentator.Backend;
using Sentimentator.Backend.Api.EndpointGroups;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDbContext<SentimentatorContext>();
builder.Services.AddIdentity<User, IdentityRole<Guid>>().AddEntityFrameworkStores<SentimentatorContext>();
builder.Services.AddScoped<AnalysisService>();

builder.Services.AddOpenApi(options =>
{
    options.AddDocumentTransformer((document, _, _) =>
    {
        document.Info = new OpenApiInfo
        {
            Title = "Sentimentator backend API",
            Version = "v1",
            Description = """
                          API for retrieving sentiment analysis data
                          """,
            Contact = new OpenApiContact
            {
                Name = "API Support",
                Email = "api@example.com",
                Url = new Uri("https://api.example.com/support")
            }
        };
        return Task.CompletedTask;
    });
});

var app = builder.Build();
app.MapOpenApi().CacheOutput();
app.MapScalarApiReference();

app.MapGet("/", () => Results.Redirect("/scalar/v1"))
    .ExcludeFromDescription();

app.MapGroup("/analyses")
    .AnalysisPostGroup()
    .AnalysisGetGroup()
    .AnalysisPutGroup()
    .AnalysisDeleteGroup();

app.Run();