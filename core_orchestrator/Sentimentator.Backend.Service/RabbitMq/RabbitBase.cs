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

    private readonly IConnection _connection;

    private readonly ConnectionFactory _factory = new() { Uri = new Uri(RabbitMqUri) };
    internal readonly IChannel Channel;

    internal readonly ILogger<RabbitBase> Logger;

    internal RabbitBase(ILogger<RabbitBase> logger)
    {
        _connection = _factory.CreateConnectionAsync().Result;
        Channel = _connection.CreateChannelAsync().Result;

        Logger = logger;
    }

    public void Dispose()
    {
        Channel.CloseAsync().Wait();
        _connection.CloseAsync().Wait();
        Channel.DisposeAsync().AsTask().Wait();
        _connection.DisposeAsync().AsTask().Wait();

        GC.SuppressFinalize(this);
    }
}