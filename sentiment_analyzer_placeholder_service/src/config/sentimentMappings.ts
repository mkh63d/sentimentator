/**
 * Predefined sentiment analysis payloads based on video URLs
 * Maps specific video URLs to hardcoded sentiment analysis results
 */

export interface SentimentFeatures {
    camera?: number;
    battery?: number;
    screen?: number;
    performance?: number;
    design?: number;
    price?: number;
    [key: string]: number | undefined;
}

export interface SentimentResult {
    overall: number;
    features: SentimentFeatures;
    confidence: number;
}

export interface PredefinedPayload {
    urlPattern: string | RegExp; // URL to match
    sentiment: SentimentResult;
    summary: string;
}

/**
 * Predefined sentiment payloads for specific video URLs
 * First exact match wins, then falls back to pattern matching, then default
 */
export const predefinedPayloads: PredefinedPayload[] = [
    // iPhone 17 Review: No Asterisks! - VERY POSITIVE
    {
        urlPattern: "rng_yUSwrgU",
        sentiment: {
            overall: 0.92,
            features: {
                camera: 0.98,
                battery: 0.75,
                screen: 0.95,
                performance: 0.96,
                design: 0.93,
                price: 0.68,
            },
            confidence: 0.97,
        },
        summary:
            "Outstanding iPhone 17 review highlighting exceptional camera capabilities and revolutionary design with no significant compromises. Performance and display are industry-leading. Pricing premium is justified by the innovation.",
    },

    // Samsung Galaxy S26 Ultra - VERY POSITIVE
    {
        urlPattern: "TTBkme1KYi0",
        sentiment: {
            overall: 0.89,
            features: {
                camera: 0.92,
                battery: 0.88,
                screen: 0.94,
                performance: 0.91,
                design: 0.87,
                price: 0.72,
            },
            confidence: 0.96,
        },
        summary:
            "Excellent Samsung Galaxy S26 Ultra delivers premium features across the board. The 7-year wait was worthwhile with significant improvements in camera system, battery endurance, and display quality. Represents a major leap forward in technology.",
    },

    // Generic positive review
    {
        urlPattern: "8kYyKpB4R0Y",
        sentiment: {
            overall: 0.78,
            features: {
                camera: 0.82,
                battery: 0.8,
                screen: 0.85,
                performance: 0.8,
                design: 0.78,
                price: 0.72,
            },
            confidence: 0.88,
        },
        summary:
            "Solid device with well-rounded performance across all categories. Good balance between features and price point.",
    },

    // Generic positive review
    {
        urlPattern: "yJbnqSf-Q_I",
        sentiment: {
            overall: 0.75,
            features: {
                camera: 0.78,
                battery: 0.76,
                screen: 0.82,
                performance: 0.74,
                design: 0.75,
                price: 0.7,
            },
            confidence: 0.85,
        },
        summary:
            "Competent smartphone with reliable performance. Camera system is impressive for the category. Battery life meets expectations.",
    },

    // Generic positive review
    {
        urlPattern: "KV8RF-PcV8I",
        sentiment: {
            overall: 0.81,
            features: {
                camera: 0.85,
                battery: 0.82,
                screen: 0.88,
                performance: 0.83,
                design: 0.8,
                price: 0.75,
            },
            confidence: 0.89,
        },
        summary:
            "Excellent all-around smartphone with standout features in photography and display quality.",
    },

    // Generic positive review
    {
        urlPattern: "2F0WGf9c4xk",
        sentiment: {
            overall: 0.76,
            features: {
                camera: 0.8,
                battery: 0.77,
                screen: 0.83,
                performance: 0.76,
                design: 0.74,
                price: 0.68,
            },
            confidence: 0.86,
        },
        summary:
            "Good smartphone offering solid performance with particular strength in visual display quality.",
    },

    // Generic positive review
    {
        urlPattern: "nBIaZxXldh4",
        sentiment: {
            overall: 0.73,
            features: {
                camera: 0.75,
                battery: 0.72,
                screen: 0.8,
                performance: 0.73,
                design: 0.71,
                price: 0.65,
            },
            confidence: 0.84,
        },
        summary:
            "Decent device with acceptable performance across key areas. Display quality is a highlight.",
    },

    // Generic positive review
    {
        urlPattern: "hU5Tx4MJ0w0",
        sentiment: {
            overall: 0.79,
            features: {
                camera: 0.83,
                battery: 0.81,
                screen: 0.86,
                performance: 0.79,
                design: 0.77,
                price: 0.73,
            },
            confidence: 0.87,
        },
        summary:
            "Strong performer with particularly good camera and screen quality. Well-designed device.",
    },

    // Generic positive review
    {
        urlPattern: "F9k7Q2FxEQs",
        sentiment: {
            overall: 0.77,
            features: {
                camera: 0.81,
                battery: 0.78,
                screen: 0.84,
                performance: 0.77,
                design: 0.75,
                price: 0.71,
            },
            confidence: 0.86,
        },
        summary:
            "Reliable smartphone with good camera performance and solid display.",
    },

    // Generic positive review
    {
        urlPattern: "ITJVzQXIUco",
        sentiment: {
            overall: 0.8,
            features: {
                camera: 0.84,
                battery: 0.81,
                screen: 0.87,
                performance: 0.81,
                design: 0.79,
                price: 0.74,
            },
            confidence: 0.88,
        },
        summary:
            "Excellent device with balanced features and good overall value proposition.",
    },

    // Generic positive review
    {
        urlPattern: "9-8MOcnxIFk",
        sentiment: {
            overall: 0.74,
            features: {
                camera: 0.77,
                battery: 0.74,
                screen: 0.81,
                performance: 0.74,
                design: 0.72,
                price: 0.66,
            },
            confidence: 0.85,
        },
        summary: "Good smartphone with capable camera and display system.",
    },
];

