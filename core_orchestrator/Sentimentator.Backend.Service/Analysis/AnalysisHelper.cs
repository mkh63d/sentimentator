using JetBrains.Annotations;

namespace Sentimentator.Backend.Service;

[UsedImplicitly]
public record AnalysisDetails(string Transcription, string Summary);

public enum UrlResult
{
    Valid,
    InvalidFormat,
    UnsupportedDomain
}

public static class AnalysisHelper
{
    private static readonly List<string> SupportedYoutubeDomains =
    [
        "www.youtube.com", "youtube.com",
        "www.youtu.be", "youtu.be",
        "m.youtube.com"
    ];
    
    private static readonly List<string> SupportedInstagramDomains =
    [
        "www.instagram.com", "instagram.com"
    ];
    
    private static readonly List<string> SupportedTikTokDomains =
    [
        "www.tiktok.com", "tiktok.com",
        "vm.tiktok.com"
    ];
    
    public static (Uri? uri, UrlResult result) ToUri(this string url)
    {
        Uri uri;
        try
        {
            uri = new Uri(url);
        }
        catch
        {
            return (null, UrlResult.InvalidFormat);
        }
        if (!SupportedYoutubeDomains.Contains(uri.Host) &&
            !SupportedInstagramDomains.Contains(uri.Host) &&
            !SupportedTikTokDomains.Contains(uri.Host))
        {
            return (null, UrlResult.UnsupportedDomain);
        }
        return (uri, UrlResult.Valid);
    }
}