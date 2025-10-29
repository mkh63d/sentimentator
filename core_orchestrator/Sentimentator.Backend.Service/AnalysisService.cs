using Sentimentator.Backend.Model;

namespace Sentimentator.Backend.Service;

public class AnalysisService(SentimentatorContext context)
{
    public Analysis? Add(string url)
    {
        try
        {
            var analysis = new Analysis { Id = Guid.NewGuid(), Url = url };
            context.Analyses.Add(analysis);
            context.SaveChanges();
            return analysis;
        }
        catch
        {
            return null;
        }
    }

    public Analysis Get(Guid id)
    {
        var a = context.Analyses.FirstOrDefault(x => x.Id == id);
        return a ?? throw new KeyNotFoundException("Analysis not found");
    }
    
    public Analysis Get(string url)
    {
        var a = context.Analyses.FirstOrDefault(x => x.Url == url);
        return a ?? throw new KeyNotFoundException("Analysis not found");
    }

    public List<Analysis> Get()
    {
        return context.Analyses.ToList();
    }

    public bool Update(Guid id, string transcription, string summary)
    {
        try
        {
            var a = Get(id);
            a.Transcription = transcription;
            a.Summary = summary;
            context.SaveChanges();
            return true;
        }
        catch
        {
            return false;
        }
    }
    
    public bool Update(string url, string transcription, string summary)
    {
        try
        {
            var a = Get(url);
            a.Transcription = transcription;
            a.Summary = summary;
            context.SaveChanges();
            return true;
        }
        catch
        {
            return false;
        }
    }

    public bool Remove(Guid id)
    {
        try
        {
            var analysis = new Analysis { Id = id, Url = "" };
            context.Analyses.Attach(analysis);
            context.Analyses.Remove(analysis);
            context.SaveChanges();
            return true;
        }
        catch
        {
            return false;
        }
    }
}