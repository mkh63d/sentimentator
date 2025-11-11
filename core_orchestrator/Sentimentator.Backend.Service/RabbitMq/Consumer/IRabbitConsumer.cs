namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public interface IRabbitConsumer
{
    Task Receive(string queue);
}