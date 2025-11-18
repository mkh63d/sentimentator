using System.Security.Claims;
using System.Text.Json;
using System.Web;
using Microsoft.AspNetCore.Identity;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;
using Sentimentator.Backend.Service.Analysis;
using Sentimentator.Backend.Service.RabbitMq.Publisher;
using Sentimentator.Backend.Service.RabbitMq.Publisher.Model;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class AnalysisGroup
{
    public static RouteGroupBuilder MapAnalysis(this RouteGroupBuilder group)
    {
        group.AnalysisGetGroup()
            .AnalysisPostGroup()
            .AnalysisPutGroup()
            .AnalysisDeleteGroup();
        return group;
    }

    private static RouteGroupBuilder AnalysisGetGroup(this RouteGroupBuilder group)
    {
        group.MapGet("/{id:guid}",
                async (Guid id, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
                {
                    var user = await userManager.GetUserAsync(principal);
                    var analysis = await service.Get(id);
                    if (analysis is null) return Results.NotFound();
                    return user!.Analyses.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
                })
            .Produces<Analysis>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/{url}",
                async (string url, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
                {
                    var uri = HttpUtility.UrlDecode(url).ToUri();
                    switch (uri.result)
                    {
                        case UrlResult.InvalidFormat:
                            return Results.BadRequest("Invalid URL format");
                        case UrlResult.UnsupportedDomain:
                            return Results.BadRequest("Unsupported URL domain");
                        case UrlResult.Valid:
                        default:
                            break;
                    }
                    
                    var user = await userManager.GetUserAsync(principal);
                    var analysis = await service.Get(uri.uri!);
                    if (analysis is null) return Results.NotFound();
                    return user!.Analyses.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
                })
            .Produces<Analysis>()
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/", async (AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
            {
                var user = await userManager.GetUserAsync(principal);
                await service.Get();
                var analyses = user!.Analyses.Select(x => x).ToList();
                return Results.Ok(analyses);
            })
            .Produces<List<Analysis>>();
        return group;
    }

    private static RouteGroupBuilder AnalysisPostGroup(this RouteGroupBuilder group)
    {
        group.MapPost("/{url}",
                async (string url, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal, IRabbitPublisher publisher) =>
                {
                    var uri = HttpUtility.UrlDecode(url).ToUri();
                    switch (uri.result)
                    {
                        case UrlResult.InvalidFormat:
                            return Results.BadRequest("Invalid URL format");
                        case UrlResult.UnsupportedDomain:
                            return Results.BadRequest("Unsupported URL domain");
                        case UrlResult.Valid:
                        default:
                            break;
                    }
                    
                    var user = await userManager.GetUserAsync(principal);
                    await service.Get();
                    var a = user!.Analyses.FirstOrDefault(x => x.Url.Equals(uri.uri));
                    if (a is { Summary: not null, Transcription: not null }) return Results.Ok(a);
                    
                    var result = await service.Add(principal, uri.uri!);
                    if (result is null) return Results.InternalServerError();

                    var downloadRequest = new TranscriptRequest
                    {
                        AnalysisId = result.Id,
                        Url = result.Url,
                        Timestamp = DateTime.UtcNow
                    };
                    await publisher.Send("video.processing", "video.transcript.requested", JsonSerializer.Serialize(downloadRequest));

                    return Results.Accepted($"/analyses/{result.Id}", result);
                })
            .Produces<Analysis>()
            .Produces<Analysis>(StatusCodes.Status202Accepted)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status500InternalServerError);

        return group;
    }

    private static RouteGroupBuilder AnalysisPutGroup(this RouteGroupBuilder group)
    {
        group.MapPut("/{id:guid}",
                async (Guid id, AnalysisDetails content, AnalysisService service, UserManager<User> userManager,
                    ClaimsPrincipal principal) =>
                {
                    var user = await userManager.GetUserAsync(principal);
                    var result = false;
                    var analysis = await service.Get(id);
                    if (analysis is null) return Results.NotFound();
                    if (user!.Analyses.Contains(analysis))
                        result = await service.Update(id, content.Transcription, content.Summary, DateTime.UtcNow);
                    return result ? Results.Ok() : Results.NotFound();
                })
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        group.MapPut("/{url}",
                async (string url, AnalysisDetails content, AnalysisService service, UserManager<User> userManager,
                    ClaimsPrincipal principal) =>
                {
                    var uri = HttpUtility.UrlDecode(url).ToUri();
                    switch (uri.result)
                    {
                        case UrlResult.InvalidFormat:
                            return Results.BadRequest("Invalid URL format");
                        case UrlResult.UnsupportedDomain:
                            return Results.BadRequest("Unsupported URL domain");
                        case UrlResult.Valid:
                        default:
                            break;
                    }
                    
                    var user = await userManager.GetUserAsync(principal);
                    var result = false;
                    var analysis = await service.Get(uri.uri!);
                    if (analysis is null) return Results.NotFound();
                    if (user!.Analyses.Contains(analysis))
                        result = await service.Update(uri.uri!, content.Transcription, content.Summary, DateTime.UtcNow);
                    return result ? Results.Ok() : Results.NotFound();
                })
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status400BadRequest)
            .Produces(StatusCodes.Status404NotFound);

        return group;
    }

    private static RouteGroupBuilder AnalysisDeleteGroup(this RouteGroupBuilder group)
    {
        group.MapDelete("/{id:guid}",
                async (Guid id, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
                {
                    var user = await userManager.GetUserAsync(principal);
                    var result = false;
                    var analysis = await service.Get(id);
                    if (analysis is null) return Results.NotFound();
                    if (user!.Analyses.Contains(analysis)) result = await service.Remove(id);
                    return result ? Results.Ok() : Results.NotFound();
                })
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        return group;
    }
}