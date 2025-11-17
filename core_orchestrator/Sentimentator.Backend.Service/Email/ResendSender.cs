using Microsoft.AspNetCore.Identity.UI.Services;
using Microsoft.Extensions.Logging;
using Resend;

namespace Sentimentator.Backend.Service.Email;

public class ResendSender(ILogger<ResendSender> logger) : IEmailSender
{
    public async Task SendEmailAsync(string email, string subject, string htmlMessage)
    {
        const string resendApiKeyEnvVar = "ORCHESTRATOR_RESEND_API_KEY";
        var resendApiKey = Environment.GetEnvironmentVariable(resendApiKeyEnvVar);
        if (resendApiKey == null)
        {
            logger.LogError("Resend API key not found in environment variable {EnvVar}", resendApiKeyEnvVar);
            return;
        }

        IResend resend;
        try
        {
            resend = ResendClient.Create(resendApiKey);
        }
        catch (ResendException e)
        {
            logger.LogError("Initializing resend client failed: {Error}", e.Message);
            return;
        }

        ResendResponse<Guid> response;
        try
        {
            response = await resend.EmailSendAsync( new EmailMessage
            {
                From = "sentimentator@bdymek.com",
                To = email,
                Subject = subject,
                HtmlBody = htmlMessage
            } );
        }
        catch (ResendException e)
        {
            logger.LogError("Resend email failed: {Error}", e.Message);
            return;
        }
        
        logger.LogInformation("Sent email to {Email} with subject {Subject}. Successful: {Success}", email, subject, response.Success);
    }
}