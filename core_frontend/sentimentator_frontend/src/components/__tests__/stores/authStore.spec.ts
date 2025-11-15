import { setActivePinia, createPinia } from 'pinia';
import { describe, it, expect, beforeEach, vi } from 'vitest';

// Mock the account API and client token helpers
const loginMock = vi.fn();
const registerMock = vi.fn();
const refreshAccountMock = vi.fn();
const logoutMock = vi.fn();
vi.mock('@/api/account', () => ({
  login: (...args: unknown[]) => (loginMock as unknown as (...a: unknown[]) => unknown)(...args),
  register: (...args: unknown[]) =>
    (registerMock as unknown as (...a: unknown[]) => unknown)(...args),
  refreshAccount: (...args: unknown[]) =>
    (refreshAccountMock as unknown as (...a: unknown[]) => unknown)(...args),
  logout: (...args: unknown[]) => (logoutMock as unknown as (...a: unknown[]) => unknown)(...args),
}));

const getTokensMock = vi.fn();
const setTokensMock = vi.fn();
vi.mock('@/lib/client', () => ({
  getTokens: () => getTokensMock(),
  setTokens: (...args: unknown[]) =>
    (setTokensMock as unknown as (...a: unknown[]) => unknown)(...args),
}));

describe('authStore', () => {
  beforeEach(() => {
    setActivePinia(createPinia());
    loginMock.mockReset();
    registerMock.mockReset();
    refreshAccountMock.mockReset();
    logoutMock.mockReset();
    getTokensMock.mockReset();
    setTokensMock.mockReset();
    // make sure localStorage is clean for logout tests
    localStorage.clear();
  });

  it('initializes tokens from getTokens and calls setTokens via watcher', async () => {
    getTokensMock.mockReturnValue({ accessToken: 'a1', refreshToken: 'r1' });
    // lazy import store to ensure mocks are active
    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    expect(store.accessToken).toBeDefined();
    expect(store.accessToken).toBe('a1');
    expect(store.refreshToken).toBe('r1');
    expect(store.isAuthenticated).toBeTruthy();
    // Watcher is immediate and should call setTokens
    expect(setTokensMock).toHaveBeenCalledWith({ accessToken: 'a1', refreshToken: 'r1' });
  });

  it('loginUser sets tokens on successful token-auth response', async () => {
    getTokensMock.mockReturnValue({ accessToken: null, refreshToken: null });
    loginMock.mockResolvedValue({
      authenticatedVia: 'token',
      accessToken: 'newA',
      refreshToken: 'newR',
    });

    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    await store.loginUser('e', 'p');

    expect(store.accessToken).toBe('newA');
    expect(store.refreshToken).toBe('newR');
    expect(store.isAuthenticated).toBeTruthy();
    expect(setTokensMock).toHaveBeenCalled();
  });

  it('loginUser throws and sets error on failure', async () => {
    getTokensMock.mockReturnValue({ accessToken: null, refreshToken: null });
    loginMock.mockRejectedValue(new Error('bad creds'));

    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    let threw = false;
    try {
      await store.loginUser('e', 'p');
    } catch {
      threw = true;
    }
    expect(threw).toBe(true);
    expect(store.error).toBeDefined();
    expect(store.loading).toBe(false);
  });

  it('registerUser calls register then loginUser', async () => {
    getTokensMock.mockReturnValue({ accessToken: null, refreshToken: null });
    registerMock.mockResolvedValue(undefined);
    loginMock.mockResolvedValue({
      authenticatedVia: 'token',
      accessToken: 'regA',
      refreshToken: 'regR',
    });

    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    await store.registerUser('e', 'p');

    expect(registerMock).toHaveBeenCalled();
    expect(loginMock).toHaveBeenCalled();
    expect(store.accessToken).toBe('regA');
  });

  it('refreshAccountTokens updates tokens on success', async () => {
    getTokensMock.mockReturnValue({ accessToken: 'a', refreshToken: 'r' });
    refreshAccountMock.mockResolvedValue({ accessToken: 'a2', refreshToken: 'r2' });

    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    const res = await store.refreshAccountTokens();
    expect(res).toBeDefined();
    expect(store.accessToken).toBe('a2');
    expect(store.refreshToken).toBe('r2');
  });

  it('logoutUser clears tokens and localStorage', async () => {
    getTokensMock.mockReturnValue({ accessToken: 'a', refreshToken: 'r' });
    logoutMock.mockResolvedValue(undefined);

    // put tokens into mocked localStorage to test removal (vitest.setup provides a fresh mock)
    globalThis.localStorage.setItem('accessToken', 'a');
    globalThis.localStorage.setItem('refreshToken', 'r');

    const { useAuthStore } = await import('@/stores/authStore');
    const store = useAuthStore();

    await store.logoutUser();

    expect(store.accessToken).toBeNull();
    expect(store.refreshToken).toBeNull();
    expect(localStorage.getItem('accessToken')).toBeNull();
    expect(localStorage.getItem('refreshToken')).toBeNull();
    // vitest.setup will reset mocked localStorage between tests
  });
});
