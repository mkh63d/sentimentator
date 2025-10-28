import { apiFetch, apiPost, Endpoints } from './api';

export interface Analysis {
  id: string;
  text: string;
  sentiment: 'positive' | 'negative' | 'neutral';
  score: number;
  createdAt: string;
}

// Fetch all analyses but not in their full structural data detail
export const getAnalyses = async (): Promise<Analysis[]> => {
  return apiFetch(Endpoints.ANALYSES, { method: 'GET' });
};

// Fetch a single analysis by ID with full structural data detail
export const getAnalysis = async (id: string): Promise<Analysis> => {
  return apiFetch(Endpoints.ANALYSES + `/${id}`, { method: 'GET' });
};

// Send request to create a new analysis based on a link
export const createAnalysisFromLink = async (link: string): Promise<Analysis> => {
  return apiPost(Endpoints.ANALYSES, JSON.stringify({ link }));
};

export const createAnalysis = async (text: string): Promise<Analysis> => {
  return apiPost(Endpoints.ANALYSES, JSON.stringify({ text }));
};
