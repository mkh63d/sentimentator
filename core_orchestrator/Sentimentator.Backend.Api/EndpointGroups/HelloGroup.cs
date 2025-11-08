using Sentimentator.Backend.Api.Transformers;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class HelloGroup
{
    public static RouteGroupBuilder HelloGetGroup(this RouteGroupBuilder group)
    {
        group.MapGet("/", () => Results.Ok("Hello World!"))
            .Produces(StatusCodes.Status200OK)
            .AllowAnonymous();

        group.MapGet("/with-auth", () => Results.Ok("Hello World! You're authenticated!"))
            .Produces(StatusCodes.Status200OK)
            .AddOperationBearerSecuritySchemeWithAuthorization();

        return group;
    }
}