using System.Security.Claims;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Service;

public class AnalysisService(
    SentimentatorContext context,
    UserManager<User> userManager,
    ILogger<AnalysisService> logger)
{
    public async Task<Analysis?> Add(string url)
    {
        try
        {
            var analysis = new Analysis { Id = Guid.NewGuid(), Url = url };
            context.Analyses.Add(analysis);
            await context.SaveChangesAsync();
            logger.LogInformation("Added analysis {AnalysisId} for URL {Url}", analysis.Id, url);
            return analysis;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to add analysis for URL {Url}: {Error}", url, ex.Message);
            return null;
        }
    }

    public async Task<Analysis?> Add(ClaimsPrincipal principal, string url)
    {
        try
        {
            var analysis = await Add(url);
            if (analysis is null) throw new Exception();
            var user = await userManager.GetUserAsync(principal);
            user!.Analyses.Add(analysis);
            await userManager.UpdateAsync(user);
            logger.LogInformation("Added analysis {AnalysisId} for URL {URL} associated with user {UserId}",
                analysis.Id, url, user.Id);
            return analysis;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to add analysis for URL {Url}: {Error}", url, ex.Message);
            return null;
        }
    }

    public async Task<Analysis?> Get(Guid id)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Id == id);
        if (a is null) logger.LogWarning("Analysis with ID {AnalysisId} not found", id);
        else logger.LogInformation("Retrieved analysis {AnalysisId}", id);
        return a;
    }

    public async Task<Analysis?> Get(string url)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Url == url);
        if (a is null) logger.LogWarning("Analysis with URL {Url} not found", url);
        else logger.LogInformation("Retrieved analysis {AnalysisId} for URL {Url}", a.Id, url);
        return a;
    }

    public async Task<List<Analysis>> Get()
    {
        return await context.Analyses.ToListAsync();
    }

    public async Task<bool> Update(Guid id, string transcription, string summary)
    {
        try
        {
            var a = await Get(id);
            if (a is null)
            {
                logger.LogError("Failed to update analysis {AnalysisId}: not found", id);
                return false;
            }

            a.Transcription = transcription;
            a.Summary = summary;
            await context.SaveChangesAsync();
            logger.LogInformation("Updated analysis {AnalysisId}", id);
            return true;
        }
        catch
        {
            logger.LogError("Failed to update analysis {AnalysisId}", id);
            return false;
        }
    }

    public async Task<bool> Update(string url, string transcription, string summary)
    {
        try
        {
            var a = await Get(url);
            if (a is null)
            {
                logger.LogError("Failed to update analysis for URL {Url}: not found", url);
                return false;
            }

            a.Transcription = transcription;
            a.Summary = summary;
            await context.SaveChangesAsync();
            logger.LogInformation("Updated analysis {AnalysisId} for URL {Url}", a.Id, url);
            return true;
        }
        catch
        {
            logger.LogError("Failed to update analysis for URL {Url}", url);
            return false;
        }
    }

    public async Task<bool> Remove(Guid id)
    {
        try
        {
            var analysis = new Analysis { Id = id, Url = "" };
            context.Analyses.Attach(analysis);
            context.Analyses.Remove(analysis);
            await context.SaveChangesAsync();
            logger.LogInformation("Removed analysis {AnalysisId}", id);
            return true;
        }
        catch
        {
            logger.LogError("Failed to remove analysis {AnalysisId}", id);
            return false;
        }
    }
}