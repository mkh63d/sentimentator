using Sentimentator.Backend.Api.Transformers;

namespace Sentimentator.Backend.Api.EndpointGroups;

public static class HelloGroup
{
    extension(RouteGroupBuilder group)
    {
        public RouteGroupBuilder MapHello()
        {
            group.HelloGetGroup();
            return group;
        }

        private RouteGroupBuilder HelloGetGroup()
        {
            group.MapGet("/", () => Results.Ok("Hello World!"))
                .Produces<string>()
                .AllowAnonymous();

            group.MapGet("/with-auth", () => Results.Ok("Hello World! You're authenticated!"))
                .Produces<string>()
                .AddOperationBearerSecuritySchemeWithAuthorization();

            group.MapGet("/with-upload-auth", () => Results.Ok("Hello World! You can request analyses!"))
                .Produces<string>()
                .Produces(StatusCodes.Status403Forbidden)
                .AddOperationBearerSecuritySchemeWithAuthorization("CanRequestAnalysis");

            return group;
        }
    }
}