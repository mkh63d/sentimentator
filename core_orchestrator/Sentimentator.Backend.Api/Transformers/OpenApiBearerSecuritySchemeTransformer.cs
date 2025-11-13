using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.Identity;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace Sentimentator.Backend.Api.Transformers;

public static class OpenApiBearerSecuritySchemeTransformer
{
    public static OpenApiOptions AddBearerSecurityScheme(this OpenApiOptions options)
    {
        return options.AddDocumentTransformer<BearerSecuritySchemeTransformer>();
    }

    internal sealed class BearerSecuritySchemeTransformer(
        IAuthenticationSchemeProvider authenticationSchemeProvider)
        : IOpenApiDocumentTransformer
    {
        public async Task TransformAsync(
            OpenApiDocument document, OpenApiDocumentTransformerContext context,
            CancellationToken cancellationToken)
        {
            var schemes = await authenticationSchemeProvider.GetAllSchemesAsync();
            if (schemes.All(s => s.Name != IdentityConstants.BearerScheme)) return;

            document.Components ??= new OpenApiComponents();
            document.Components.SecuritySchemes = new Dictionary<string, IOpenApiSecurityScheme>
            {
                [IdentityConstants.BearerScheme] = new OpenApiSecurityScheme
                {
                    Name = "Authorization",
                    Type = SecuritySchemeType.Http,
                    Scheme = "bearer",
                    In = ParameterLocation.Header,
                    Description = "Identity Bearer Token. Paste the Bearer token below.",
                    BearerFormat = "Opaque"
                }
            };
        }
    }
}