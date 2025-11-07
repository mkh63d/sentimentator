using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace Sentimentator.Backend.Api.Transformers;

public static class OpenApiTagTransformer
{
    public static OpenApiOptions AddTags(this OpenApiOptions options)
    {
        return options.AddDocumentTransformer<TagTransformer>();
    }

    internal sealed class TagTransformer(
        IAuthenticationSchemeProvider authenticationSchemeProvider)
        : IOpenApiDocumentTransformer
    {
        public async Task TransformAsync(
            OpenApiDocument document, OpenApiDocumentTransformerContext context,
            CancellationToken cancellationToken)
        {
            document.Tags = new HashSet<OpenApiTag>
            {
                new()
                {
                    Name = "Hello",
                    Description = "Test endpoints"
                },
                new()
                {
                    Name = "Analysis",
                    Description = "Endpoints for analysis storage and retrieval"
                },
                new()
                {
                    Name = "Identity",
                    Description = "Auto-generated endpoints for user registration, login, and account management",
                    ExternalDocs = new OpenApiExternalDocs
                    {
                        Description = "Bearer Token Authentication in ASP.NET Core",
                        Url = new Uri(
                            "https://devblogs.microsoft.com/dotnet/bearer-token-authentication-in-asp-net-core/")
                    }
                }
            };
        }
    }
}