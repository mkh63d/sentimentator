using Microsoft.AspNetCore.Authorization;
using Microsoft.AspNetCore.Identity;
using Scalar.AspNetCore;
using Sentimentator.Backend;
using Sentimentator.Backend.Api.EndpointGroups;
using Sentimentator.Backend.Api.Transformers;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;
using Sentimentator.Backend.Service.RabbitMq.Consumer;
using Sentimentator.Backend.Service.RabbitMq.Publisher;

var builder = WebApplication.CreateBuilder(args);
builder.Logging.ClearProviders();
builder.Logging.AddConsole();

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

        options.SignIn.RequireConfirmedEmail = false;
    })
    .AddRoles<IdentityRole<Guid>>()
    .AddEntityFrameworkStores<SentimentatorContext>()
    .AddApiEndpoints();
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
builder.Services.AddScoped<RabbitPublisher>();
builder.Services.AddScoped<AnalysisService>();
builder.Services.AddSingleton<IRabbitConsumer, RabbitConsumer>();
builder.Services.AddHostedService<HostedRabbitConsumer>();
#endregion

var app = builder.Build();

app.UseAuthentication();
app.UseAuthorization();

app.MapOpenApi()
    .AllowAnonymous()
    .CacheOutput();

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

app.UseHttpsRedirection();

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

app.Run();