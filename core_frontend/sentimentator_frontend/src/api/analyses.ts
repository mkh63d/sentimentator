import { api } from '@/lib/client';
import type { components } from '@/generated/schema';

export type Analysis = components['schemas']['Analysis'];
export type AnalysisDetails = components['schemas']['AnalysisDetails'];

/**
 * GET /analyses
 * Returns all analyses for the authenticated user.
 */
export async function getAnalyses(): Promise<Analysis[]> {
  const { data, error } = await api.GET('/analyses');
  if (error) throw new Error(error.detail ?? 'Failed to fetch analyses');
  if (!data) throw new Error('No data returned from server');
  return data;
}

/**
 * POST /analyses/{url}
 * Creates or triggers a new analysis.
 */
export async function createAnalysis(url: string): Promise<Analysis> {
  const { data, error } = await api.POST('/analyses/{url}', {
    params: { path: { url } },
  });
  if (error) throw new Error(error.detail ?? 'Failed to create analysis');
  if (!data) throw new Error('No data returned from server');
  return data;
}

/**
 * GET /analyses/{id}
 * Fetch a single analysis by its UUID.
 */
export async function getAnalysisById(id: string): Promise<Analysis> {
  const { data, error } = await api.GET('/analyses/{id}', {
    params: { path: { id } },
  });
  if (error) throw new Error(error.detail ?? 'Failed to fetch analysis');
  if (!data) throw new Error('No data returned from server');
  return data;
}

/**
 * PUT /analyses/{id}
 * Updates an analysis.
 */
export async function updateAnalysis(id: string, details: AnalysisDetails): Promise<void> {
  const { error } = await api.PUT('/analyses/{id}', {
    params: { path: { id } },
    body: details,
  });
  if (error) throw new Error(error.detail ?? 'Failed to update analysis');
}

/**
 * DELETE /analyses/{id}
 * Deletes an analysis by UUID.
 */
export async function deleteAnalysis(id: string): Promise<void> {
  const { error } = await api.DELETE('/analyses/{id}', {
    params: { path: { id } },
  });
  if (error) throw new Error(error.detail ?? 'Failed to delete analysis');
}
