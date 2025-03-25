import { useEffect } from "react"
import { useRouter, usePathname } from "next/navigation"
import { useAuth } from "../hooks"

export function useProtectedRoute() {
  const router = useRouter()
  const pathname = usePathname()
  const { isAuthenticated, isLoading } = useAuth()

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      const returnUrl = encodeURIComponent(pathname)
      router.push(`/auth/signin?returnUrl=${returnUrl}`)
    }
  }, [isAuthenticated, isLoading, router, pathname])

  return { isLoading, isAuthenticated }
}
