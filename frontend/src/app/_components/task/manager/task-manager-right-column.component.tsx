"use client"

import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import TaskSummaryComponent from "@/app/_components/task/task-summary.component"
import { useRef } from "react"
import TaskSystemEventsComponent from "@/app/_components/task/events/task-system-events.component"

export default function TaskManagerRightColumnComponent({
  tasks,
  selectedTaskId,
}: {
  tasks: GetOneTaskInterface[]
  selectedTaskId: string
}) {
  const rightColumnRef = useRef<HTMLDivElement>(null)

  return (
    <div className="flex flex-col gap-4" ref={rightColumnRef}>
      <TaskSummaryComponent tasks={tasks} />
      <TaskSystemEventsComponent
        tasks={tasks}
        selectedTaskId={selectedTaskId}
        rightColumnRef={rightColumnRef}
      />
    </div>
  )
}
