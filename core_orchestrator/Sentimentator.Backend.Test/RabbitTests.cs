using Microsoft.Extensions.Logging;
using Microsoft.Extensions.Logging.Abstractions;
using Sentimentator.Backend.Service.RabbitMq.Consumer;
using Sentimentator.Backend.Service.RabbitMq.Publisher;

namespace Sentimentator.Backend.Test;

public class RabbitTests
{
    [Fact]
    public async Task ShouldCompleteWhenPublishingCorrectRequest()
    {
        var consumer = new RabbitConsumer(new Logger<RabbitConsumer>(new NullLoggerFactory()));
        var publisher = new RabbitPublisher(new Logger<RabbitPublisher>(new NullLoggerFactory()));

        var message = "";
        var exception = await Record.ExceptionAsync(async () =>
        {
            await consumer.Receive("test-queue", "test.key", async msg =>
            {
                message = msg;
                return await Task.FromResult(true);
            });
        });

        var publishException = await Record.ExceptionAsync(async () =>
        {
            await publisher.Send("test-queue", "test.key", "Hello, RabbitMQ!");
        });

        Assert.Null(exception);
        Assert.Null(publishException);
        await Task.Delay(500);
        Assert.Equal("Hello, RabbitMQ!", message);
    }
}