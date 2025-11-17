using Microsoft.Extensions.Logging;
using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq;

public class RabbitBase : IDisposable
{
    internal const string ExchangeName = "sentimentator.topic";
    
    private const string RabbitMqUriEnvVar = "ORCHESTRATOR_RABBITMQ_URI";
    private const string RabbitMqDefaultUri = "amqp://guest:guest@localhost:5672/";
    private static readonly string RabbitMqUri = Environment.GetEnvironmentVariable(RabbitMqUriEnvVar)
                                                 ?? RabbitMqDefaultUri;

    private static readonly ConnectionFactory Factory = new() { Uri = new Uri(RabbitMqUri) };
    private static readonly IConnection Connection = Factory.CreateConnectionAsync().Result;
    internal readonly IChannel Channel = Connection.CreateChannelAsync().Result;
    
    internal readonly ILogger<RabbitBase> Logger;

    internal RabbitBase(ILogger<RabbitBase> logger)
    {
        Logger = logger;
    }
    
    public void Dispose()
    {
        Channel.CloseAsync().Wait();
        Connection.CloseAsync().Wait();
        Channel.DisposeAsync().AsTask().Wait();
        Connection.DisposeAsync().AsTask().Wait();

        GC.SuppressFinalize(this);
    }
}