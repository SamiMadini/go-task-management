import { store } from "./store"
import { setCredentials, setError, setLoading, logout, User } from "./features/auth/authSlice"
import { axiosInstance } from "./http/axios"
import axios from "axios"
import { persistor } from "./store"
import { ApiError, ValidationError } from "./http/axios"

export interface SignInCredentials {
  email: string
  password: string
}

export interface SignUpCredentials extends SignInCredentials {
  handle: string
}

export interface ResetPasswordCredentials {
  token: string
  password: string
}

export interface AuthResponse {
  data: {
    user: User
    token: {
      access_token: string
      refresh_token: string
      expires_in: number
    }
  }
}

export class AuthError extends Error {
  constructor(message: string, public code?: string, public details?: string, public validationErrors?: ValidationError[]) {
    super(message)
    this.name = "AuthError"
  }
}

const handleAuthError = (error: unknown): never => {
  if (axios.isAxiosError(error)) {
    const errorData = error.response?.data as ApiError
    throw new AuthError(errorData?.message || error.message, errorData?.code, errorData?.details, errorData?.validation_errors)
  }
  if ((error as ApiError)?.code) {
    const apiError = error as ApiError
    throw new AuthError(apiError.message, apiError.code, apiError.details, apiError.validation_errors)
  }
  throw new AuthError(error instanceof Error ? error.message : "An unknown error occurred")
}

const dispatchAuthError = (error: unknown) => {
  let message = "An unknown error occurred"
  if (error instanceof AuthError) {
    message = error.validationErrors?.length ? error.validationErrors.map((e) => `${e.field}: ${e.message}`).join(", ") : error.message
  } else if (error instanceof Error) {
    message = error.message
  }
  store.dispatch(setError(message))
  return handleAuthError(error)
}

const withLoading = async <T>(operation: () => Promise<T>): Promise<T> => {
  store.dispatch(setLoading(true))
  try {
    return await operation()
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function signIn(credentials: SignInCredentials): Promise<void> {
  return withLoading(async () => {
    try {
      const response = await axiosInstance.post<AuthResponse>("/api/v1/auth/signin", credentials)
      const { token, user } = response.data.data

      if (!token.access_token || !token.refresh_token) {
        throw new AuthError("Invalid response: missing tokens")
      }

      store.dispatch(
        setCredentials({
          user,
          accessToken: token.access_token,
          refreshToken: token.refresh_token,
        })
      )

      await Promise.all([new Promise((resolve) => setTimeout(resolve, 0)), persistor.flush()])

      if (process.env.NODE_ENV === "development") {
        const state = store.getState()
        console.log("Auth state after signin:", {
          hasUser: !!state.auth.user,
          hasTokens: !!state.auth.accessToken && !!state.auth.refreshToken,
        })
      }
    } catch (error) {
      return dispatchAuthError(error)
    }
  })
}

export async function signUp(credentials: SignUpCredentials): Promise<void> {
  return withLoading(async () => {
    try {
      const response = await axiosInstance.post<AuthResponse>("/api/v1/auth/signup", credentials)
      const { token, user } = response.data.data

      store.dispatch(
        setCredentials({
          user,
          accessToken: token.access_token,
          refreshToken: token.refresh_token,
        })
      )

      await persistor.flush()
    } catch (error) {
      return dispatchAuthError(error)
    }
  })
}

export async function signOut(): Promise<void> {
  return withLoading(async () => {
    try {
      // await axiosInstance.post("/api/v1/auth/signout")
      store.dispatch(logout())
      await persistor.purge()
    } catch (error) {
      return dispatchAuthError(error)
    }
  })
}

export async function forgotPassword(email: string): Promise<void> {
  return withLoading(async () => {
    try {
      await axiosInstance.post("/api/v1/auth/forgot-password", { email })
    } catch (error) {
      return dispatchAuthError(error)
    }
  })
}

export async function resetPassword(credentials: ResetPasswordCredentials): Promise<void> {
  return withLoading(async () => {
    try {
      await axiosInstance.post("/api/v1/auth/reset-password", credentials)
    } catch (error) {
      return dispatchAuthError(error)
    }
  })
}

export async function refreshToken(token: string): Promise<{ accessToken: string }> {
  try {
    const response = await axiosInstance.post<{ data: AuthResponse["data"]["token"] }>("/api/v1/auth/refresh-token", {
      refresh_token: token,
    })

    if (!response.data.data.access_token) {
      throw new AuthError("No access token received in refresh response")
    }

    const currentUser = store.getState().auth.user
    if (!currentUser) {
      throw new AuthError("No user found in state during token refresh")
    }

    store.dispatch(
      setCredentials({
        user: currentUser,
        accessToken: response.data.data.access_token,
        refreshToken: response.data.data.refresh_token,
      })
    )

    return { accessToken: response.data.data.access_token }
  } catch (error) {
    console.error("RefreshToken::Token refresh failed:", error)
    if (axios.isAxiosError(error)) {
      console.error("RefreshToken::Response details:", {
        status: error.response?.status,
        data: error.response?.data,
      })
    }
    return dispatchAuthError(error)
  }
}
