using System.Security.Claims;
using JetBrains.Annotations;
using Microsoft.AspNetCore.Identity;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Api.EndpointGroups;

[UsedImplicitly]
public record UserApiKey(string ApiKey);

public static class UserGroup
{
    extension(RouteGroupBuilder group)
    {
        public RouteGroupBuilder MapUser()
        {
            group.UserPostGroup()
                .UserDeleteGroup();
            return group;
        }

        private RouteGroupBuilder UserPostGroup()
        {
            group.MapPost("/openai-api-key",
                    async (UserApiKey apiKey, UserManager<User> userManager,
                        RoleManager<IdentityRole<Guid>> roleManager, ClaimsPrincipal principal) =>
                    {
                        if (string.IsNullOrWhiteSpace(apiKey.ApiKey))
                            return Results.BadRequest("API key cannot be empty");
                        var user = await userManager.GetUserAsync(principal);
                        user!.OpenAiApiKey = apiKey.ApiKey;
                        var result = await userManager.UpdateAsync(user);
                        if (!result.Succeeded)
                            return Results.InternalServerError(result.Errors);

                        var canRequestAnalysisRoleExists = await roleManager.RoleExistsAsync("CanRequestAnalysis");
                        if (!canRequestAnalysisRoleExists)
                        {
                            var canRequestAnalysisRole = new IdentityRole<Guid>("CanRequestAnalysis");
                            _ = await roleManager.CreateAsync(canRequestAnalysisRole);
                        }

                        var canRequestAnalysisUsers = await userManager.GetUsersInRoleAsync("CanRequestAnalysis");
                        if (canRequestAnalysisUsers.Contains(user)) return Results.Ok();

                        result = await userManager.AddToRoleAsync(user, "CanRequestAnalysis");
                        return result.Succeeded ? Results.Ok() : Results.InternalServerError(result.Errors);
                    })
                .Produces(StatusCodes.Status200OK)
                .Produces(StatusCodes.Status500InternalServerError);

            return group;
        }

        private RouteGroupBuilder UserDeleteGroup()
        {
            group.MapDelete("/openai-api-key",
                    async (UserManager<User> userManager, RoleManager<IdentityRole<Guid>> roleManager,
                        ClaimsPrincipal principal) =>
                    {
                        var user = await userManager.GetUserAsync(principal);
                        user!.OpenAiApiKey = null;
                        var result = await userManager.UpdateAsync(user);
                        if (!result.Succeeded) return Results.InternalServerError(result.Errors);

                        var canRequestAnalysisRoleExists = await roleManager.RoleExistsAsync("CanRequestAnalysis");
                        if (!canRequestAnalysisRoleExists) return Results.Ok();

                        var adminRoleExists = await roleManager.RoleExistsAsync("Admin");
                        if (adminRoleExists)
                        {
                            var adminUsers = await userManager.GetUsersInRoleAsync("Admin");
                            if (adminUsers.Contains(user)) return Results.Ok();
                        }

                        var canRequestAnalysisUsers = await userManager.GetUsersInRoleAsync("CanRequestAnalysis");
                        if (!canRequestAnalysisUsers.Contains(user)) return Results.Ok();

                        result = await userManager.RemoveFromRoleAsync(user, "CanRequestAnalysis");
                        return result.Succeeded ? Results.Ok() : Results.InternalServerError(result.Errors);
                    })
                .Produces(StatusCodes.Status200OK)
                .Produces(StatusCodes.Status500InternalServerError);

            return group;
        }
    }
}