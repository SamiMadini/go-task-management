"use client"

import { useEffect } from "react"
import { usePathname } from "next/navigation"
import { useAuth } from "@/lib/hooks"
// import { refreshToken } from "@/lib/auth"
import { store } from "@/lib/store"
// import { setAccessToken, logout } from "@/lib/features/auth/authSlice"

// const REFRESH_INTERVAL = 4 * 60 * 1000 // 4 minutes
// const TOKEN_EXPIRY_BUFFER = 5 * 60 * 1000 // 5 minutes

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { accessToken } = useAuth()
  const pathname = usePathname()

  // useEffect(() => {
  //   let refreshInterval: NodeJS.Timeout

  //   const setupTokenRefresh = async () => {
  //     if (!storedRefreshToken) return

  //     try {
  //       const response = await refreshToken(storedRefreshToken)
  //       store.dispatch(setAccessToken(response.accessToken))

  //       // Set up periodic refresh
  //       refreshInterval = setInterval(async () => {
  //         try {
  //           const response = await refreshToken(storedRefreshToken)
  //           store.dispatch(setAccessToken(response.accessToken))
  //         } catch (error) {
  //           console.error("Failed to refresh token:", error)
  //           store.dispatch(logout())
  //           clearInterval(refreshInterval)
  //         }
  //       }, REFRESH_INTERVAL)
  //     } catch (error) {
  //       console.error("Failed to refresh token:", error)
  //       store.dispatch(logout())
  //     }
  //   }

  //   // Only set up token refresh on initial mount or when refresh token changes
  //   if (isInitialMount.current) {
  //     setupTokenRefresh()
  //     isInitialMount.current = false
  //   }

  //   // Cleanup
  //   return () => {
  //     if (refreshInterval) {
  //       clearInterval(refreshInterval)
  //     }
  //   }
  // }, [storedRefreshToken])

  // // Handle token expiry
  // useEffect(() => {
  //   if (!accessToken || isInitialMount.current) return

  //   try {
  //     // Decode JWT to get expiry time
  //     const payload = JSON.parse(atob(accessToken.split(".")[1]))
  //     const expiryTime = payload.exp * 1000 // Convert to milliseconds
  //     const timeUntilExpiry = expiryTime - Date.now()

  //     // If token is about to expire, try to refresh it
  //     if (timeUntilExpiry < TOKEN_EXPIRY_BUFFER && storedRefreshToken) {
  //       refreshToken(storedRefreshToken)
  //         .then((response) => {
  //           store.dispatch(setAccessToken(response.accessToken))
  //         })
  //         .catch(() => {
  //           store.dispatch(logout())
  //         })
  //     }
  //   } catch (error) {
  //     console.error("Failed to process token:", error)
  //   }
  // }, [accessToken, storedRefreshToken])

  return <>{children}</>
}
