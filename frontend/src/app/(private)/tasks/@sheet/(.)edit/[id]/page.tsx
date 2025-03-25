"use client"

import { useEffect, useState } from "react"
import RenderTaskFormSheet from "@/app/_components/task/sheets/render-task-form-sheet.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { axiosInstance } from "@/lib/http/axios"
import { useRouter } from "next/navigation"
import { toast } from "sonner"

interface TaskResponse {
  data: GetOneTaskInterface
}

type Props = {
  params: {
    id: string
  }
}

export default function TaskEditSheetPage({ params }: Props) {
  const { id } = params
  const router = useRouter()
  const [task, setTask] = useState<GetOneTaskInterface | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const fetchTask = async () => {
      try {
        const response = await axiosInstance.get<TaskResponse>(`/api/v1/tasks/${id}`)
        setTask(response.data.data)
      } catch (error: any) {
        console.error("Error fetching task:", error)
        toast.error("Error fetching task", {
          description: error.message || "Failed to load task details",
        })
        router.back()
      } finally {
        setIsLoading(false)
      }
    }

    fetchTask()
  }, [id, router])

  if (isLoading) {
    return <div>Loading...</div>
  }

  if (!task) {
    return null
  }

  return <RenderTaskFormSheet task={task} title="Edit Task" description="Edit the task" overlay={false} isModal={true} open={true} />
}
