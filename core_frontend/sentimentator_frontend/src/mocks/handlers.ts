import { http, HttpResponse } from 'msw';

// Example data types
export interface User {
  id: number;
  name: string;
}

// API Response Types
interface AccessTokenResponse {
  tokenType?: string | null;
  accessToken: string;
  expiresIn: number;
  refreshToken: string;
}

interface AnalysisDetails {
  transcription: string;
  summary: string;
}

export const handlers = [
  // Authentication endpoints
  http.post('/account/login', async () => {
    return HttpResponse.json<AccessTokenResponse>({
      accessToken: 'mock-access-token',
      expiresIn: 3600,
      refreshToken: 'mock-refresh-token',
      tokenType: 'Bearer',
    });
  }),

  http.post('/account/register', async () => {
    return new HttpResponse(null, { status: 200 });
  }),

  // Analysis endpoints
  http.get('/analyses', () => {
    return HttpResponse.json([
      {
        id: '123e4567-e89b-12d3-a456-426614174000',
        url: 'https://example.com',
        transcription: 'Sample transcription',
        summary: 'Sample summary',
      },
    ]);
  }),

  http.post('/analyses/:url', ({ params }) => {
    return HttpResponse.json({
      id: '123e4567-e89b-12d3-a456-426614174000',
      url: params.url,
      transcription: '',
      summary: '',
    });
  }),

  http.get('/analyses/:id', ({ params }) => {
    return HttpResponse.json({
      id: params.id,
      url: 'https://example.com',
      transcription: 'Sample transcription',
      summary: 'Sample summary',
    });
  }),

  http.put('/analyses/:id', async ({ params, request }) => {
    const body = (await request.json()) as AnalysisDetails;
    return HttpResponse.json({
      id: params.id,
      ...body,
    });
  }),
];
