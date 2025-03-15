"use client"

import { SheetRight } from "@/app/_components/common/sheet-right"
import { TaskForm } from "@/app/_components/task/forms/task-form"
import { GetOneTaskInterface } from "@/app/tasks/domain/interfaces.task"
import { axiosInstance } from "@/lib/http/axios"
import { useRouter } from "next/navigation"

type Props = {
  params: {
    id: string
  }
}

export default async function TaskEditSheetPage({ params }: Props) {
  const { id } = params
  const router = useRouter()

  let task: GetOneTaskInterface | null = null
  try {
    const response = await axiosInstance.get(`/api/tasks/${id}`)
    task = response.data
    console.log("task")
    console.log(task)
  } catch (error) {
    console.error("Error fetching task:", error)
  }

  if (!task) {
    return <div>Task not found</div>
  }

  return (
    <SheetRight title="Edit a Task" description="Fill all fields to edit the task" overlay={false} isModal={true} open={true}>
      <div className="mt-8">
        <TaskForm
          initialData={{
            title: task.title,
            description: task.description,
            dueDate: new Date(task.due_date),
            priority: task.priority,
            status: task.status,
          }}
          onSubmit={async (data) => {
            console.log(data)
            try {
              const body = {
                title: data.title,
                description: data.description,
                dueDate: data.dueDate,
                priority: data.priority,
                status: data.status,
              }

              const response = await axiosInstance.put(`/api/tasks/${id}`, body, {
                headers: {
                  "Content-Type": "application/json",
                  Accept: "application/json",
                },
              })

              if (!response.data) {
                throw new Error("Failed to edit task")
              }

              router.back()
            } catch (error) {
              console.error("Error creating task:", error)
              throw error
            }
          }}
        />
      </div>
    </SheetRight>
  )
}
