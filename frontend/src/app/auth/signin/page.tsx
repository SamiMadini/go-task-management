"use client"

import { useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import Link from "next/link"
import { signIn } from "@/lib/auth"
import { useAuth } from "@/lib/hooks"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { persistor } from "@/lib/store"

export default function SignInPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { isLoading, error } = useAuth()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      await signIn({ email, password })
      // Wait for the store to be persisted
      await persistor.flush()
      const returnUrl = searchParams.get("returnUrl")
      router.push(returnUrl ? decodeURIComponent(returnUrl) : "/tasks")
    } catch (err) {
      // Error is handled by the Redux store
      console.error("Sign in failed:", err)
    }
  }

  return (
    <div className="container flex items-center justify-center min-h-screen py-12">
      <Card className="w-full max-w-md">
        <CardHeader>
          <CardTitle>Sign In</CardTitle>
          <CardDescription>Enter your credentials to access your account</CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            {error && <div className="p-3 text-sm text-red-500 bg-red-50 rounded-md">{error}</div>}
            <div className="space-y-2">
              <Label htmlFor="email">Email</Label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="Enter your email"
                required
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="Enter your password"
                required
              />
            </div>
          </CardContent>
          <CardFooter className="flex flex-col space-y-4">
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Signing in..." : "Sign In"}
            </Button>
            <div className="text-sm text-center text-muted-foreground">
              Don&#39;t have an account?{" "}
              <Link href="/auth/signup" className="text-primary hover:underline">
                Sign up
              </Link>
              <span className="mx-2">•</span>
              <Link href="/auth/forgot-password" className="text-primary hover:underline">
                Forgot password?
              </Link>
            </div>
          </CardFooter>
        </form>
      </Card>
    </div>
  )
}
