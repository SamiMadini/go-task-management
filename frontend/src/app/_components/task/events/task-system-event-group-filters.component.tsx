"use client"

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { Checkbox } from "@/components/ui/checkbox"
import { Button } from "@/components/ui/button"
import { ListFilter } from "lucide-react"
import { GetOneTaskInterface, GetTaskSystemEventInterface, TaskSystemEventGroupInterface } from "@/app/domain/task/interfaces.task"
import { useState } from "react"

export default function TaskSystemEventGroupFilterComponent({
  selectedTask,
  selectedTaskEvents,
}: {
  selectedTask: GetOneTaskInterface | undefined
  selectedTaskEvents: {
    taskId: string
    events: { [key: string]: TaskSystemEventGroupInterface }
  }
}) {
  const [originFilters, setOriginFilters] = useState<Record<string, boolean>>({})
  const [actionFilters, setActionFilters] = useState<Record<string, boolean>>({})

  if (!selectedTask) {
    return null
  }

  const allEvents = selectedTask?.events || []
  const uniqueOrigins = Array.from(new Set(allEvents.map((event: GetTaskSystemEventInterface) => event.origin)))
  const uniqueActions = Array.from(new Set(allEvents.map((event: GetTaskSystemEventInterface) => event.action)))

  // const filteredEvents = selectedTask?.events
  //   ? selectedTask.events.filter((event: GetTaskSystemEventInterface) => originFilters[event.origin] && actionFilters[event.action])
  //   : []

  const toggleAllOrigins = (value: boolean) => {
    const newFilters = { ...originFilters }
    Object.keys(newFilters).forEach((key) => {
      newFilters[key] = value
    })
    setOriginFilters(newFilters)
  }

  const toggleAllActions = (value: boolean) => {
    const newFilters = { ...actionFilters }
    Object.keys(newFilters).forEach((key) => {
      newFilters[key] = value
    })
    setActionFilters(newFilters)
  }

  const hasActiveFilters = () => {
    return Object.values(originFilters).some((value) => !value) || Object.values(actionFilters).some((value) => !value)
  }

  return (
    <div className="flex items-center gap-2">
      <span className="text-xs text-muted-foreground">
        {Object.values(selectedTaskEvents.events).length} {Object.values(selectedTaskEvents.events).length === 1 ? "group" : "groups"}
      </span>
      <span className="text-xs text-muted-foreground">
        {Object.values(selectedTaskEvents.events).reduce((acc, group) => acc + group.events.length, 0)}{" "}
        {Object.values(selectedTaskEvents.events).reduce((acc, group) => acc + group.events.length, 0) === 1 ? "event" : "events"}
        {hasActiveFilters() && " (filtered)"}
      </span>

      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button size="sm" className={`${hasActiveFilters() ? "bg-primary/10" : ""} mr-2`}>
            <ListFilter className="h-4 w-4" /> Filters
            <span className="sr-only">Filter events</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-56">
          <DropdownMenuLabel>Filter Events</DropdownMenuLabel>

          <DropdownMenuGroup>
            <DropdownMenuLabel className="text-xs font-normal text-muted-foreground pt-2">By Origin</DropdownMenuLabel>
            <div className="px-2 py-1">
              <div className="flex items-center mb-2">
                <Button variant="ghost" size="sm" className="h-6 text-xs" onClick={() => toggleAllOrigins(true)}>
                  Select All
                </Button>
                <Button variant="ghost" size="sm" className="h-6 text-xs" onClick={() => toggleAllOrigins(false)}>
                  Clear All
                </Button>
              </div>
              {uniqueOrigins.map((origin) => (
                <div key={origin} className="flex items-center space-x-2 py-1">
                  <Checkbox
                    id={`origin-${origin}`}
                    checked={originFilters[origin]}
                    onCheckedChange={(checked) => {
                      setOriginFilters({
                        ...originFilters,
                        [origin]: !!checked,
                      })
                    }}
                  />
                  <label
                    htmlFor={`origin-${origin}`}
                    className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
                  >
                    {origin}
                  </label>
                </div>
              ))}
            </div>
          </DropdownMenuGroup>

          <DropdownMenuSeparator />

          <DropdownMenuGroup>
            <DropdownMenuLabel className="text-xs font-normal text-muted-foreground pt-2">By Action</DropdownMenuLabel>
            <div className="px-2 py-1">
              <div className="flex items-center mb-2">
                <Button variant="ghost" size="sm" className="h-6 text-xs" onClick={() => toggleAllActions(true)}>
                  Select All
                </Button>
                <Button variant="ghost" size="sm" className="h-6 text-xs" onClick={() => toggleAllActions(false)}>
                  Clear All
                </Button>
              </div>
              {uniqueActions.map((action) => (
                <div key={action} className="flex items-center space-x-2 py-1">
                  <Checkbox
                    id={`action-${action}`}
                    checked={actionFilters[action]}
                    onCheckedChange={(checked) => {
                      setActionFilters({
                        ...actionFilters,
                        [action]: !!checked,
                      })
                    }}
                  />
                  <label
                    htmlFor={`action-${action}`}
                    className="text-sm font-normal leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 truncate"
                  >
                    {action}
                  </label>
                </div>
              ))}
            </div>
          </DropdownMenuGroup>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
