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

async function doRefresh(refreshToken?: string | null): Promise<boolean> {
  // If refreshToken is provided, prefer token-based refresh.
  // Otherwise, attempt cookie-based refresh (server should use httpOnly cookie).
  if (refreshing) return refreshing;
  refreshing = (async () => {
    try {
      let res: Response;

      if (refreshToken) {
        res = await fetch(`${getBaseUrl()}/account/refresh`, {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refreshToken }),
        });
      } else {
        // Attempt cookie-based refresh: no body, server should read httpOnly cookie
        res = await fetch(`${getBaseUrl()}/account/refresh`, {
          method: 'POST',
          credentials: 'include',
          headers: { 'Content-Type': 'application/json' },
        });
      }

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
      // attempt refresh even when no refresh token is stored (cookie-based refresh)
      const refreshed = await doRefresh(refreshToken ?? null);
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

let _api: ReturnType<typeof createClient<paths>> | null = null;

// Initialize and return API client (creates once after config loads)
export function initializeApiClient() {
  if (!_api) {
    _api = createClient<paths>({
      baseUrl: getBaseUrl(),
      credentials: 'include',
    });
    _api.use(authMiddleware);
  }
  return _api;
}

// Export api getter that ensures client is initialized
export const api = new Proxy({} as ReturnType<typeof createClient<paths>>, {
  get(target, prop) {
    const client = _api || initializeApiClient();
    return client[prop as keyof typeof client];
  },
});
