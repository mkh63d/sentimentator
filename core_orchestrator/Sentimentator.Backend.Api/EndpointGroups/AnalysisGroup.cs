using System.Security.Claims;
using System.Text.Json;
using System.Web;
using JetBrains.Annotations;
using Microsoft.AspNetCore.Identity;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;
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
                    url = HttpUtility.UrlDecode(url);
                    var user = await userManager.GetUserAsync(principal);
                    var analysis = await service.Get(url);
                    if (analysis is null) return Results.NotFound();
                    return user!.Analyses.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
                })
            .Produces<Analysis>()
            .Produces(StatusCodes.Status404NotFound);

        group.MapGet("/", async (UserManager<User> userManager, ClaimsPrincipal principal) =>
            {
                var user = await userManager.GetUserAsync(principal);
                return Results.Ok(user!.Analyses.ToList());
            })
            .Produces<List<Analysis>>();
        return group;
    }

    private static RouteGroupBuilder AnalysisPostGroup(this RouteGroupBuilder group)
    {
        group.MapPost("/{url}",
                async (string url, AnalysisService service, ClaimsPrincipal principal, RabbitPublisher publisher) =>
                {
                    url = HttpUtility.UrlDecode(url);
                    var result = await service.Add(principal, url);
                    if (result is null) return Results.InternalServerError();

                    var downloadRequest = new TranscriptRequest
                    {
                        AnalysisId = result.Id,
                        Url = new Uri(result.Url),
                        Timestamp = DateTime.UtcNow
                    };
                    await publisher.Send("video.transcript.requested", JsonSerializer.Serialize(downloadRequest));

                    return Results.Accepted($"/analyses/{result.Id}", result);
                })
            .Produces<Analysis>(StatusCodes.Status202Accepted)
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
                        result = await service.Update(id, content.Transcription, content.Summary);
                    return result ? Results.Ok() : Results.NotFound();
                })
            .Produces(StatusCodes.Status200OK)
            .Produces(StatusCodes.Status404NotFound);

        group.MapPut("/{url}",
                async (string url, AnalysisDetails content, AnalysisService service, UserManager<User> userManager,
                    ClaimsPrincipal principal) =>
                {
                    url = HttpUtility.UrlDecode(url);
                    var user = await userManager.GetUserAsync(principal);
                    var result = false;
                    var analysis = await service.Get(url);
                    if (analysis is null) return Results.NotFound();
                    if (user!.Analyses.Contains(analysis))
                        result = await service.Update(url, content.Transcription, content.Summary);
                    return result ? Results.Ok() : Results.NotFound();
                })
            .Produces(StatusCodes.Status200OK)
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

    [UsedImplicitly]
    private record AnalysisDetails(string Transcription, string Summary);
}