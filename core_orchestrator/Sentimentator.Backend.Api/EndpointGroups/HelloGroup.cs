using Sentimentator.Backend.Api.Transformers;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class HelloGroup
{
    public static RouteGroupBuilder HelloGetGroup(this RouteGroupBuilder group)
    {
        group.MapGet("/", () => Results.Ok("Hello World!"))
            .Produces<string>()
            .AllowAnonymous();

        group.MapGet("/with-auth", () => Results.Ok("Hello World! You're authenticated!"))
            .Produces<string>()
            .AddOperationBearerSecuritySchemeWithAuthorization();

        return group;
    }
}