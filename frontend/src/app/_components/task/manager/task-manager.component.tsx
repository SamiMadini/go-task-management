"use client"

import TaskManagerLeftColumnComponent from "@/app/_components/task/manager/task-manager-left-column.component"
import TaskManagerRightColumnComponent from "@/app/_components/task/manager/task-manager-right-column.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { useState } from "react"

export default function TaskManagerComponent({ tasks }: { tasks: GetOneTaskInterface[] }) {
  const [selectedTaskId, setSelectedTaskId] = useState<string>(tasks[0]?.id.toString() || "")

  const handleShowEvents = (taskId: string) => {
    const newTaskId = taskId.toString()
    if (newTaskId !== selectedTaskId) {
      setSelectedTaskId(newTaskId)
    }
  }

  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
      <TaskManagerLeftColumnComponent
        tasks={tasks}
        selectedTaskId={selectedTaskId}
        setSelectedTaskId={setSelectedTaskId}
        handleShowEvents={(taskId: string) => handleShowEvents(taskId)}
      />
      <TaskManagerRightColumnComponent tasks={tasks} selectedTaskId={selectedTaskId} setSelectedTaskId={setSelectedTaskId} />
    </div>
  )
}
