import { store } from "./store"
import { setCredentials, setError, setLoading, logout, User } from "./features/auth/authSlice"
import { axiosInstance } from "./http/axios"
import axios from "axios"
import { persistor } from "./store"

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
  user: User
  access_token: string
  refresh_token: string
  expires_in: number
}

// Helper function to get headers with optional auth token
// const getHeaders = (includeAuth: boolean = false) => {
//   const headers: Record<string, string> = {
//     "Content-Type": "application/json",
//   }

//   if (includeAuth) {
//     const state = store.getState()
//     const token = state.auth.accessToken
//     if (token) {
//       headers.Authorization = `Bearer ${token}`
//     }
//   }

//   return headers
// }

export async function signIn(credentials: SignInCredentials): Promise<void> {
  try {
    store.dispatch(setLoading(true))
    console.log("Attempting to sign in...")
    const response = await axiosInstance.post<AuthResponse>("/api/v1/auth/signin", credentials)

    if (!response.data.access_token || !response.data.refresh_token) {
      throw new Error("Invalid response: missing tokens")
    }

    console.log(">> Sign in successful, received tokens:", {
      hasAccessToken: !!response.data.access_token,
      hasRefreshToken: !!response.data.refresh_token,
      user: response.data.user,
    })

    // Update the store with tokens
    store.dispatch(
      setCredentials({
        user: response.data.user,
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
      })
    )

    // Wait for the store to be updated
    await new Promise((resolve) => setTimeout(resolve, 0))

    // Verify the store was updated
    const state = store.getState()
    console.log("Redux store state after update:", {
      user: state.auth.user,
      hasAccessToken: !!state.auth.accessToken,
      hasRefreshToken: !!state.auth.refreshToken,
      accessToken: state.auth.accessToken?.substring(0, 10) + "...",
      refreshToken: state.auth.refreshToken?.substring(0, 10) + "...",
    })

    // Verify localStorage
    const persistedState = localStorage.getItem("persist:root")
    console.log("LocalStorage state:", persistedState ? JSON.parse(persistedState) : null)

    // Force persist the state
    await persistor.flush()
  } catch (error) {
    console.error("Sign in failed:", error)
    store.dispatch(setError(error instanceof Error ? error.message : "An error occurred"))
    throw error
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function signUp(credentials: SignUpCredentials): Promise<void> {
  try {
    store.dispatch(setLoading(true))
    const response = await axiosInstance.post<AuthResponse>("/api/v1/auth/signup", credentials)

    store.dispatch(
      setCredentials({
        user: response.data.user,
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
      })
    )
  } catch (error) {
    store.dispatch(setError(error instanceof Error ? error.message : "An error occurred"))
    throw error
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function signOut(): Promise<void> {
  try {
    store.dispatch(setLoading(true))
    // await axiosInstance.post("/api/v1/auth/signout")
    store.dispatch(logout())
    await persistor.purge()
  } catch (error) {
    store.dispatch(setError(error instanceof Error ? error.message : "An error occurred"))
    throw error
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function forgotPassword(email: string): Promise<void> {
  try {
    store.dispatch(setLoading(true))
    await axiosInstance.post("/api/v1/auth/forgot-password", { email })
  } catch (error) {
    store.dispatch(setError(error instanceof Error ? error.message : "An error occurred"))
    throw error
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function resetPassword(credentials: ResetPasswordCredentials): Promise<void> {
  try {
    store.dispatch(setLoading(true))
    await axiosInstance.post("/api/v1/auth/reset-password", credentials)
  } catch (error) {
    store.dispatch(setError(error instanceof Error ? error.message : "An error occurred"))
    throw error
  } finally {
    store.dispatch(setLoading(false))
  }
}

export async function refreshToken(token: string): Promise<{ accessToken: string }> {
  try {
    console.log("Attempting to refresh token...")
    const response = await axiosInstance.post<{ access_token: string; refresh_token: string; expires_in: number }>(
      "/api/v1/auth/refresh-token",
      {
        refresh_token: token,
      }
    )

    if (!response.data.access_token) {
      throw new Error("No access token received in refresh response")
    }

    // Update both tokens in the store
    store.dispatch(
      setCredentials({
        user: store.getState().auth.user!,
        accessToken: response.data.access_token,
        refreshToken: response.data.refresh_token,
      })
    )

    console.log("Token refresh successful")
    return { accessToken: response.data.access_token }
  } catch (error) {
    console.error("Token refresh failed:", error)
    if (axios.isAxiosError(error)) {
      console.error("Response data:", error.response?.data)
      console.error("Response status:", error.response?.status)
    }
    throw error
  }
}
