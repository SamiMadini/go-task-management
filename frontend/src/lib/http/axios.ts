import axios from "axios"
import { store } from "@/lib/store"

const BASE_URL =
  typeof window === "undefined"
    ? process.env.NEXT_PUBLIC_BACKEND_API_URL
    : process.env.NEXT_PUBLIC_BACKEND_API_URL?.replace("gateway", "localhost") || "http://localhost:3012"

export const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: { "Content-Type": "application/json" },
  withCredentials: false,
})

axiosInstance.interceptors.request.use(
  (config) => {
    const state = store.getState()
    const token = state.auth.accessToken

    if (token) {
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

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    if (error.response?.status === 401) {
      console.log("Received 401 error, token might be invalid or expired")
    }
    return Promise.reject(error)
  }
)
