"use client"

import { SheetRight } from "@/app/_components/common/sheet-right"
import { TaskForm } from "@/app/_components/task/forms/task-form"
import { axiosInstance } from "@/lib/http/axios"
import { useRouter } from "next/navigation"

export default function TaskCreateSheetPage() {
  const router = useRouter()

  return (
    <SheetRight title="Create a Task" description="Fill all fields to create the task" overlay={false} isModal={true} open={true}>
      <div className="mt-8">
        <TaskForm
          onSubmit={async (data) => {
            console.log(data)
            try {
              const params = {
                title: data.title,
                description: data.description,
                dueDate: data.dueDate,
                priority: data.priority,
                status: data.status,
              }

              const response = await axiosInstance.post("/api/tasks", params, {
                headers: {
                  "Content-Type": "application/json",
                  Accept: "application/json",
                },
              })

              if (!response.data) {
                throw new Error("Failed to create task")
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
