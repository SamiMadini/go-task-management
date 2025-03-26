"use client"

import { ConfirmDialog } from "@/app/_components/common/confirm-dialog"
import { axiosInstance, ApiError } from "@/lib/http/axios"
import { toast } from "sonner"

type Props = {
  params: {
    id: string
  }
}

export default function TaskDeleteModalPage({ params }: Props) {
  const { id } = params

  const handleError = (error: unknown) => {
    const apiError = error as ApiError
    const message = apiError.message || (error instanceof Error ? error.message : "An unknown error occurred")
    const details = apiError.details

    toast.error("Error deleting task", {
      description: details || message,
    })
    console.error("Error deleting task:", error)
  }

  return (
    <ConfirmDialog
      key={`delete-task-modal-${id}`}
      title="Delete your task"
      description={<p>Are you sure you want to delete this task?</p>}
      open={true}
      onConfirm={async () => {
        try {
          const resp = await axiosInstance.delete(`/api/v1/tasks/${id}`)
          if (!resp.data?.success) {
            throw new Error("Failed to delete task")
          }
          toast.success("Task deleted successfully")
        } catch (error) {
          handleError(error)
        }
      }}
    />
  )
}
