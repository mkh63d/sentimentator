import { apiFetch, apiPost, Endpoints } from './api';

export interface LoginCredentials {
  username: string;
  password: string;
}

export interface RegisterData {
  username: string;
  email: string;
  password: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: {
    id: string;
    username: string;
    email: string;
  };
}

export const authApi = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    return apiPost(Endpoints.AUTH_LOGIN, credentials);
  },

  async register(data: RegisterData): Promise<AuthResponse> {
    return apiPost(Endpoints.AUTH_REGISTER, data);
  },

  async logout(token: string): Promise<void> {
    await apiFetch(Endpoints.AUTH_LOGOUT, { method: 'POST' }, token);
  },

  async refreshToken(refreshToken: string): Promise<AuthResponse> {
    return apiPost(Endpoints.AUTH_REFRESH, { refresh_token: refreshToken });
  },
};
