using System.Security.Claims;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Service.Analysis;

public class AnalysisService(
    SentimentatorContext context,
    UserManager<User> userManager,
    ILogger<AnalysisService> logger)
{
    private async Task<Model.Analysis?> Add(Uri uri)
    {
        try
        {
            var analysis = new Model.Analysis { Id = Guid.NewGuid(), Url = uri, CreatedAt = DateTime.UtcNow };
            context.Analyses.Add(analysis);
            await context.SaveChangesAsync();
            logger.LogInformation("Added analysis {AnalysisId} for URL {Url}", analysis.Id, uri.ToString());
            return analysis;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to add analysis for URL {Url}: {Error}", uri.ToString(), ex.Message);
            return null;
        }
    }

    public async Task<Model.Analysis?> Add(ClaimsPrincipal principal, Uri uri)
    {
        try
        {
            var analysis = await Add(uri);
            if (analysis is null) throw new Exception();
            var user = await userManager.GetUserAsync(principal);
            user!.Analyses.Add(analysis);
            await userManager.UpdateAsync(user);
            logger.LogInformation("Added analysis {AnalysisId} for URL {URL} associated with user {UserId}",
                analysis.Id, uri.ToString(), user.Id);
            return analysis;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to add analysis for URL {Url}: {Error}", uri.ToString(), ex.Message);
            return null;
        }
    }

    public async Task<Model.Analysis?> Get(Guid id)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Id == id);
        if (a is null) logger.LogWarning("Analysis with ID {AnalysisId} not found", id);
        else logger.LogInformation("Retrieved analysis {AnalysisId}", id);
        return a;
    }

    public async Task<Model.Analysis?> Get(Uri uri)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Url.Equals(uri));
        if (a is null) logger.LogWarning("Analysis with URL {Url} not found", uri.ToString());
        else logger.LogInformation("Retrieved analysis {AnalysisId} for URL {Url}", a.Id, uri.ToString());
        return a;
    }

    public async Task<List<Model.Analysis>> Get()
    {
        return await context.Analyses.ToListAsync();
    }

    public async Task<bool> Update(Guid id, string transcription, string summary, DateTime updatedAt)
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
            a.UpdatedAt = updatedAt;
            await context.SaveChangesAsync();
            logger.LogInformation("Updated analysis {AnalysisId}", id);
            return true;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to update analysis {AnalysisId}: {Exception}", id, ex.Message);
            return false;
        }
    }

    public async Task<bool> Update(Uri uri, string transcription, string summary, DateTime updatedAt)
    {
        try
        {
            var a = await Get(uri);
            if (a is null)
            {
                logger.LogError("Failed to update analysis for URL {Url}: not found", uri.ToString());
                return false;
            }

            a.Transcription = transcription;
            a.Summary = summary;
            a.UpdatedAt = updatedAt;
            await context.SaveChangesAsync();
            logger.LogInformation("Updated analysis {AnalysisId} for URL {Url}", a.Id, uri.ToString());
            return true;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to update analysis for URL {Url}: {Exception}", uri.ToString(), ex.Message);
            return false;
        }
    }

    public async Task<bool> Remove(Guid id)
    {
        try
        {
            var analysis = new Model.Analysis { Id = id, Url = new Uri("about:blank") };
            context.Analyses.Attach(analysis);
            context.Analyses.Remove(analysis);
            await context.SaveChangesAsync();
            logger.LogInformation("Removed analysis {AnalysisId}", id);
            return true;
        }
        catch (Exception ex)
        {
            logger.LogError("Failed to remove analysis {AnalysisId}: {Exception}", id, ex.Message);
            return false;
        }
    }
}