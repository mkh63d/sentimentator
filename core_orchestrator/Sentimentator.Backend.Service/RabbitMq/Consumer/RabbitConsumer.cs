using System.Text;
using Microsoft.Extensions.Logging;
using RabbitMQ.Client;
using RabbitMQ.Client.Events;

namespace Sentimentator.Backend.Service.RabbitMq.Consumer;

public class RabbitConsumer : RabbitBase, IRabbitConsumer, IDisposable
{
    private readonly ILogger<RabbitConsumer> _logger;
    private readonly IConnection _connection;
    private readonly IChannel _channel;

    public RabbitConsumer(ILogger<RabbitConsumer> logger)
    {
        _logger = logger;
        _connection = Factory.CreateConnectionAsync().Result;
        _channel = _connection.CreateChannelAsync().Result;
    }
    
    public async Task Receive(string queue)
    {
        await _channel.QueueDeclareAsync(queue: queue, durable: false, exclusive: false, autoDelete: false,
            arguments: null);
        
        var consumer = new AsyncEventingBasicConsumer(_channel);
        consumer.ReceivedAsync += async (model, ea) =>
        {
            var body = ea.Body.ToArray();
            var message = Encoding.UTF8.GetString(body);
            _logger.LogInformation("Received message from queue {} at {} with message \"{}\"", queue, DateTime.UtcNow.ToLongTimeString(), message);

            await _channel.BasicAckAsync(ea.DeliveryTag, false);
        };

        await _channel.BasicConsumeAsync(queue: queue, autoAck: false, consumer: consumer);
    }

    public void Dispose()
    {
        _channel.CloseAsync().Wait();
        _connection.CloseAsync().Wait();
        _channel.DisposeAsync();
        _connection.DisposeAsync();
    }
}