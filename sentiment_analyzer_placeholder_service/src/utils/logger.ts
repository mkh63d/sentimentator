export type LogLevel = "debug" | "info" | "warn" | "error";

const logLevels: Record<LogLevel, number> = {
    debug: 0,
    info: 1,
    warn: 2,
    error: 3,
};

export class Logger {
    private currentLevel: LogLevel;

    constructor(level: LogLevel = "info") {
        this.currentLevel = level;
    }

    private shouldLog(level: LogLevel): boolean {
        return logLevels[level] >= logLevels[this.currentLevel];
    }

    private formatTimestamp(): string {
        return new Date().toISOString();
    }

    debug(message: string, meta?: Record<string, unknown>): void {
        if (this.shouldLog("debug")) {
            console.log(
                `[${this.formatTimestamp()}] DEBUG: ${message}`,
                meta ? JSON.stringify(meta) : ""
            );
        }
    }

    info(message: string, meta?: Record<string, unknown>): void {
        if (this.shouldLog("info")) {
            console.log(
                `[${this.formatTimestamp()}] INFO: ${message}`,
                meta ? JSON.stringify(meta) : ""
            );
        }
    }

    warn(message: string, meta?: Record<string, unknown>): void {
        if (this.shouldLog("warn")) {
            console.warn(
                `[${this.formatTimestamp()}] WARN: ${message}`,
                meta ? JSON.stringify(meta) : ""
            );
        }
    }

    error(message: string, error?: Error | Record<string, unknown>): void {
        if (this.shouldLog("error")) {
            const errorMeta =
                error instanceof Error
                    ? { message: error.message, stack: error.stack }
                    : error;
            console.error(
                `[${this.formatTimestamp()}] ERROR: ${message}`,
                errorMeta ? JSON.stringify(errorMeta) : ""
            );
        }
    }
}

export const logger = new Logger((process.env.LOG_LEVEL as LogLevel) || "info");
