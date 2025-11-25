namespace Sentimentator.Backend.Service.RabbitMq.Publisher;

public interface IRabbitPublisher
{
    Task Send(string queue, string routingKey, string message);
}