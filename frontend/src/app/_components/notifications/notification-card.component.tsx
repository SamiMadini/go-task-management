"use client"

import { GetOneNotificationInterface } from "@/app/domain/notification/interfaces.notification"
import { Button } from "@/components/ui/button"
import { CheckCircle2, CheckIcon } from "lucide-react"

export function NotificationCard({
  notifications,
  onRead,
}: {
  notifications: GetOneNotificationInterface[]
  onRead: (id: number) => void
}) {
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
            <div className="mt-2 flex justify-end">
              {false === notification.is_read ? (
                <Button variant="outline" size="sm" onClick={() => onRead(notification.id)}>
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
