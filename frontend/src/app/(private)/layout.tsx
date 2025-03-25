import { cookies } from "next/headers"
import MainLayout from "@/app/_components/layout/MainLayout"
import { AuthProvider } from "@/app/_components/auth/AuthProvider"

export default function PrivateLayout({
  children,
}: Readonly<{
  children: React.ReactNode
}>) {
  const cookieStore = cookies()
  const isExpanded = cookieStore.get("app:sidebarExpanded") !== undefined ? cookieStore.get("app:sidebarExpanded")?.value === "true" : true

  return (
    <AuthProvider>
      <MainLayout isExpanded={isExpanded}>{children}</MainLayout>
    </AuthProvider>
  )
}
