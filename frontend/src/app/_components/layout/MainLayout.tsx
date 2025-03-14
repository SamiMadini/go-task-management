"use client"

import { SidebarInset, SidebarProvider } from "@/components/ui/sidebar"
import { TopBar } from "@/app/_components/layout/top-bar"
import { AppSidebar } from "@/app/_components/layout/app-sidebar"

export default function MainLayout({ children, isExpanded }: { children: React.ReactNode; isExpanded: boolean }) {
  return (
    <SidebarProvider defaultOpen={isExpanded}>
      <AppSidebar />
      <SidebarInset>
        <TopBar />
        <main className="flex-1 p-4 md:p-6">{children}</main>
      </SidebarInset>
    </SidebarProvider>
  )
}
