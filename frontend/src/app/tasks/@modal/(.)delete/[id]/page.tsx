"use client"

import { ConfirmDialog } from "@/app/_components/common/confirm-dialog"
import { axiosInstance } from "@/lib/http/axios"

type Props = {
  params: {
    id: string
  }
}

export default function TaskDeleteModalPage({ params }: Props) {
  const { id } = params

  return (
    <ConfirmDialog
      key={`delete-task-modal-${id}`}
      title="Delete your task"
      description={<p>Are you sure you want to delete this task?</p>}
      open={true}
      onConfirm={async () => {
        try {
          const resp = await axiosInstance.delete(`/api/tasks/${id}`)

          if (undefined === resp) {
            throw new Error("No response")
          }
        } catch (error) {
          console.log(error)
        }
      }}
    />
  )
}
