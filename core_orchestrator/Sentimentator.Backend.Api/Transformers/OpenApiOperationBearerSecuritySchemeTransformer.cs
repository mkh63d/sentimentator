using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace Sentimentator.Backend.Api.Transformers;

public static class OpenApiOperationBearerSecuritySchemeTransformer
{
    public static TBuilder AddOperationBearerSecuritySchemeWithAuthorization<TBuilder>(this TBuilder builder)
        where TBuilder : IEndpointConventionBuilder
    {
        return builder
            .AddOpenApiOperationTransformer(new OperationBearerSecuritySchemeTransformer().TransformAsync)
            .RequireAuthorization()
            .ProducesProblem(StatusCodes.Status401Unauthorized);
    }

    internal sealed class OperationBearerSecuritySchemeTransformer : IOpenApiOperationTransformer
    {
        public async Task TransformAsync(
            OpenApiOperation operation, OpenApiOperationTransformerContext context,
            CancellationToken cancellationToken)
        {
            operation.Security ??= [];
            operation.Security.Add(new OpenApiSecurityRequirement
            {
                [new OpenApiSecuritySchemeReference(IdentityConstants.BearerScheme, context.Document)] = []
            });
        }
    }
}