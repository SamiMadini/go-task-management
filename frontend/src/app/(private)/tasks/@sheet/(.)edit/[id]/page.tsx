import RenderTaskFormSheet from "@/app/_components/task/sheets/render-task-form-sheet.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { axiosInstance } from "@/lib/http/axios"

type Props = {
  params: {
    id: string
  }
}

export default async function TaskEditSheetPage({ params }: Props) {
  const { id } = params

  let task: GetOneTaskInterface | null = null
  try {
    const response = await axiosInstance.get(`/api/v1/tasks/${id}`)
    task = response.data
  } catch (error) {
    console.error("Error fetching task:", error)
  }

  if (!task) {
    return <div>Task not found</div>
  }

  return <RenderTaskFormSheet task={task} title="Edit Task" description="Edit the task" overlay={false} isModal={true} open={true} />
}
