"use client"

import { HelpCircle } from "lucide-react"
import { Button } from "@/components/ui/button"
import { SidebarTrigger, useSidebar } from "@/components/ui/sidebar"
import { NotificationsComponent } from "@/app/_components/notifications/notifications.component"
import { AvatarMenuComponent } from "@/app/_components/layout/avatar-menu.component"

export function TopBar() {
  const { open } = useSidebar()

  return (
    <header className="sticky top-0 z-10 flex h-16 items-center gap-4 border-b bg-background px-4 md:px-6">
      <SidebarTrigger
        onClick={() => {
          document.cookie = `app:sidebarExpanded=${!open}; path=/; max-age=31536000; SameSite=Lax`
        }}
      />

      <div className="hidden w-full max-w-sm md:flex">
        <div className="relative flex-1">
          {/* <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input type="search" placeholder="Search..." className="w-full rounded-lg bg-background pl-8 md:w-[240px] lg:w-[280px]" /> */}
        </div>
      </div>

      <div className="ml-auto flex items-center gap-2">
        <Button variant="ghost" size="icon" className="text-muted-foreground">
          <HelpCircle className="h-5 w-5" />
          <span className="sr-only">Help</span>
        </Button>

        <NotificationsComponent />

        <AvatarMenuComponent />
      </div>
    </header>
  )
}
