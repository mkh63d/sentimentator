import dotenv from "dotenv";
dotenv.config({ path: ".env" });
import amqp from "amqplib";
import type { ConsumeMessage } from "amqplib";
import { logger } from "../utils/logger.js";
import {
    TranscriptCompletedEvent,
    AnalysisCompletedEvent,
} from "../types/events.js";
import { getSentimentForUrl } from "../config/sentimentMappings.js";

const RABBITMQ_URL = process.env.RABBITMQ_URL;
if (!RABBITMQ_URL) {
    throw new Error("RABBITMQ_URL environment variable is not set");
}

const TRANSCRIPT_QUEUE = process.env.TRANSCRIPT_QUEUE || "video.transcript";
const ANALYSIS_QUEUE = process.env.ANALYSIS_QUEUE || "video.analysis";
const ANALYSIS_EXCHANGE =
    process.env.ANALYSIS_EXCHANGE || "sentimentator.topic";
const ANALYSIS_ROUTING_KEY =
    process.env.ANALYSIS_ROUTING_KEY || "video.analysis.completed";

export class SentimentAnalyzerService {
    private connection: any = null;
    private channel: any = null;
    private isRunning = false;

    async connect(): Promise<void> {
        try {
            logger.info("Connecting to RabbitMQ...", { url: RABBITMQ_URL });
            this.connection = await amqp.connect(RABBITMQ_URL as string);
            this.channel = await this.connection.createChannel();

            await this.channel.prefetch(1);

            logger.info("Connected to RabbitMQ successfully");
        } catch (error) {
            logger.error("Failed to connect to RabbitMQ", error as Error);
            throw error;
        }
    }

    async setupQueuesAndExchanges(): Promise<void> {
        if (!this.channel) {
            throw new Error("Channel not initialized");
        }

        try {
            await this.channel.assertExchange(ANALYSIS_EXCHANGE, "topic", {
                durable: true,
                exclusive: false,
            });

            await this.channel.assertQueue(TRANSCRIPT_QUEUE, {
                durable: true,
                exclusive: false,
            });

            await this.channel.bindQueue(
                TRANSCRIPT_QUEUE,
                ANALYSIS_EXCHANGE,
                "video.transcript.completed"
            );

            await this.channel.assertQueue(ANALYSIS_QUEUE, {
                durable: true,
                exclusive: false,
            });

            await this.channel.bindQueue(
                ANALYSIS_QUEUE,
                ANALYSIS_EXCHANGE,
                ANALYSIS_ROUTING_KEY
            );

            logger.info("Queues and exchanges configured", {
                exchange: ANALYSIS_EXCHANGE,
                inputQueue: TRANSCRIPT_QUEUE,
                outputQueue: ANALYSIS_QUEUE,
                incomingRoutingKey: "video.transcript.completed",
                outgoingRoutingKey: ANALYSIS_ROUTING_KEY,
            });
        } catch (error) {
            logger.error(
                "Failed to setup queues and exchanges",
                error as Error
            );
            throw error;
        }
    }

    async start(): Promise<void> {
        if (!this.channel) {
            throw new Error("Channel not initialized");
        }

        try {
            this.isRunning = true;
            logger.info("Starting to consume messages from queue...", {
                queue: TRANSCRIPT_QUEUE,
            });

            await this.channel.consume(
                TRANSCRIPT_QUEUE,
                (msg: any) => this.handleMessage(msg),
                { noAck: false }
            );

            logger.info("Sentiment analyzer service started successfully");
        } catch (error) {
            logger.error("Failed to start consuming messages", error as Error);
            this.isRunning = false;
            throw error;
        }
    }

    private async handleMessage(msg: ConsumeMessage | null): Promise<void> {
        if (!msg) {
            logger.warn("Received null message from queue");
            return;
        }

        if (!this.channel) {
            logger.error("Channel not available for processing message");
            this.channel?.nack(msg, false, true);
            return;
        }

        try {
            const transcriptEvent: TranscriptCompletedEvent = JSON.parse(
                msg.content.toString()
            );

            logger.info("Processing transcript event", {
                videoId: transcriptEvent.video_id,
                videoUrl: transcriptEvent.original_url,
            });

            const { result: sentiment, summary } = getSentimentForUrl(
                transcriptEvent.original_url
            );

            const analysisEvent: AnalysisCompletedEvent = {
                analysis_id: transcriptEvent.analysis_id,
                original_url: transcriptEvent.original_url,
                video_id: transcriptEvent.video_id,
                transcript: transcriptEvent.text,
                sentiment: sentiment,
            };

            await this.publishAnalysisResult(analysisEvent);

            this.channel.ack(msg);

            logger.info("Successfully processed and published analysis", {
                videoId: transcriptEvent.video_id,
                analysisId: transcriptEvent.analysis_id,
            });
        } catch (error) {
            logger.error("Error processing message", error as Error);
            if (this.channel) {
                this.channel.nack(msg, false, true);
            }
        }
    }

    private async publishAnalysisResult(
        event: AnalysisCompletedEvent
    ): Promise<void> {
        if (!this.channel) {
            throw new Error("Channel not initialized");
        }

        try {
            const message = Buffer.from(JSON.stringify(event));

            logger.info("Publishing analysis result to RabbitMQ", {
                exchange: ANALYSIS_EXCHANGE,
                routingKey: ANALYSIS_ROUTING_KEY,
                videoId: event.video_id,
                analysisId: event.analysis_id,
            });

            const published = this.channel.publish(
                ANALYSIS_EXCHANGE,
                ANALYSIS_ROUTING_KEY,
                message,
                {
                    persistent: true,
                    contentType: "application/json",
                    timestamp: Date.now(),
                }
            );

            if (!published) {
                logger.warn(
                    "Message may not have been published successfully - channel buffer full",
                    {
                        videoId: event.video_id,
                        routingKey: ANALYSIS_ROUTING_KEY,
                    }
                );
            } else {
                logger.info("Analysis result published successfully", {
                    videoId: event.video_id,
                    routingKey: ANALYSIS_ROUTING_KEY,
                    queueName: ANALYSIS_QUEUE,
                });
            }
        } catch (error) {
            logger.error("Failed to publish analysis result", error as Error);
            throw error;
        }
    }

    async stop(): Promise<void> {
        this.isRunning = false;

        try {
            if (this.channel) {
                await this.channel.close();
                logger.info("Channel closed");
            }

            if (this.connection) {
                await this.connection.close();
                logger.info("Connection closed");
            }

            logger.info("Sentiment analyzer service stopped");
        } catch (error) {
            logger.error("Error during shutdown", error as Error);
        }
    }
}
