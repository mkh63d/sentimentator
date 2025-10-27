const API_URL = import.meta.env.VITE_BACKEND_API || 'http://localhost:8000';

export interface Analysis {
  id: string;
  text: string;
  sentiment: 'positive' | 'negative' | 'neutral';
  score: number;
  createdAt: string;
}

// Fetch all analyses but not in their full structural data detail
export const getAnalyses = async (): Promise<Analysis[]> => {
  const response = await fetch(`${API_URL}/api/analyses`);

  if (!response.ok) {
    throw new Error('Failed to fetch analyses');
  }

  return response.json();
};

// Fetch a single analysis by ID with full structural data detail
export const getAnalysis = async (id: string): Promise<Analysis> => {
  const response = await fetch(`${API_URL}/api/analyses/${id}`);

  if (!response.ok) {
    throw new Error('Failed to fetch analysis');
  }

  return response.json();
};

// Send request to create a new analysis based on a link
export const createAnalysisFromLink = async (link: string): Promise<Analysis> => {
  const response = await fetch(`${API_URL}/api/analyses`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ link }),
  });

  if (!response.ok) {
    throw new Error('Failed to create analysis from link');
  }

  return response.json();
};
export const createAnalysis = async (text: string): Promise<Analysis> => {
  const response = await fetch(`${API_URL}/api/analyses`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ text }),
  });

  if (!response.ok) {
    throw new Error('Failed to create analysis');
  }

  return response.json();
};
