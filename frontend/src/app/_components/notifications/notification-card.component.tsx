"use client"

import { GetOneNotificationInterface } from "@/app/domain/notification/interfaces.notification"
import { Button } from "@/components/ui/button"
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover"
import { BellDot, CheckCircle2, CheckIcon, MoreVertical, Trash2 } from "lucide-react"

export function NotificationCard({
  notifications,
  onRead,
  onDelete,
}: {
  notifications: GetOneNotificationInterface[]
  onRead: (id: number, isRead: boolean) => void
  onDelete: (id: number) => void
}) {
  const renderSubMenu = (notification: GetOneNotificationInterface) => {
    return (
      <Popover key={`notification-card-menu-${notification.id}`}>
        <PopoverTrigger>
          <MoreVertical className="h-5 w-5" />
        </PopoverTrigger>

        <PopoverContent className="w-48 space-y-2" align="end">
          {notification.is_read && (
            <div className="flex items-center gap-2 cursor-pointer">
              <BellDot className="h-5 w-5" /> Mark as unread
            </div>
          )}

          <div className="flex items-center gap-2 cursor-pointer text-red-700" onClick={() => onDelete(notification.id)}>
            <Trash2 className="h-5 w-5" /> Remove
          </div>
        </PopoverContent>
      </Popover>
    )
  }

  return (
    <div className="mt-6 space-y-6">
      {notifications.map((notification) => (
        <div key={notification.id} className="rounded-lg border p-4 relative">
          {!notification.is_read && <span className="absolute -top-1 -left-1 h-3 w-3 rounded-full bg-red-500 animate-pulse" />}

          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-semibold">{notification.title}</h3>
              <p className="text-sm text-muted-foreground">
                {new Date(notification.created_at).toLocaleDateString("en-US", {
                  month: "short",
                  day: "numeric",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </p>
              <p className="text-sm text-muted-foreground">{notification.description}</p>
            </div>

            <div className="mt-2 flex justify-end ">
              {false === notification.is_read ? (
                <Button variant="outline" size="sm" onClick={() => onRead(notification.id, true)} className="absolute bottom-4 right-4">
                  <CheckIcon className="h-4 w-4" />
                </Button>
              ) : (
                <CheckCircle2 className="h-6 w-6 text-green-600 absolute bottom-4 right-4" />
              )}
            </div>
          </div>

          <div className="absolute top-2 right-2">{renderSubMenu(notification)}</div>
        </div>
      ))}
    </div>
  )
}
