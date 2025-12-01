import { defineStore } from 'pinia';
import { ref, computed, watch } from 'vue';
import { login, logout, refreshAccount, register } from '@/api/account';
import { getTokens, setTokens, clearTokens } from '@/lib/tokenStorage';

// Pinia for exposing the tokens and authentication state to the app
export const useAuthStore = defineStore('auth', () => {
  const accessToken = ref<string | null>(getTokens().accessToken);
  const refreshToken = ref<string | null>(getTokens().refreshToken);
  const cookieAuthenticated = ref(false);
  const isAuthenticated = computed(() => !!accessToken.value || cookieAuthenticated.value);
  const loading = ref(false);
  const error = ref<string | null>(null);

  watch(
    [accessToken, refreshToken],
    ([access, refresh]) => {
      if (access && refresh) {
        setTokens({ accessToken: access, refreshToken: refresh });
      } else {
        // clear if either token is missing
        clearTokens();
      }
    },
    { immediate: true },
  );

  async function loginUser(
    email: string,
    password: string,
    useCookies = false,
    useSessionCookies = false,
  ) {
    loading.value = true;
    error.value = null;

    try {
      const result = await login(email, password, useCookies, useSessionCookies);

      if (result.authenticatedVia === 'token') {
        accessToken.value = result.accessToken ?? null;
        refreshToken.value = result.refreshToken ?? null;
      } else if (result.authenticatedVia === 'cookie') {
        // backend set httpOnly cookies; mark store as authenticated via cookie
        cookieAuthenticated.value = true;
      }

      return result;
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      error.value = message || 'Login failed';
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function registerUser(email: string, password: string) {
    try {
      await register(email, password);
    } catch (err) {
      error.value = 'Registration failed: ' + (err instanceof Error ? err.message : '');
      throw err;
    }
  }

  async function refreshAccountTokens() {
    loading.value = true;
    error.value = null;

    try {
      const result = await refreshAccount(refreshToken.value!);

      if (!result || !result.accessToken || !result.refreshToken) {
        throw new Error('No token data returned from API');
      }

      accessToken.value = result.accessToken;
      refreshToken.value = result.refreshToken;
      cookieAuthenticated.value = false;

      return result;
    } catch (err) {
      const message = err instanceof Error ? err.message : String(err);
      error.value = message || 'Token refresh failed';
      accessToken.value = null;
      refreshToken.value = null;
      cookieAuthenticated.value = false;
      throw err;
    } finally {
      loading.value = false;
    }
  }

  async function logoutUser() {
    await logout();

    accessToken.value = null;
    refreshToken.value = null;
    clearTokens();
    cookieAuthenticated.value = false;
  }

  return {
    accessToken,
    refreshToken,
    refreshAccountTokens,
    isAuthenticated,
    loading,
    error,
    loginUser,
    registerUser,
    logoutUser,
  };
});
