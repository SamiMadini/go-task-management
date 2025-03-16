"use client"

import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs"
import { NotificationCard } from "@/app/_components/notifications/notification-card.component"
import { GetOneNotificationInterface } from "@/app/domain/notification/interfaces.notification"

export function NotificationsTabsComponent({
  notifications,
  onRead,
  onDelete,
}: {
  notifications: GetOneNotificationInterface[]
  onRead: (notificationId: number, isRead: boolean) => void
  onDelete: (notificationId: number) => void
}) {
  return (
    <>
      {notifications.length > 0 ? (
        <Tabs defaultValue="unread" className="w-full mt-4">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="unread">Unread</TabsTrigger>
            <TabsTrigger value="read">Read</TabsTrigger>
          </TabsList>

          <div className="mt-2">
            <TabsContent value="unread">
              {notifications.filter((notification) => !notification.is_read).length > 0 ? (
                <NotificationCard
                  notifications={notifications.filter((notification) => !notification.is_read)}
                  onRead={(notificationId, isRead) => onRead(notificationId, isRead)}
                  onDelete={(notificationId) => onDelete(notificationId)}
                />
              ) : (
                <p className="text-center text-muted-foreground py-4">No unread notifications</p>
              )}
            </TabsContent>

            <TabsContent value="read">
              {notifications.filter((notification) => notification.is_read).length > 0 ? (
                <NotificationCard
                  notifications={notifications.filter((notification) => notification.is_read)}
                  onRead={(notificationId, isRead) => onRead(notificationId, isRead)}
                  onDelete={(notificationId) => onDelete(notificationId)}
                />
              ) : (
                <p className="text-center text-muted-foreground py-4">No read notifications</p>
              )}
            </TabsContent>
          </div>
        </Tabs>
      ) : (
        <p className="text-muted-foreground py-4">No notifications yet</p>
      )}
    </>
  )
}
