import dotenv from "dotenv";
import { logger } from "./utils/logger.js";
import { SentimentAnalyzerService } from "./services/SentimentAnalyzerService.js";

dotenv.config({ path: ".env.local" });
dotenv.config({ path: ".env" });

const service = new SentimentAnalyzerService();

async function shutdown(signal: string): Promise<void> {
    logger.info(`Received ${signal}, shutting down gracefully...`);
    await service.stop();
    process.exit(0);
}

async function main(): Promise<void> {
    try {
        logger.info("Starting Sentiment Analyzer Placeholder Service");

        await service.connect();

        await service.setupQueuesAndExchanges();

        await service.start();

        logger.info("Service is running. Press Ctrl+C to exit.");
    } catch (error) {
        logger.error("Fatal error during startup", error as Error);
        process.exit(1);
    }
}

process.on("SIGTERM", () => shutdown("SIGTERM"));
process.on("SIGINT", () => shutdown("SIGINT"));

process.on("uncaughtException", (error: Error) => {
    logger.error("Uncaught exception", error);
    process.exit(1);
});

process.on("unhandledRejection", (reason: unknown) => {
    logger.error("Unhandled rejection", reason as Error);
    process.exit(1);
});

main().catch((error) => {
    logger.error("Application startup failed", error as Error);
    process.exit(1);
});
