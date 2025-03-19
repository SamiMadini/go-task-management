"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible"
import { useEffect, useRef, useState } from "react"
import { Button } from "@/components/ui/button"
import { ChevronDown, ChevronUp, Grid, List } from "lucide-react"
import { GetOneTaskInterface, GetTaskSystemEventInterface, TaskSystemEventGroupInterface } from "@/app/domain/task/interfaces.task"
import TaskSystemEventGroupComponent from "@/app/_components/task/events/task-system-event-group.component"
import TaskSystemEventGroupFilterComponent from "@/app/_components/task/events/task-system-event-group-filters.component"

export default function TaskSystemEventsComponent({
  tasks,
  selectedTaskId,
  rightColumnRef,
}: {
  tasks: GetOneTaskInterface[]
  selectedTaskId: string
  rightColumnRef: React.RefObject<HTMLDivElement>
}) {
  const [eventsOpen, setEventsOpen] = useState(true)
  const [fullDetails, setFullDetails] = useState(false)
  const [isEventsFixed, setIsEventsFixed] = useState(false)
  const eventsCardRef = useRef<HTMLDivElement>(null)

  const buildEventsGroupsDictionary = (events: GetTaskSystemEventInterface[]): { [key: string]: TaskSystemEventGroupInterface } => {
    const eventsGroups: { [key: string]: TaskSystemEventGroupInterface } = {}

    const sortedEvents = events.sort((a, b) => new Date(a.emit_at).getTime() - new Date(b.emit_at).getTime())

    sortedEvents.forEach((event) => {
      if (!eventsGroups[event.correlation_id]) {
        eventsGroups[event.correlation_id] = {
          task_id: event.task_id,
          correlation_id: event.correlation_id,
          events: [event],
        }
      } else {
        eventsGroups[event.correlation_id].events.push(event)
      }
    })

    return eventsGroups
  }

  const selectedTask = tasks.find((t) => t.id.toString() === selectedTaskId) || tasks[0]
  const selectedTaskEvents = {
    taskId: selectedTask?.id.toString() || "",
    events: buildEventsGroupsDictionary(selectedTask?.events || []),
  }

  // useEffect(() => {
  //   const newOriginFilters: Record<string, boolean> = {}
  //   uniqueOrigins.forEach((origin) => {
  //     newOriginFilters[origin] = true // All enabled by default
  //   })
  //   setOriginFilters(newOriginFilters)

  //   const newActionFilters: Record<string, boolean> = {}
  //   uniqueActions.forEach((action) => {
  //     newActionFilters[action] = true // All enabled by default
  //   })
  //   setActionFilters(newActionFilters)
  // }, [selectedTaskId, uniqueOrigins, uniqueActions])

  useEffect(() => {
    const handleScroll = () => {
      if (!eventsCardRef.current || !rightColumnRef.current) return

      const rightColumnRect = rightColumnRef.current.getBoundingClientRect()
      const topBarHeight = 64

      if (rightColumnRect.top < topBarHeight) {
        setIsEventsFixed(true)
      } else {
        setIsEventsFixed(false)
      }
    }

    window.addEventListener("scroll", handleScroll)
    return () => window.removeEventListener("scroll", handleScroll)
  }, [])

  return (
    <div
      ref={eventsCardRef}
      className={`${isEventsFixed ? "lg:sticky lg:top-[80px]" : ""} transition-all duration-200`}
      style={{
        width: isEventsFixed ? rightColumnRef.current?.offsetWidth || "auto" : "auto",
        zIndex: 10,
      }}
    >
      <Card className={`${isEventsFixed ? "lg:shadow-md" : ""}`}>
        <Collapsible open={eventsOpen} onOpenChange={setEventsOpen}>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle>System Events</CardTitle>
              <CollapsibleTrigger asChild>
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  {eventsOpen ? <ChevronUp className="h-4 w-4" /> : <ChevronDown className="h-4 w-4" />}
                  <span className="sr-only">Toggle</span>
                </Button>
              </CollapsibleTrigger>
            </div>
          </CardHeader>

          <CollapsibleContent>
            {selectedTaskEvents && selectedTaskEvents.taskId && (
              <CardContent>
                {selectedTaskEvents && selectedTaskEvents.events && (
                  <div className="flex items-center justify-between mb-6">
                    <h3 className="text-lg">History</h3>

                    <div className="flex items-center gap-2">
                      <Button
                        variant={!fullDetails ? "outline" : "ghost"}
                        size="sm"
                        className="h-8 w-8 p-0"
                        aria-label="List view"
                        onClick={() => setFullDetails(false)}
                      >
                        <List className="h-4 w-4" />
                        <span className="sr-only">List view</span>
                      </Button>

                      <Button
                        variant={fullDetails ? "outline" : "ghost"}
                        size="sm"
                        className="h-8 w-8 p-0"
                        aria-label="Grid view"
                        onClick={() => setFullDetails(true)}
                      >
                        <Grid className="h-4 w-4" />
                        <span className="sr-only">Grid view</span>
                      </Button>
                    </div>

                    <TaskSystemEventGroupFilterComponent selectedTask={selectedTask} selectedTaskEvents={selectedTaskEvents} />
                  </div>
                )}

                {selectedTaskEvents && (
                  <div className="space-y-6 max-h-[1200px] overflow-y-auto pr-2">
                    {Object.values(selectedTaskEvents.events).map((group, groupIndex) => (
                      <TaskSystemEventGroupComponent key={group.correlation_id} group={group} groupIndex={groupIndex} full={fullDetails} />
                    ))}
                  </div>
                )}
              </CardContent>
            )}
          </CollapsibleContent>
        </Collapsible>
      </Card>
    </div>
  )
}
