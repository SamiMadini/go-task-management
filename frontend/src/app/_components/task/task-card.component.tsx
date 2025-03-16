"use client"

import { Badge } from "@/components/ui/badge"
import { Card, CardContent } from "@/components/ui/card"
import { CalendarIcon, CheckCircle2, CircleAlert, Clock, Pencil, Trash2, Eye } from "lucide-react"
import { Button } from "@/components/ui/button"
import Link from "next/link"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { TaskPriority, TaskStatus } from "@/app/domain/task/enums.task"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { getPriorityLabel, getStatusLabel } from "@/app/_components/task/task-helper"

export default function TaskCard({
  task,
  isSelected,
  onShowEvents,
}: {
  task: GetOneTaskInterface
  isSelected: boolean
  onShowEvents: (taskId: string) => void
}) {
  return (
    <Card className={isSelected ? "ring-2 ring-primary/50" : ""}>
      <CardContent className="p-4">
        <div className="flex items-start justify-between gap-4">
          <div className="flex-1 space-y-3">
            <div className="flex items-center gap-2">
              {task.status === TaskStatus.DONE ? (
                <CheckCircle2 className="h-5 w-5 text-green-500" />
              ) : task.status === TaskStatus.IN_PROGRESS ? (
                <Clock className="h-5 w-5 text-amber-500" />
              ) : (
                <CircleAlert className="h-5 w-5 text-blue-500" />
              )}
              <h3 className="font-medium">{task.title}</h3>
            </div>

            <p className="text-sm text-muted-foreground">{task.description}</p>

            <div className="flex flex-wrap items-center gap-4 text-sm text-muted-foreground">
              <Badge variant={task.status === TaskStatus.DONE ? "outline" : "secondary"}>{getStatusLabel(task.status)}</Badge>
              <Badge
                variant="outline"
                className={
                  task.priority === TaskPriority.HIGH
                    ? "border-red-200 bg-red-100 text-red-800 dark:border-red-800 dark:bg-red-950 dark:text-red-300"
                    : task.priority === TaskPriority.MEDIUM
                    ? "border-amber-200 bg-amber-100 text-amber-800 dark:border-amber-800 dark:bg-amber-950 dark:text-amber-300"
                    : "border-green-200 bg-green-100 text-green-800 dark:border-green-800 dark:bg-green-950 dark:text-green-300"
                }
              >
                {getPriorityLabel(task.priority)} Priority
              </Badge>

              <div className="flex items-center gap-1">
                <CalendarIcon className="h-3.5 w-3.5" />
                <span>Due {new Date(task.due_date).toLocaleDateString("en-US", { month: "short", day: "numeric" })}</span>
              </div>
            </div>
          </div>

          <div className="flex gap-3 self-start">
            <TooltipProvider>
              <Tooltip>
                <TooltipTrigger asChild>
                  <Button variant={"outline"} disabled={isSelected} size="sm" onClick={() => onShowEvents(task.id)}>
                    <Eye className="h-4 w-4" /> Events
                    <span className="sr-only">Show Events</span>
                  </Button>
                </TooltipTrigger>
                <TooltipContent>
                  <p>Show Events</p>
                </TooltipContent>
              </Tooltip>
            </TooltipProvider>
            <Link href={`/tasks/edit/${task.id}`}>
              <Button variant="outline" size="icon" className="h-8 w-8" aria-label="Edit task">
                <Pencil className="h-4 w-4" />
              </Button>
            </Link>
            <Link href={`/tasks/delete/${task.id}`}>
              <Button variant="outline" size="icon" className="h-8 w-8 text-destructive" aria-label="Delete task">
                <Trash2 className="h-4 w-4" />
              </Button>
            </Link>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}
