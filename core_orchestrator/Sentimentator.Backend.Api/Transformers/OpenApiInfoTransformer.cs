using Microsoft.AspNetCore.Authentication;
using Microsoft.AspNetCore.OpenApi;
using Microsoft.OpenApi;

namespace Sentimentator.Backend.Api.Transformers;

public static class OpenApiInfoTransformer
{
    public static OpenApiOptions AddInfo(this OpenApiOptions options)
    {
        return options.AddDocumentTransformer<InfoTransformer>();
    }

    internal sealed class InfoTransformer(
        IAuthenticationSchemeProvider authenticationSchemeProvider)
        : IOpenApiDocumentTransformer
    {
        public async Task TransformAsync(
            OpenApiDocument document, OpenApiDocumentTransformerContext context,
            CancellationToken cancellationToken)
        {
            document.Info = new OpenApiInfo
            {
                Title = "Sentimentator backend API",
                Version = "v1",
                Description = "API for retrieving sentiment analysis data"
            };
        }
    }
}