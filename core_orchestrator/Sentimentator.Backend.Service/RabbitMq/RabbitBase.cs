using RabbitMQ.Client;

namespace Sentimentator.Backend.Service.RabbitMq;

public class RabbitBase
{
    internal readonly ConnectionFactory Factory = new() { Uri = new Uri("amqp://guest:guest@localhost:5672/") }; // TODO: change to environment variable
}