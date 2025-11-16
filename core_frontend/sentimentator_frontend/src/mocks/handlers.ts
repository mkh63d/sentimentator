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
  http.post('*/account/login', async () => {
    return HttpResponse.json<AccessTokenResponse>({
      accessToken: 'mock-access-token',
      expiresIn: 3600,
      refreshToken: 'mock-refresh-token',
      tokenType: 'Bearer',
    });
  }),

  http.post('*/account/register', async () => {
    return new HttpResponse(null, { status: 200 });
  }),

  // Analysis endpoints
  http.get('*/analyses', () => {
    return HttpResponse.json([
      {
        id: '123e4567-e89b-12d3-a456-426614174000',
        url: 'https://example.com',
        transcription:
          'Just picked up the new ZetaPhone X2, and after a week, I’m pretty impressed. The battery easily lasts me a full day with streaming and social apps. Performance is smooth overall, though I did notice slight stutters when switching between heavy games. The camera takes sharp photos in daylight, but nighttime shots are just okay. The build quality feels premium, but the phone is a bit heavier than I expected.',
        summary: {
          battery: 'positive',
          performance: 'neutral',
          camera: 'positive',
          camera_by_night: 'neutral',
          build_quality: 'positive',
          weight: 'negative',
        },
      },
    ]);
  }),

  http.post('*/analyses/:url', ({ params }) => {
    return HttpResponse.json({
      id: '123e4567-e89b-12d3-a456-426614174000',
      url: params.url,
    });
  }),

  http.get('*/analyses/:id', ({ params }) => {
    return HttpResponse.json({
      id: params.id,
      url: 'https://example.com',
      transcription:
        'Just picked up the new ZetaPhone X2, and after a week, I’m pretty impressed. The battery easily lasts me a full day with streaming and social apps. Performance is smooth overall, though I did notice slight stutters when switching between heavy games. The camera takes sharp photos in daylight, but nighttime shots are just okay. The build quality feels premium, but the phone is a bit heavier than I expected.',
      summary: {
        battery: 'positive',
        performance: 'neutral',
        camera: 'positive',
        camera_by_night: 'neutral',
        build_quality: 'positive',
        weight: 'negative',
      },
    });
  }),

  http.delete('*/analyses/:id', ({ params }) => {
    return new HttpResponse(null, { status: 200 });
  }),
];
