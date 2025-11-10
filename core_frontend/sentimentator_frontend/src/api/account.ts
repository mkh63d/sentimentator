import { api, setTokens, clearTokens } from '@/lib/client';

export interface InfoRequest {
  newEmail?: string | null;
  newPassword?: string | null;
  oldPassword?: string | null;
}

/**
 * Logs in a user with the provided email and password.
 *
 * @param email
 * @param password
 * @param useCookies
 * @param useSessionCookies
 * @returns
 */
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

  if (useCookies) {
    return { authenticatedVia: 'cookie' };
  }

  const tokenResponse = AccessTokenResponse;
  if (!tokenResponse || !tokenResponse.accessToken || !tokenResponse.refreshToken) {
    throw new Error('No token data returned from API');
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
  const { data, error: fetchError } = await api.POST('/account/register', {
    body: {
      email,
      password,
    },
  });

  if (fetchError) {
    throw new Error(fetchError.detail ?? 'Failed to register user');
  }

  return data;
}

/**
 * Refreshes the access token using the provided refresh token.
 *
 * @param refreshToken - The refresh token to use for refreshing the access token
 * @returns The new access token
 */
export async function refreshAccount(refreshToken: string) {
  const { data: accessTokenResponse, error: fetchError } = await api.POST('/account/refresh', {
    body: {
      refreshToken,
    },
  });

  if (fetchError) {
    throw new Error(fetchError ?? 'Failed to refresh access token');
  }

  return accessTokenResponse;
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
  if (fetchError) {
    throw new Error(fetchError.detail ?? 'Failed to reset password');
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
  if (fetchError) {
    throw new Error(fetchError.detail ?? 'Failed to fetch account info');
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
  if (fetchError) {
    throw new Error(fetchError.detail ?? 'Failed to update account info');
  }
  return data;
}
