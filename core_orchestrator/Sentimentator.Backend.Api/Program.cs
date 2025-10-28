using System.Web;
using Microsoft.AspNetCore.Identity;
using Scalar.AspNetCore;
using Sentimentator.Backend;
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
        document.Info = new()
        {
            Title = "Sentimentator backup API",
            Version = "v1",
            Description = """
                          API for retrieving sentiment analisis data
                          """,
            Contact = new()
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

app.MapGet("/analyses/{id:guid}", (Guid id) =>
{
    using var scope = app.Services.CreateScope();
    var service = scope.ServiceProvider.GetService<AnalysisService>();
    try
    {
        var analysis = service!.Get(id);
        return Results.Ok(analysis);
    }
    catch
    {
        return Results.NotFound();
    }
});

app.MapGet("/analyses/{url}", (string url) =>
{
    url = HttpUtility.UrlDecode(url);
    using var scope = app.Services.CreateScope();
    var service = scope.ServiceProvider.GetService<AnalysisService>();
    try
    {
        var analysis = service!.Get(url);
        return Results.Ok(analysis);
    }
    catch
    {
        return Results.NotFound();
    }
});

app.Run();