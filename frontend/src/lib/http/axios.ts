import axios from "axios"

const BASE_URL = process.env.NEXT_PUBLIC_BACKEND_API_URL

const getBaseUrl = () => {
  if (typeof window !== "undefined") {
    if (BASE_URL && BASE_URL.includes("gateway")) {
      return BASE_URL.replace("http://gateway", "http://localhost")
    }
  } else {
    if (BASE_URL && BASE_URL.includes("localhost")) {
      return BASE_URL.replace("http://localhost", "http://gateway")
    }
  }
  return BASE_URL
}

export const axiosInstance = axios.create({
  baseURL: getBaseUrl(),
  headers: { "Content-Type": "application/json" },
})
