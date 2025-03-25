import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

// Add paths that should be accessible only to authenticated users
const protectedPaths = ["/tasks", "/profile"]

// Add paths that should be accessible only to non-authenticated users
const authPaths = ["/auth/signin", "/auth/signup", "/auth/forgot-password", "/auth/reset-password"]

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname

  // Skip middleware for API routes and static files
  if (path.startsWith("/api") || path.startsWith("/_next") || path.startsWith("/favicon.ico")) {
    return NextResponse.next()
  }

  // For protected paths, let the client handle authentication
  if (protectedPaths.some((pp) => path.startsWith(pp))) {
    return NextResponse.next()
  }

  // For auth paths, let the client handle authentication
  if (authPaths.some((ap) => path.startsWith(ap))) {
    return NextResponse.next()
  }

  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    "/((?!api|_next/static|_next/image|favicon.ico).*)",
  ],
}
