import createClient, { type Middleware } from 'openapi-fetch';
import type { paths } from '@generated/schema';
import { getTokens, setTokens, clearTokens } from '@/lib/tokenStorage';
import { getConfig } from '@/lib/runtimeConfig';

// Lazy-load BASE_URL
let BASE_URL: string;
function getBaseUrl(): string {
  if (!BASE_URL) {
    BASE_URL = getConfig().VITE_BACKEND_API;
  }
  return BASE_URL;
}

// Single-flight refresh lock
let refreshing: Promise<boolean> | null = null;

async function doRefresh(refreshToken: string): Promise<boolean> {
  if (!refreshToken) return false;
  if (refreshing) return refreshing;

  refreshing = (async () => {
    try {
      const res = await fetch(`${getBaseUrl()}/account/refresh`, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ refreshToken }),
      });

      if (!res.ok) {
        clearTokens();
        return false;
      }

      const data = await res.json();
      if (!data?.accessToken || !data?.refreshToken) {
        clearTokens();
        return false;
      }

      setTokens({ accessToken: data.accessToken, refreshToken: data.refreshToken });
      return true;
    } catch {
      clearTokens();
      return false;
    } finally {
      refreshing = null;
    }
  })();

  return refreshing;
}

const authMiddleware: Middleware = {
  async onRequest({ request }) {
    const { accessToken } = getTokens();
    if (accessToken) {
      request.headers.set('Authorization', `Bearer ${accessToken}`);
    }
    return request;
  },
  async onResponse({ request, response }) {
    if (response.status === 401) {
      const { refreshToken } = getTokens();
      if (!refreshToken) return response;

      const refreshed = await doRefresh(refreshToken);
      if (!refreshed) return response;

      // retry the original request with updated token
      const { accessToken } = getTokens();
      if (accessToken) request.headers.set('Authorization', `Bearer ${accessToken}`);
      // re-send the request using Fetch API
      const retried = await fetch(request);
      return retried;
    }
    return response;
  },
  async onError({ error }) {
    // wrap or rethrow; keep original behaviour by rethrowing
    throw error;
  },
};

export const api = createClient<paths>({
  baseUrl: '', // Will be set dynamically
  credentials: 'include',
});

// Update baseUrl after config is loaded
export function initializeApiClient() {
  api.baseUrl = getBaseUrl();
}

api.use(authMiddleware);
