"use client"

import { SheetRight } from "@/app/_components/common/sheet-right"
import { TaskForm } from "@/app/_components/task/forms/task-form"
import { axiosInstance } from "@/lib/http/axios"
import { useRouter } from "next/navigation"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { AxiosResponse } from "axios"

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
            console.log(data)
            try {
              let response: AxiosResponse<GetOneTaskInterface> | null = null
              if (task) {
                const body = {
                  title: data.title,
                  description: data.description,
                  dueDate: data.dueDate,
                  priority: data.priority,
                  status: data.status,
                }

                response = await axiosInstance.put(`/api/tasks/${task.id}`, body, {
                  headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                  },
                })
              } else {
                const body = {
                  title: data.title,
                  description: data.description,
                  dueDate: data.dueDate,
                  priority: data.priority,
                  status: data.status,
                }

                response = await axiosInstance.post("/api/tasks", body, {
                  headers: {
                    "Content-Type": "application/json",
                    Accept: "application/json",
                  },
                })
              }

              if (!response?.data) {
                throw new Error(`Failed to ${task ? "edit" : "create"} task`)
              }

              router.back()
            } catch (error) {
              console.error(`Error ${task ? "editing" : "creating"} task:`, error)
              throw error
            }
          }}
        />
      </div>
    </SheetRight>
  )
}
