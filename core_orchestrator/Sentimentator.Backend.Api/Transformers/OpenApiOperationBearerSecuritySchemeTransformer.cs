using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace Sentimentator.Backend.Api.Transformers;

public static class OpenApiOperationBearerSecuritySchemeTransformer
{
    public static TBuilder AddOperationBearerSecuritySchemeWithAuthorization<TBuilder>(this TBuilder builder,
        string policyName = "")
        where TBuilder : IEndpointConventionBuilder
    {
        builder.AddOpenApiOperationTransformer(new OperationBearerSecuritySchemeTransformer().TransformAsync)
            .ProducesProblem(StatusCodes.Status401Unauthorized);

        if (!string.IsNullOrEmpty(policyName)) builder.RequireAuthorization(policyName);

        return builder;
    }

    private sealed class OperationBearerSecuritySchemeTransformer : IOpenApiOperationTransformer
    {
        public Task TransformAsync(
            OpenApiOperation operation, OpenApiOperationTransformerContext context,
            CancellationToken cancellationToken)
        {
            operation.Security ??= [];
            var reference = new OpenApiSecuritySchemeReference(IdentityConstants.BearerScheme, context.Document);
            if (operation.Security.Any(r => r.ContainsKey(reference))) return Task.CompletedTask;

            operation.Security.Add(new OpenApiSecurityRequirement
            {
                [reference] = []
            });
            return Task.CompletedTask;
        }
    }
}