using System.Text.Json.Serialization;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.Mvc;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class AdminGroup
{
    extension(RouteGroupBuilder group)
    {
        public RouteGroupBuilder MapAdmin()
        {
            group.AdminPostGroup()
                .AdminDeleteGroup();
            return group;
        }
        
        private RouteGroupBuilder AdminPostGroup()
        {
            group.MapPost("/role", async ([FromBody] RoleDetails roleDetails, UserManager<User> userManager, RoleManager<IdentityRole<Guid>> roleManager) =>
            {
                var user = await userManager.FindByEmailAsync(roleDetails.UserEmail);
                if (user is null) return Results.NotFound("User not found");
                
                var roleExists = await roleManager.RoleExistsAsync(roleDetails.RoleName);
                if (!roleExists)
                {
                    var role = new IdentityRole<Guid>(roleDetails.RoleName);
                    var roleResult = await roleManager.CreateAsync(role);
                    if (!roleResult.Succeeded) return Results.InternalServerError(roleResult.Errors);
                }
                
                var result = await userManager.AddToRoleAsync(user, roleDetails.RoleName);
                return !result.Succeeded ? Results.InternalServerError(result.Errors) : Results.Ok();
            })
            .Produces(StatusCodes.Status404NotFound)
            .Produces(StatusCodes.Status500InternalServerError)
            .Produces(StatusCodes.Status200OK);
            
            return group;
        }
        
        private RouteGroupBuilder AdminDeleteGroup()
        {
            group.MapDelete("/role",
                    async ([FromBody] RoleDetails roleDetails, UserManager<User> userManager,
                        RoleManager<IdentityRole<Guid>> roleManager) =>
                    {
                        var user = await userManager.FindByEmailAsync(roleDetails.UserEmail);
                        if (user is null) return Results.NotFound("User not found");

                        var roleExists = await roleManager.RoleExistsAsync(roleDetails.RoleName);
                        if (!roleExists) return Results.BadRequest("Role does not exist");

                        var result = await userManager.RemoveFromRoleAsync(user, roleDetails.RoleName);
                        return !result.Succeeded ? Results.InternalServerError(result.Errors) : Results.Ok();
                    })
                .Produces(StatusCodes.Status404NotFound)
                .Produces(StatusCodes.Status500InternalServerError)
                .Produces(StatusCodes.Status400BadRequest)
                .Produces(StatusCodes.Status200OK);
            
            return group;
        }
    }
}

internal class RoleDetails
{
    [JsonPropertyName("user_email")] public required string UserEmail { get; set; }
    [JsonPropertyName("role_name")] public required string RoleName { get; set; }
}