"use client"

import TaskCard from "@/app/_components/task/task-card.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"

export default function TaskListComponent({
  tasks,
  selectedTaskId,
  onShowEvents,
}: {
  tasks: GetOneTaskInterface[]
  selectedTaskId: string
  onShowEvents: (taskId: string) => void
}) {
  if (tasks.length === 0) {
    return <p className="text-center text-sm text-muted-foreground">No tasks</p>
  }

  return (
    <div className="space-y-6">
      {tasks.map((task) => (
        <TaskCard key={task.id} task={task} isSelected={selectedTaskId === task.id.toString()} onShowEvents={onShowEvents} />
      ))}
    </div>
  )
}
