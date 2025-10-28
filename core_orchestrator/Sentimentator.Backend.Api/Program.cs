using Microsoft.AspNetCore.Identity;
using Sentimentator.Backend;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDbContext<SentimentatorContext>();
builder.Services.AddIdentity<User, IdentityRole<Guid>>().AddEntityFrameworkStores<SentimentatorContext>();
builder.Services.AddScoped<AnalysisService>();
var app = builder.Build();

app.MapGet("/", () => "Hello World!");

app.MapGet("/analyses/{Guid: id}", (Guid id) =>
{
    var service = app.Services.GetService<AnalysisService>();
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

app.Run();