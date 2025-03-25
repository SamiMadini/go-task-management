"use client"

import { SheetRight } from "@/app/_components/common/sheet-right"
import { TaskForm } from "@/app/_components/task/forms/task-form"
import { axiosInstance, ApiError } from "@/lib/http/axios"
import { useRouter } from "next/navigation"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { AxiosResponse } from "axios"
import { DateTime } from "luxon"
import { toast } from "sonner"

interface TaskResponse {
  data: GetOneTaskInterface
}

export default function RenderTaskFormSheet({
  task,
  title,
  description,
  overlay,
  isModal,
  open,
}: {
  task: GetOneTaskInterface | null
  title: string
  description: string
  overlay: boolean
  isModal: boolean
  open: boolean
}) {
  const router = useRouter()

  const handleError = (error: unknown, action: string) => {
    const apiError = error as ApiError
    const message = apiError.message || (error instanceof Error ? error.message : "An unknown error occurred")
    const details = apiError.details
    const validationErrors = apiError.validation_errors?.map((err) => `${err.field}: ${err.message}`).join("\n")

    toast.error(`Error ${action}`, {
      description: validationErrors || details || message,
    })
    console.error(`Error ${action}:`, error)
    throw error
  }

  return (
    <SheetRight title={title} description={description} overlay={overlay} isModal={isModal} open={open}>
      <div className="mt-8">
        <TaskForm
          initialData={
            task
              ? {
                  title: task.title,
                  description: task.description,
                  dueDate: new Date(task.due_date),
                  priority: task.priority,
                  status: task.status,
                }
              : undefined
          }
          onSubmit={async (data) => {
            try {
              let response: AxiosResponse<TaskResponse> | null = null
              if (task) {
                const body = {
                  title: data.title,
                  description: data.description,
                  due_date: data.dueDate ? DateTime.fromJSDate(data.dueDate).toISO() : null,
                  priority: data.priority,
                  status: data.status,
                }

                response = await axiosInstance.put(`/api/v1/tasks/${task.id}`, body, {
                  headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                  },
                })
              } else {
                const body = {
                  title: data.title,
                  description: data.description,
                  due_date: data.dueDate ? DateTime.fromJSDate(data.dueDate).toISO() : null,
                  priority: data.priority,
                  status: data.status,
                }

                response = await axiosInstance.post("/api/v1/tasks", body, {
                  headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                  },
                })
              }

              if (!response?.data?.data) {
                throw new Error(`Failed to ${task ? "edit" : "create"} task`)
              }

              toast.success(`Task ${task ? "updated" : "created"} successfully`)
              router.back()
            } catch (error) {
              handleError(error, task ? "editing" : "creating task")
            }
          }}
        />
      </div>
    </SheetRight>
  )
}
