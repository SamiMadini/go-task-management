import { NextResponse } from "next/server"
import type { NextRequest } from "next/server"

const protectedPaths = ["/tasks", "/profile"]
const authPaths = ["/auth/signin", "/auth/signup", "/auth/forgot-password", "/auth/reset-password"]

export function middleware(request: NextRequest) {
  const path = request.nextUrl.pathname

  if (path.startsWith("/api") || path.startsWith("/_next") || path.startsWith("/favicon.ico")) {
    return NextResponse.next()
  }

  if (protectedPaths.some((pp) => path.startsWith(pp))) {
    return NextResponse.next()
  }

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
