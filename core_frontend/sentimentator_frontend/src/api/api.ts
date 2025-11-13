const API_URL = import.meta.env.VITE_BACKEND_API || 'http://localhost:8081/api';

export enum Endpoints {
  USERS = '/users',
  ANALYSES = '/analyses',
  AUTH_LOGIN = '/auth/login',
  AUTH_REGISTER = '/auth/register',
  AUTH_LOGOUT = '/auth/logout',
  AUTH_REFRESH = '/auth/refresh',
  // Add more endpoints as needed
}

export async function apiFetch(
  endpoint: Endpoints | string,
  options: RequestInit = {},
  token?: string,
) {
  const endpointStr = typeof endpoint === 'string' ? endpoint : endpoint;
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
    ...(token && { Authorization: `Bearer ${token}` }),
    ...options.headers,
  };

  const response = await fetch(`${API_URL}${endpointStr}`, {
    ...options,
    headers,
  });

  if (!response.ok) {
    throw new Error(`API error: ${response.status}`);
  }
  return response.json();
}

export async function apiPost(
  endpoint: Endpoints | string,
  data: unknown,
  token?: string,
  options: RequestInit = {},
) {
  return apiFetch(
    endpoint,
    {
      method: 'POST',
      body: JSON.stringify(data),
      ...options,
    },
    token,
  );
}
