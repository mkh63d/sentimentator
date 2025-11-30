import { api } from '@/lib/client';
import { setTokens, clearTokens } from '@/lib/tokenStorage';
import { getConfig } from '@/lib/runtimeConfig';

function isRegisterEnabled(): boolean {
  try {
    const cfg = getConfig();
    return cfg.VITE_REGISTER_ENABLED === 'true';
  } catch {
    return import.meta.env.VITE_REGISTER_ENABLED === 'true';
  }
}

export interface InfoRequest {
  newEmail?: string | null;
  newPassword?: string | null;
  oldPassword?: string | null;
}

export async function login(
  email: string,
  password: string,
  useCookies: boolean = false,
  useSessionCookies: boolean = false,
) {
  clearTokens();

  const { data: AccessTokenResponse } = await api.POST('/account/login', {
    params: {
      query: {
        useCookies,
        useSessionCookies,
      },
    },
    body: {
      email,
      password,
    },
  });

  if (!AccessTokenResponse && !useCookies) {
    throw new Error('Wrong email or password');
  }

  if (useCookies) {
    return { authenticatedVia: 'cookie' };
  }

  const tokenResponse = AccessTokenResponse;
  if (!tokenResponse || !tokenResponse.accessToken || !tokenResponse.refreshToken) {
    throw new Error('Connection error');
  }

  setTokens(tokenResponse);

  return {
    authenticatedVia: 'token',
    accessToken: tokenResponse.accessToken,
    refreshToken: tokenResponse.refreshToken,
  };
}

/** Logs out the current user by clearing stored tokens.
 */
export async function logout() {
  clearTokens();
}

/**
 * Registers a new user account.
 *
 * @param email - The user's email
 * @param password - The user's password
 */
export async function register(email: string, password: string) {
  if (!isRegisterEnabled()) {
    throw new Error('Registration is disabled');
  }

  const { data, error: fetchError } = await api.POST('/account/register', {
    body: {
      email,
      password,
    },
  });

  if (fetchError?.status === 400) {
    const errorsObj = fetchError.errors ?? {};
    const firstKey = Object.keys(errorsObj)[0];
    const errorMessage = firstKey ? errorsObj[firstKey][0] : 'Failed to register user (Http error)';
    throw new Error(errorMessage);
  } else if (fetchError) {
    throw new Error(fetchError.title ?? 'Failed to register user');
  }

  return data;
}

/**
 * Refreshes the access token using the provided refresh token.
 *
 * @param refreshToken - The refresh token to use for refreshing the access token
 * @returns The new access token
 */
export async function refreshAccount(refreshToken?: string | null) {
  // If a refresh token string is provided, send it in the body (token-based refresh).
  // Otherwise, rely on cookie-based refresh where the server reads httpOnly cookie.
  if (refreshToken) {
    const { data: accessTokenResponse } = await api.POST('/account/refresh', {
      body: { refreshToken },
    });
    return accessTokenResponse;
  }
  // If no refresh token is provided, attempt cookie-based refresh using fetch
  try {
    const baseUrl = getConfig().VITE_BACKEND_API;
    const res = await fetch(`${baseUrl}/account/refresh`, {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
    });

    if (!res.ok) {
      return null;
    }

    const data = await res.json();
    return data;
  } catch {
    return null;
  }
}

/**
 * Confirms the user's email address.
 *
 * @param userId
 * @param code
 * @param changedEmail
 * @returns
 */
export async function confirmEmail(userId: string, code: string, changedEmail?: string) {
  const { data, error: fetchError } = await api.GET('/account/confirmEmail', {
    params: {
      query: {
        userId,
        code,
        changedEmail,
      },
    },
  });

  if (fetchError) {
    throw new Error(fetchError ?? 'Failed to confirm email');
  }

  return data;
}

/**
 * Resends the confirmation email to the user.
 *
 * @param email
 * @returns
 */
export async function resendConfirmationEmail(email: string) {
  const { data, error: fetchError } = await api.POST('/account/resendConfirmationEmail', {
    body: {
      email,
    },
  });

  if (fetchError) {
    throw new Error(fetchError ?? 'Failed to resend confirmation email');
  }

  return data;
}

/**
 * Sends a password reset email to the user.
 *
 * @param email
 * @returns
 */
export async function forgotPassword(email: string) {
  try {
    const { error: fetchError } = await api.POST('/account/forgotPassword', {
      body: {
        email,
      },
    });
    if (fetchError?.status === 400) {
      const errorsObj = fetchError.errors ?? {};
      const firstKey = Object.keys(errorsObj)[0];
      const errorMessage = firstKey
        ? errorsObj[firstKey][0]
        : 'Failed to register user (Http error)';
      throw new Error(errorMessage);
    } else if (fetchError) {
      throw new Error(fetchError.title ?? 'Failed to send password reset email');
    }
  } catch (err) {
    console.error('Error in forgotPassword API call:', err);
    throw err;
  }
}

/**
 * Resets the user's password.
 *
 * @param email
 * @param resetCode
 * @param newPassword
 * @returns
 */
export async function resetPassword(email: string, resetCode: string, newPassword: string) {
  const { data, error: fetchError } = await api.POST('/account/resetPassword', {
    body: {
      email,
      resetCode,
      newPassword,
    },
  });
  if (fetchError?.status === 400) {
    throw new Error(fetchError.title ?? 'Failed to reset password (Http error)');
  } else if (fetchError) {
    throw new Error(fetchError.title ?? 'Failed to reset password');
  }
  return data;
}

/**
 * Fetches the account information for the logged-in user.
 *
 * @returns The account information
 */
export async function getAccountInfo() {
  const { data, error: fetchError } = await api.GET('/account/manage/info');
  if (fetchError?.status === 400) {
    throw new Error(fetchError.title ?? 'Failed to obtain account info (Http error)');
  } else if (fetchError?.status === 404) {
    throw new Error(fetchError.title ?? 'Account not found');
  } else if (fetchError) {
    throw new Error(fetchError.title ?? 'Failed to obtain account info');
  }
  return data;
}

/**
 * Updates the account information for the logged-in user.
 *
 * @param request
 * @returns
 */
export async function updateAccountInfo(request: InfoRequest) {
  const { data, error: fetchError } = await api.POST('/account/manage/info', {
    body: request,
  });
  if (fetchError?.status === 400) {
    throw new Error(fetchError.title ?? 'Failed to update account info (Http error)');
  } else if (fetchError?.status === 404) {
    throw new Error(fetchError.title ?? 'Account not found');
  } else if (fetchError) {
    throw new Error(fetchError.title ?? 'Failed to update account info');
  }
  return data;
}
