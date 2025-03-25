import axios from "axios"
import { store } from "@/lib/store"

// When running in Docker, the backend URL will be set via environment variable
// When running in browser, we need to use localhost
const BASE_URL =
  typeof window === "undefined"
    ? process.env.NEXT_PUBLIC_BACKEND_API_URL
    : process.env.NEXT_PUBLIC_BACKEND_API_URL?.replace("gateway", "localhost") || "http://localhost:3012"

export const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
  withCredentials: false, // We don't need credentials since we're using token-based auth
})

// Add a request interceptor to add the bearer token
axiosInstance.interceptors.request.use(
  (config) => {
    const state = store.getState()
    const token = state.auth.accessToken

    if (token) {
      console.log("Adding bearer token to request:", {
        url: config.url,
        hasToken: !!token,
        tokenPreview: token.substring(0, 10) + "...",
      })
      config.headers.Authorization = `Bearer ${token}`
    } else {
      console.log("No token available for request:", {
        url: config.url,
      })
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Add a response interceptor to handle 401 errors
axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      console.log("Received 401 error, token might be invalid or expired")
    }
    return Promise.reject(error)
  }
)
