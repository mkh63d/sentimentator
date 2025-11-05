using Microsoft.AspNetCore.Authentication.BearerToken;
using Microsoft.AspNetCore.Authorization;
using Microsoft.OpenApi;
using Scalar.AspNetCore;
using Sentimentator.Backend;
using Sentimentator.Backend.Api.EndpointGroups;
using Sentimentator.Backend.Api.Helpers;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;

var builder = WebApplication.CreateBuilder(args);
builder.Services.AddDbContext<SentimentatorContext>();
builder.Services
    .AddIdentityApiEndpoints<User>()
    .AddEntityFrameworkStores<SentimentatorContext>();
builder.Services.AddScoped<AnalysisService>();

builder.Services.AddAuthentication(options =>
{
    options.DefaultAuthenticateScheme = "Bearer";
    options.DefaultChallengeScheme = "Bearer";
});

builder.Services.AddAuthorizationBuilder();
builder.Services.AddOpenApi(options =>
{
    options
        .AddDocumentTransformer((document, _, _) =>
        {
            document.Info = new OpenApiInfo
            {
                Title = "Sentimentator backend API",
                Version = "v1",
                Description = "API for retrieving sentiment analysis data",
                Contact = new OpenApiContact
                {
                    Name = "API Support",
                    Email = "api@example.com",
                    Url = new Uri("https://api.example.com/support")
                }
            };
            return Task.CompletedTask;
        })
        .AddBearerSecurityScheme();
});

var app = builder.Build();
app.MapOpenApi().CacheOutput();
app.MapScalarApiReference(options =>
{
    options.Authentication = new ScalarAuthenticationOptions
    {
        PreferredSecuritySchemes = ["Bearer"]
    };
    options.Theme = ScalarTheme.BluePlanet;
});

app.UseAuthentication();
app.UseAuthorization();

app.MapGet("/", () => Results.Redirect("/scalar/v1"))
    .ExcludeFromDescription();

app.MapGet("/hello", () => Results.Ok("Hello World!"))
    .WithTags("Hello")
    .RequireAuthorization();

app.MapGroup("/analyses")
    .AnalysisPostGroup()
    .AnalysisGetGroup()
    .AnalysisPutGroup()
    .AnalysisDeleteGroup()
    .WithTags("Analysis")
    .RequireAuthorization();

app.MapGroup("/account")
    .MapIdentityApi<User>()
    .WithTags("Identity");

app.Run();