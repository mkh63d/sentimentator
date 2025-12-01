import { api } from '@/lib/client';

export async function changeAPIKey(newApiKey: string) {
  const { data, error } = await api.PUT('/users/openai-api-key', {
    body: { apiKey: newApiKey },
  });
  if (error) {
    throw new Error(error.title ?? 'Failed to update API key');
  }
  return data;
}
