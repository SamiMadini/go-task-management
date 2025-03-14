"use client"

import TaskManagerLeftColumnComponent from "@/app/_components/task/manager/task-manager-left-column.component"
import TaskManagerRightColumnComponent from "@/app/_components/task/manager/task-manager-right-column.component"
import { GetOneTaskInterface } from "@/app/tasks/domain/interfaces.task"

export default function TaskManagerComponent({ tasks }: { tasks: GetOneTaskInterface[] }) {
  return (
    <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
      <TaskManagerLeftColumnComponent tasks={tasks} />
      <TaskManagerRightColumnComponent tasks={tasks} />
    </div>
  )
}
