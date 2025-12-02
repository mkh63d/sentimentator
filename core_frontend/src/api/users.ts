import { api } from '@/lib/client';

export async function changeAPIKey(newApiKey: string) {
  const { data, error } = await api.POST('/users/openai-api-key', {
    body: { apiKey: newApiKey },
  });
  if (error) {
    throw new Error(error.title ?? 'Failed to update API key');
  }
  return data;
}

export async function resetAPIKey() {
  const { data, error } = await api.DELETE('/users/openai-api-key');
  if (error) {
    throw new Error(error.title ?? 'Failed to reset API key');
  }
  return data;
}
