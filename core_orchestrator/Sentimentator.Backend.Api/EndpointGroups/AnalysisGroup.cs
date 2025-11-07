using System.Security.Claims;
using System.Web;
using JetBrains.Annotations;
using Microsoft.AspNetCore.Mvc;
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
                var analysis = service.Get(id);
                return user!.Analyses!.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
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
                var analysis = service.Get(url);
                return user!.Analyses!.Contains(analysis) ? Results.Ok(analysis) : Results.NotFound();
            }
            catch
            {
                return Results.NotFound();
            }
        });

        group.MapGet("/", async (AnalysisService service, UserManager<User> userManager, ClaimsPrincipal principal) =>
        {
            try
            {
                var user = await userManager.GetUserAsync(principal);
                return Results.Ok(user!.Analyses!.ToList());
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
        group.MapPost("/{url}", (string url, AnalysisService service) =>
        {
            url = HttpUtility.UrlDecode(url);
            var result = service.Add(url);
            return result is not null ? Results.Created($"/analyses/{result.Id}", result) : Results.InternalServerError();
        });
        
        return group;
    }

    [UsedImplicitly]
    private record AnalysisDetails(string Transcription, string Summary);
    
    public static RouteGroupBuilder AnalysisPutGroup(this RouteGroupBuilder group)
    {
        group.MapPut("/{id:guid}", (Guid id, AnalysisDetails content, AnalysisService service) =>
        {
            var result = service.Update(id, content.Transcription, content.Summary);
            return result ? Results.Ok() : Results.NotFound();
        });
        
        group.MapPut("/{url}", (string url, AnalysisDetails content, AnalysisService service) =>
        {
            url = HttpUtility.UrlDecode(url);
            var result = service.Update(url, content.Transcription, content.Summary);
            return result ? Results.Ok() : Results.NotFound();
        });
        
        return group;
    }

    public static RouteGroupBuilder AnalysisDeleteGroup(this RouteGroupBuilder group)
    {
        group.MapDelete("/{id:guid}", (Guid id, AnalysisService service) =>
        {
            var result = service.Remove(id);
            return result ? Results.Ok() : Results.NotFound();
        });
        
        return group;
    }
}