import { axiosInstance } from "@/lib/http/axios"
import TaskManagerComponent from "@/app/_components/task/manager/task-manager.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"

async function getData(): Promise<{ tasks: GetOneTaskInterface[] }> {
  try {
    const res = await axiosInstance.get("/api/tasks")

    if (!res.data) {
      throw new Error("Failed to fetch data")
    }

    return {
      tasks: res.data,
    }
  } catch (error) {
    return {
      tasks: [],
    }
  }
}

export default async function TasksPage() {
  const data = await getData()
  return <TaskManagerComponent tasks={data.tasks} />
}
