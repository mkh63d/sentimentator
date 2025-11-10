using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq;

public class RabbitBase
{
    private const string RabbitMqUriEnvVar = "ORCHESTRATOR_RABBITMQ_URI";
    private const string RabbitMqDefaultUri = "amqp://guest:guest@localhost:5672/";

    private static readonly string RabbitMqUri = Environment.GetEnvironmentVariable(RabbitMqUriEnvVar)
                                                 ?? RabbitMqDefaultUri;

    internal readonly ConnectionFactory
        Factory = new() { Uri = new Uri(RabbitMqUri) };
}