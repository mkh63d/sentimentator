using System.Security.Claims;
using System.Web;
using JetBrains.Annotations;
using Sentimentator.Backend.Model;
using Sentimentator.Backend.Service;
using Microsoft.AspNetCore.Identity;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class AnalysisGroup
{
    public static RouteGroupBuilder AnalysisGetGroup(this RouteGroupBuilder group)
    {
        group.MapGet("/{id:guid}", async (Guid id, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            try
            {
                var user = await userManager.GetUserAsync(principal);
                var analysis = await service.Get(id);
                return user!.Analyses.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
            }
            catch
            {
                return Results.NotFound();
            }
        });

        group.MapGet("/{url}", async (string url, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            try
            {
                url = HttpUtility.UrlDecode(url);
                var user = await userManager.GetUserAsync(principal);
                var analysis = await service.Get(url);
                return user!.Analyses.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
            }
            catch
            {
                return Results.NotFound();
            }
        });

        group.MapGet("/", async (UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            try
            {
                var user = await userManager.GetUserAsync(principal);
                return Results.Ok(user!.Analyses.ToList());
            }
            catch
            {
                return Results.InternalServerError();
            }
        });

        return group;
    }

    public static RouteGroupBuilder AnalysisPostGroup(this RouteGroupBuilder group)
    {
        group.MapPost("/{url}", async (string url, AnalysisService service, ClaimsPrincipal principal) =>
        {
            url = HttpUtility.UrlDecode(url);
            var result = await service.Add(principal, url);
            return result is not null
                ? Results.Created($"/analyses/{result.Id}", result)
                : Results.InternalServerError();
        });

        return group;
    }

    public static RouteGroupBuilder AnalysisPutGroup(this RouteGroupBuilder group)
    {
        group.MapPut("/{id:guid}", async (Guid id, AnalysisDetails content, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            var user = await userManager.GetUserAsync(principal);
            var result = false;
            if (user!.Analyses.Contains(await service.Get(id)))
            {
                result = await service.Update(id, content.Transcription, content.Summary);
            }
            return result ? Results.Ok() : Results.NotFound();
        });

        group.MapPut("/{url}", async (string url, AnalysisDetails content, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            url = HttpUtility.UrlDecode(url);
            var user = await userManager.GetUserAsync(principal);
            var result = false;
            if (user!.Analyses.Contains(await service.Get(url)))
            {
                result = await service.Update(url, content.Transcription, content.Summary);
            }
            return result ? Results.Ok() : Results.NotFound();
        });

        return group;
    }

    public static RouteGroupBuilder AnalysisDeleteGroup(this RouteGroupBuilder group)
    {
        group.MapDelete("/{id:guid}", async (Guid id, AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            var user = await userManager.GetUserAsync(principal);
            var result = false;
            if (user!.Analyses.Contains(await service.Get(id)))
            {
                result = await service.Remove(id);
            }
            return result ? Results.Ok() : Results.NotFound();
        });

        return group;
    }

    [UsedImplicitly]
    private record AnalysisDetails(string Transcription, string Summary);
}