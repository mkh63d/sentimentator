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
            group.UserPutGroup();
            return group;
        }

        private RouteGroupBuilder UserPutGroup()
        {
            group.MapPut("/openai-api-key",
                async (UserApiKey apiKey, UserManager<User> userManager, RoleManager<IdentityRole<Guid>> roleManager, ClaimsPrincipal principal) =>
                {
                    if (string.IsNullOrWhiteSpace(apiKey.ApiKey)) return Results.BadRequest("API key cannot be empty");
                    var user = await userManager.GetUserAsync(principal);
                    user!.OpenAiApiKey = apiKey.ApiKey;
                    var result = await userManager.UpdateAsync(user);
                    if (!result.Succeeded)
                    {
                        return Results.InternalServerError(result.Errors);
                    }
                    
                    var canRequestAnalysisRoleExists = await roleManager.RoleExistsAsync("CanRequestAnalysis");
                    if (!canRequestAnalysisRoleExists)
                    {
                        var canRequestAnalysisRole = new IdentityRole<Guid>("CanRequestAnalysis");
                        _ = await roleManager.CreateAsync(canRequestAnalysisRole);
                    }
                    result = await userManager.AddToRoleAsync(user, "CanRequestAnalysis");
                    
                    return result.Succeeded ? Results.Ok() : Results.InternalServerError(result.Errors);
                })
                .Produces(StatusCodes.Status200OK)
                .Produces(StatusCodes.Status500InternalServerError);
            
            return group;
        }
    }
}