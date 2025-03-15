"use client"

import { useEffect, useState } from "react"
import { Bell } from "lucide-react"

import { Button } from "@/components/ui/button"
import { axiosInstance } from "@/lib/http/axios"
import { SheetRight } from "@/app/_components/common/sheet-right"
import { GetOneNotificationInterface } from "@/app/domain/notification/interfaces.notification"
import { NotificationsTabsComponent } from "@/app/_components/notifications/notifications-tabs.component"

export function NotificationsComponent() {
  const [notificationOpen, setNotificationOpen] = useState(false)
  const [notifications, setNotifications] = useState<GetOneNotificationInterface[]>([])

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await axiosInstance.get("/api/notifications")
        if (response.data) {
          setNotifications(response.data)
        }
      } catch (error) {
        console.error("Error fetching notifications:", error)
      }
    }

    const timer = setTimeout(() => {
      fetchNotifications()
    }, 5000)

    return () => clearTimeout(timer)
  }, [])

  const handleReadNotification = async (notificationId: number, isRead: boolean) => {
    try {
      await axiosInstance.post(`/api/notifications/${notificationId}/read`, {
        isRead: isRead,
      })
      setNotifications(
        notifications.map((notification) => {
          if (notification.id === notificationId) {
            return {
              ...notification,
              is_read: isRead,
            }
          }
          return notification
        })
      )
    } catch (error) {
      console.error("Error reading notification:", error)
    }
  }

  return (
    <>
      <Button variant="ghost" size="icon" className="text-muted-foreground relative" onClick={() => setNotificationOpen(true)}>
        <Bell className="h-5 w-5" />
        {notifications.some((notification) => !notification.is_read) && (
          <span className="absolute top-1 right-1 h-3 w-3 rounded-full bg-red-500 animate-pulse" />
        )}
        <span className="sr-only">Notifications</span>
      </Button>

      <SheetRight
        title="Notifications"
        description=""
        overlay={false}
        isModal={true}
        open={notificationOpen}
        setIsOpen={setNotificationOpen}
      >
        <NotificationsTabsComponent
          notifications={notifications}
          onRead={(notificationId, isRead) => handleReadNotification(notificationId, isRead)}
        />
      </SheetRight>

      {notifications.some((notification) => !notification.is_read) && (
        <div className="fixed top-20 right-6 items-center justify-center flex flex-col transition-all ease-in-out animate-in fade-in duration-1000 fill-mode-forwards">
          <img src="/arrow-rough-drawig-top-right.png" alt="New notifications" className="w-20 h-20 -rotate-45" />
          <p className="pt-5 text-foreground animate-pulse text-red-500 font-bold">New notifications</p>
        </div>
      )}
    </>
  )
}
