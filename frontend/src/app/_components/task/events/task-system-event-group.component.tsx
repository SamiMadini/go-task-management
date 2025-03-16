import { TaskSystemEventGroupInterface } from "@/app/domain/task/interfaces.task"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { ChevronDown } from "lucide-react"
import { Separator } from "@/components/ui/separator"
import { ChevronUp } from "lucide-react"
import { useState } from "react"

export default function TaskSystemEventGroupComponent({
  group,
  groupIndex,
  full,
}: {
  group: TaskSystemEventGroupInterface
  groupIndex: number
  full: boolean
}) {
  const [isOpen, setIsOpen] = useState(groupIndex === 0)

  return (
    <Collapsible open={isOpen} onOpenChange={setIsOpen}>
      <div key={group.correlation_id} className="flex flex-col gap-2 space-y-5 p-2">
        <div className="flex items-center justify-between">
          <div className="flex flex-col gap-2">
            <h4 className="flex flex-row text-sm font-medium gap-3">
              Main event
              <Badge variant={"secondary"} className="">
                {group.events[0].action}
              </Badge>
            </h4>

            <div className="flex flex-wrap items-center text-sm gap-2">
              Correlation ID
              <Badge
                variant={"outline"}
                className="border-amber-200 bg-amber-100 text-amber-800 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-300"
              >
                {group.correlation_id}
              </Badge>
            </div>
          </div>

          <div className="flex flex-row text-xs text-muted-foreground items-center gap-2">
            {group.events.length} {group.events.length === 1 ? "event" : "events"}
            <CollapsibleTrigger asChild>
              <Button variant="ghost" size="sm" className="h-8 w-8">
                {isOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                <span className="sr-only">Toggle</span>
              </Button>
            </CollapsibleTrigger>
          </div>
        </div>

        <CollapsibleContent>
          <div className="space-y-4 max-h-[350px] overflow-y-auto pr-2">
            {group.events.length > 0 ? (
              group.events.map((event, index) => (
                <>
                  {full && index !== 0 && (
                    <div className="flex flex-col items-center justify-center relative">
                      <Separator orientation="vertical" className="mx-2 h-4 bg-muted-foreground/80" />
                      <ChevronDown className="h-4 w-4 text-muted-foreground absolute -bottom-2" />
                    </div>
                  )}

                  <div key={event.id} className="rounded-md border p-3 text-sm">
                    <div className="flex items-center justify-between mb-1">
                      <div className="flex items-center gap-2">
                        <span className="text-xs text-muted-foreground">{new Date(event.emit_at).toLocaleTimeString()}</span>
                        <Badge
                          variant="outline"
                          className={
                            event.origin === "API Gateway"
                              ? "border-blue-200 bg-blue-100 text-blue-800 dark:border-blue-800 dark:bg-blue-950 dark:text-blue-300"
                              : event.origin === "Notification Service"
                              ? "border-purple-200 bg-purple-100 text-purple-800 dark:border-purple-800 dark:bg-purple-950 dark:text-purple-300"
                              : event.origin === "Email Service"
                              ? "border-green-200 bg-green-100 text-green-800 dark:border-green-800 dark:bg-green-950 dark:text-green-300"
                              : "border-orange-200 bg-orange-100 text-orange-800 dark:border-orange-800 dark:bg-orange-950 dark:text-orange-300"
                          }
                        >
                          {event.origin}
                        </Badge>
                      </div>

                      <div className="mb-1">
                        <code className="text-xs font-mono text-muted-foreground">{event.action}</code>
                      </div>
                    </div>
                    {full && <p className="text-sm">{event.message}</p>}
                  </div>
                </>
              ))
            ) : (
              <div className="text-center py-8 text-muted-foreground">
                <p>No events match the current filters</p>
                <Button
                  variant="link"
                  className="mt-2 h-auto p-0"
                  onClick={() => {
                    // toggleAllOrigins(true)
                    // toggleAllActions(true)
                  }}
                >
                  Reset filters
                </Button>
              </div>
            )}
          </div>
        </CollapsibleContent>
      </div>
      {groupIndex !== Object.values(group.events).length - 1 && <Separator orientation="horizontal" className="mt-2" />}
    </Collapsible>
  )
}
