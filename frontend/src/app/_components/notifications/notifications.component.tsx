"use client"

import { useEffect, useState } from "react"
import { Bell, CheckCircle2, CheckIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { axiosInstance } from "@/lib/http/axios"
import { SheetRight } from "@/app/_components/common/sheet-right"

interface Notification {
  id: number
  title: string
  description: string
  is_read: boolean
  created_at: string
  updated_at: string
}

export function NotificationsComponent() {
  const [loading, setLoading] = useState(false)
  const [notificationOpen, setNotificationOpen] = useState(false)
  const [notifications, setNotifications] = useState<Notification[]>([])

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        setLoading(true)
        const response = await axiosInstance.get("/api/notifications")
        if (response.data) {
          setNotifications(response.data)
        }
      } catch (error) {
        console.error("Error fetching notifications:", error)
      } finally {
        setLoading(false)
      }
    }

    const timer = setTimeout(() => {
      fetchNotifications()
    }, 5000)

    return () => clearTimeout(timer)
  }, [])

  const handleNotificationRead = async (notificationId: number) => {
    try {
      await axiosInstance.post(`/api/notifications/${notificationId}/read`)
      setNotifications(notifications.filter((notification) => notification.id !== notificationId))
    } catch (error) {
      console.error("Error reading notification:", error)
    } finally {
      setNotifications(notifications.filter((notification) => notification.id !== notificationId))
    }
  }

  const renderNotifications = (filteredNotifications: Notification[]) => {
    return (
      <div className="mt-6 space-y-6">
        {filteredNotifications.map((notification) => (
          <div key={notification.id} className="rounded-lg border p-4 relative">
            {!notification.is_read && <span className="absolute -top-1 -left-1 h-3 w-3 rounded-full bg-red-500 animate-pulse" />}
            <div className="flex items-center justify-between">
              <div>
                <h3 className="text-lg font-semibold">{notification.title}</h3>
                <p className="text-sm text-muted-foreground">
                  {new Date(notification.created_at).toLocaleDateString("en-US", { month: "short", day: "numeric", hour: "2-digit", minute: "2-digit" })}
                </p>
                <p className="text-sm text-muted-foreground">{notification.description}</p>
              </div>
              <div className="mt-2 flex justify-end">
                {false === notification.is_read ? (
                  <Button variant="outline" size="sm" onClick={() => handleNotificationRead(notification.id)}>
                    <CheckIcon className="h-4 w-4" />
                  </Button>
                ) : (
                  <CheckCircle2 className="h-6 w-6 mr-2 text-green-600" />
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    )
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
        <div className="">
          {notifications.length > 0 ? (
            <Tabs defaultValue="unread" className="w-full mt-4">
              <TabsList className="grid w-full grid-cols-2">
                <TabsTrigger value="unread">Unread</TabsTrigger>
                <TabsTrigger value="read">Read</TabsTrigger>
              </TabsList>

              <div className="mt-2">
                <TabsContent value="unread">
                  {notifications.filter((notification) => !notification.is_read).length > 0 ? (
                    renderNotifications([
                      ...notifications.filter((notification) => !notification.is_read),
                      ...notifications.filter((notification) => !notification.is_read),
                    ])
                  ) : (
                    <p className="text-center text-muted-foreground py-4">No unread notifications</p>
                  )}
                </TabsContent>

                <TabsContent value="read">
                  {notifications.filter((notification) => notification.is_read).length > 0 ? (
                    renderNotifications(notifications.filter((notification) => notification.is_read))
                  ) : (
                    <p className="text-center text-muted-foreground py-4">No read notifications</p>
                  )}
                </TabsContent>
              </div>
            </Tabs>
          ) : (
            <p className="text-muted-foreground py-4">No notifications yet</p>
          )}
        </div>
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
