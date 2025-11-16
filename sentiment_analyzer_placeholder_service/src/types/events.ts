export interface TranscriptCompletedEvent {
    analysis_id: string;
    original_url: string;
    video_id: string;
    transcript_id: string;
    text: string;
    language: string;
    duration: number;
    created_at: string;
}

export interface SentimentFeatures {
    camera?: number;
    battery?: number;
    screen?: number;
    performance?: number;
    design?: number;
    price?: number;
    [key: string]: number | undefined;
}

export interface SentimentAnalysis {
    overall: number;
    features: SentimentFeatures;
    confidence: number;
}

export interface AnalysisCompletedEvent {
    analysis_id: string;
    original_url: string;
    video_id: string;
    transcript: string;
    sentiment: SentimentAnalysis;
    data: SentimentAnalysis;
}
