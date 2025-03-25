import axios, { AxiosError, AxiosInstance, InternalAxiosRequestConfig, AxiosResponse } from "axios"
import { store } from "@/lib/store"
import { refreshToken as refreshTokenAction } from "@/lib/auth"
import { logout } from "@/lib/features/auth/authSlice"

export interface ValidationError {
  field: string
  message: string
}

export interface ApiError {
  code: string
  message: string
  details?: string
  validation_errors?: ValidationError[]
}

interface ErrorResponseData {
  code: string
  message: string
  details?: string
  validation_errors?: ValidationError[]
}

const BASE_URL =
  typeof window === "undefined"
    ? process.env.NEXT_PUBLIC_BACKEND_API_URL
    : process.env.NEXT_PUBLIC_BACKEND_API_URL?.replace("gateway", "localhost") || "http://localhost:3012"

export const axiosInstance: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
  withCredentials: false,
  timeout: 10000, // 10 second timeout
  timeoutErrorMessage: "Request timed out",
})

axiosInstance.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    const state = store.getState()
    const token = state.auth.accessToken

    if (token) {
      config.headers = config.headers || {}
      config.headers.Authorization = `Bearer ${token}`
    } else {
      console.warn("No token available for request:", {
        url: config.url,
        method: config.method,
        timestamp: new Date().toISOString(),
      })
    }

    return config
  },
  (error: AxiosError): Promise<AxiosError> => {
    console.error("Request interceptor error:", {
      message: error.message,
      code: error.code,
      timestamp: new Date().toISOString(),
    })
    return Promise.reject(error)
  }
)

axiosInstance.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => response,
  async (error: AxiosError<ErrorResponseData>): Promise<any> => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const state = store.getState()
        const refreshTokenStr = state.auth.refreshToken

        if (!refreshTokenStr) {
          store.dispatch(logout())
          return Promise.reject(error)
        }

        const { accessToken } = await refreshTokenAction(refreshTokenStr)

        if (originalRequest.headers) {
          originalRequest.headers.Authorization = `Bearer ${accessToken}`
        }

        return axiosInstance(originalRequest)
      } catch (refreshError) {
        store.dispatch(logout())
        return Promise.reject(refreshError)
      }
    }

    if (error.code === "ECONNABORTED" || !error.response) {
      console.error("Network error:", {
        message: error.message,
        code: error.code,
        timestamp: new Date().toISOString(),
      })
    }

    const apiError: ApiError = {
      code: error.response?.data?.code || "INTERNAL_ERROR",
      message: error.response?.data?.message || error.response?.statusText || error.message || "An unknown error occurred",
      details: error.response?.data?.details,
      validation_errors: error.response?.data?.validation_errors,
    }

    return Promise.reject(apiError)
  }
)

export const createCancelToken = () => {
  return axios.CancelToken.source()
}

export const setTimeout = (ms: number) => {
  axiosInstance.defaults.timeout = ms
}
