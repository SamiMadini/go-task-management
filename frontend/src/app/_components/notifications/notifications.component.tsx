"use client"

import { useEffect, useState } from "react"
import { Bell } from "lucide-react"

import { Button } from "@/components/ui/button"
import { axiosInstance, ApiError } from "@/lib/http/axios"
import { SheetRight } from "@/app/_components/common/sheet-right"
import { GetOneNotificationInterface } from "@/app/domain/notification/interfaces.notification"
import { NotificationsTabsComponent } from "@/app/_components/notifications/notifications-tabs.component"
import { toast } from "sonner"

interface NotificationsResponse {
  data: {
    in_app_notifications: GetOneNotificationInterface[]
  }
}

export function NotificationsComponent() {
  const [notificationOpen, setNotificationOpen] = useState(false)
  const [notifications, setNotifications] = useState<GetOneNotificationInterface[]>([])

  const handleError = (error: unknown, action: string) => {
    const apiError = error as ApiError
    const message = apiError.message || (error instanceof Error ? error.message : "An unknown error occurred")
    const details = apiError.details

    toast.error(`Error ${action}`, {
      description: details || message,
    })
    console.error(`Error ${action}:`, error)
  }

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await axiosInstance.get<NotificationsResponse>("/api/v1/notifications")
        if (response.data?.data?.in_app_notifications) {
          setNotifications(response.data.data.in_app_notifications)
        }
      } catch (error) {
        handleError(error, "fetching notifications")
      }
    }

    const timer = setTimeout(() => {
      fetchNotifications()
    }, 5000)

    return () => clearTimeout(timer)
  }, [notifications])

  const handleReadNotification = async (notificationId: number, isRead: boolean) => {
    try {
      const response = await axiosInstance.post<{ data: { success: boolean } }>(`/api/v1/notifications/${notificationId}/read`, {
        is_read: isRead,
      })

      if (!response.data?.data?.success) {
        throw new Error("Failed to update notification")
      }

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
      toast.success(`Notification marked as ${isRead ? "read" : "unread"}`)
    } catch (error) {
      handleError(error, "updating notification")
    }
  }

  const handleDeleteNotification = async (notificationId: number) => {
    try {
      const response = await axiosInstance.delete<{ data: { success: boolean } }>(`/api/v1/notifications/${notificationId}`)

      if (!response.data?.data?.success) {
        throw new Error("Failed to delete notification")
      }

      setNotifications(notifications.filter((notification) => notification.id !== notificationId))
      toast.success("Notification deleted successfully")
    } catch (error) {
      handleError(error, "deleting notification")
    }
  }

  return (
    <>
      <Button variant="ghost" size="icon" className="text-muted-foreground relative" onClick={() => setNotificationOpen(true)}>
        <Bell className="h-5 w-5" />
        {notifications?.some((notification) => !notification.is_read) && (
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
          onDelete={(notificationId) => handleDeleteNotification(notificationId)}
        />
      </SheetRight>

      {/* {notifications?.some((notification) => !notification.is_read) && (
        <div className="fixed top-20 right-6 items-center justify-center flex flex-col transition-all ease-in-out animate-in fade-in duration-1000 fill-mode-forwards">
          <p className="pb-5 text-foreground animate-pulse text-red-500 font-bold">New notifications</p>
          <img src="/arrow-rough-drawig-top-right.png" alt="New notifications" className="w-20 h-20 -rotate-45" />
        </div>
      )} */}
    </>
  )
}
