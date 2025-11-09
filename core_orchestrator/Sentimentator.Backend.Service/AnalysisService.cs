using System.Security.Claims;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Service;

public class AnalysisService(SentimentatorContext context, UserManager<User> userManager)
{
    public async Task<Analysis?> Add(string url)
    {
        try
        {
            var analysis = new Analysis { Id = Guid.NewGuid(), Url = url };
            context.Analyses.Add(analysis);
            await context.SaveChangesAsync();
            return analysis;
        }
        catch
        {
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
            return analysis;
        }
        catch (Exception e)
        {
            return null;
        }
    }

    public async Task<Analysis> Get(Guid id)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Id == id);
        return a ?? throw new KeyNotFoundException("Analysis not found");
    }

    public async Task<Analysis> Get(string url)
    {
        var a = await context.Analyses.FirstOrDefaultAsync(x => x.Url == url);
        return a ?? throw new KeyNotFoundException("Analysis not found");
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
            a.Transcription = transcription;
            a.Summary = summary;
            await context.SaveChangesAsync();
            return true;
        }
        catch
        {
            return false;
        }
    }

    public async Task<bool> Update(string url, string transcription, string summary)
    {
        try
        {
            var a = await Get(url);
            a.Transcription = transcription;
            a.Summary = summary;
            await context.SaveChangesAsync();
            return true;
        }
        catch
        {
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
            return true;
        }
        catch
        {
            return false;
        }
    }
}