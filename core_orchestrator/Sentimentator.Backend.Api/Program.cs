using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Identity.UI.Services;
using Scalar.AspNetCore;
using Sentimentator.Backend;
using Sentimentator.Backend.Api.EndpointGroups;
using Sentimentator.Backend.Api.Transformers;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;
using Sentimentator.Backend.Service.Analysis;
using Sentimentator.Backend.Service.Email;
using Sentimentator.Backend.Service.RabbitMq.Consumer;
using Sentimentator.Backend.Service.RabbitMq.Publisher;

var builder = WebApplication.CreateBuilder(args);

#region Logging

builder.Logging.ClearProviders();
builder.Logging.SetMinimumLevel(LogLevel.Debug);
builder.Logging.AddConsole();

#endregion

#region Identity

builder.Services.AddDbContext<SentimentatorContext>();
builder.Services
    .AddIdentityCore<User>(options =>
    {
        options.User.AllowedUserNameCharacters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-._@+";
        options.User.RequireUniqueEmail = true;

        options.Password.RequireDigit = true;
        options.Password.RequireNonAlphanumeric = true;
        options.Password.RequireUppercase = true;
        options.Password.RequireLowercase = true;

        options.Lockout.MaxFailedAccessAttempts = 5;
        options.Lockout.DefaultLockoutTimeSpan = TimeSpan.FromMinutes(5);

        options.SignIn.RequireConfirmedEmail = true;
    })
    .AddRoles<IdentityRole<Guid>>()
    .AddEntityFrameworkStores<SentimentatorContext>()
    .AddApiEndpoints();

#endregion

#region CORS

builder.Services.AddCors(options =>
{
    options.AddPolicy("AllowSpecificOrigin", policy =>
    {
        policy.WithOrigins("http://localhost:5173", "http://localhost:3000", "http://localhost:8000")
            .AllowCredentials()
            .AllowAnyHeader()
            .AllowAnyMethod()
            .AllowCredentials();
    });
});

#endregion

#region Auth

builder.Services
    .AddAuthentication()
    .AddBearerToken(IdentityConstants.BearerScheme, options =>
    {
        options.BearerTokenExpiration = TimeSpan.FromSeconds(3600);
        options.RefreshTokenExpiration = TimeSpan.FromSeconds(7200);
    });

builder.Services.AddAuthorization(options =>
{
    options.DefaultPolicy = new AuthorizationPolicyBuilder()
        .RequireAuthenticatedUser()
        .AddAuthenticationSchemes(IdentityConstants.BearerScheme)
        .Build();

    options.FallbackPolicy = options.DefaultPolicy;
});

#endregion

#region OpenAPI

builder.Services.AddOpenApi(options =>
{
    options
        .AddInfo()
        .AddTags()
        .AddBearerSecurityScheme();
});

#endregion

#region Services

builder.Services.AddTransient<IEmailSender, ResendSender>();
builder.Services.AddScoped<IRabbitPublisher, RabbitPublisher>();
builder.Services.AddScoped<AnalysisService>();
builder.Services.AddSingleton<IRabbitConsumer, RabbitConsumer>();
builder.Services.AddHostedService<HostedRabbitConsumer>();

#endregion

var app = builder.Build();

#region Database Initialization

await using var scope = app.Services.CreateAsyncScope();
var dbContext = scope.ServiceProvider.GetRequiredService<SentimentatorContext>();
await dbContext.Database.EnsureCreatedAsync();

#endregion

#region CORS

app.UseCors("AllowSpecificOrigin");

#endregion

#region Auth

app.UseAuthentication();
app.UseAuthorization();

#endregion

#region OpenAPI

app.MapOpenApi()
    .AllowAnonymous()
    .CacheOutput();

#endregion

#region Scalar

app.MapScalarApiReference(options =>
{
    options.ShowSidebar = true;
    options.TagSorter = TagSorter.Alpha;
    options.Layout = ScalarLayout.Modern;
    options.HideClientButton = true;
    options.Theme = ScalarTheme.BluePlanet;
    options.DocumentDownloadType = DocumentDownloadType.Both;
    options.DefaultOpenAllTags = false;
    options.DefaultHttpClient =
        new KeyValuePair<ScalarTarget, ScalarClient>(ScalarTarget.CSharp, ScalarClient.HttpClient);
    options.Authentication = new ScalarAuthenticationOptions
    {
        PreferredSecuritySchemes = [IdentityConstants.BearerScheme]
    };
}).AllowAnonymous();

#endregion

#region Endpoints

app.MapGet("/", () => Results.Redirect("/scalar/v1"))
    .ExcludeFromDescription()
    .AllowAnonymous();

app.MapGroup("/hello")
    .MapHello()
    .WithTags("Hello");

app.MapGroup("/analyses")
    .MapAnalysis()
    .WithTags("Analysis")
    .AddOperationBearerSecuritySchemeWithAuthorization();

app.MapGroup("/account")
    .MapIdentityApi<User>()
    .WithTags("Identity")
    .AllowAnonymous();

app.MapGroup("/users")
    .MapUser()
    .WithTags("Users")
    .AddOperationBearerSecuritySchemeWithAuthorization();

#endregion

app.Run();