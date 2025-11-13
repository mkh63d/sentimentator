import createClient from 'openapi-fetch';
import type { paths } from '@generated/schema';

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8000/api';
let accessToken: string | null = localStorage.getItem('accessToken');
let refreshToken: string | null = localStorage.getItem('refreshToken');

export const api = createClient<paths>({
  baseUrl: BASE_URL,
  credentials: 'include',
  headers: accessToken ? { Authorization: `Bearer ${accessToken}` } : {},
});

export function setTokens(tokens: { accessToken: string; refreshToken: string }) {
  accessToken = tokens.accessToken;
  refreshToken = tokens.refreshToken;

  localStorage.setItem('accessToken', tokens.accessToken);
  localStorage.setItem('refreshToken', tokens.refreshToken);
}

export function clearTokens() {
  accessToken = null;
  refreshToken = null;

  localStorage.removeItem('accessToken');
  localStorage.removeItem('refreshToken');
}

export function getTokens() {
  return {
    accessToken,
    refreshToken,
  };
}
