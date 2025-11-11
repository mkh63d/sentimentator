using Microsoft.AspNetCore.Identity;

namespace Sentimentator.Backend.Model;

public class User : IdentityUser<Guid>
{
    public List<Analysis> Analyses { get; } = [];
}