/**
 * Get sentiment result for a given video URL
 * Matches against predefined payloads first by exact video ID, then by pattern
 * Returns the matched payload or a default neutral sentiment
 */
export function getSentimentForUrl(videoUrl: string): {
    result: SentimentResult;
    summary: string;
} {
    // Extract video ID from various URL formats
    const videoId = extractVideoId(videoUrl);

    // Try exact match on video ID first (from YouTube watch?v= or Instagram reel ID)
    for (const payload of predefinedPayloads) {
        if (typeof payload.urlPattern === "string") {
            if (
                videoUrl.includes(payload.urlPattern) ||
                videoId.includes(payload.urlPattern)
            ) {
                return {
                    result: payload.sentiment,
                    summary: payload.summary,
                };
            }
        } else if (payload.urlPattern instanceof RegExp) {
            if (payload.urlPattern.test(videoUrl)) {
                return {
                    result: payload.sentiment,
                    summary: payload.summary,
                };
            }
        }
    }

    // Default neutral sentiment for unknown URLs
    return {
        result: {
            overall: 0.7,
            features: {
                camera: 0.7,
                battery: 0.7,
                screen: 0.7,
                performance: 0.7,
                design: 0.7,
                price: 0.7,
            },
            confidence: 0.8,
        },
        summary:
            "Mixed review with balanced positive and negative aspects across device features.",
    };
}

/**
 * Extract video ID from various URL formats
 * Handles YouTube, Instagram, and direct URL formats
 */
function extractVideoId(url: string): string {
    // YouTube: https://www.youtube.com/watch?v=VIDEO_ID or youtu.be/VIDEO_ID
    const youtubeMatch = url.match(
        /(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]+)/
    );
    if (youtubeMatch) return youtubeMatch[1];

    // Instagram: https://www.instagram.com/reels/REEL_ID/
    const instagramMatch = url.match(/instagram\.com\/reels\/([a-zA-Z0-9_-]+)/);
    if (instagramMatch) return instagramMatch[1];

    // Fallback: return the URL itself
    return url;
}
