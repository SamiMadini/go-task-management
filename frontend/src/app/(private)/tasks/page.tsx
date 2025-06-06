"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import { axiosInstance } from "@/lib/http/axios"
import TaskManagerComponent from "@/app/_components/task/manager/task-manager.component"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { useAuth } from "@/lib/hooks"
import { store } from "@/lib/store"

interface TasksResponse {
  data: GetOneTaskInterface[]
}

export default function TasksPage() {
  const router = useRouter()
  const { isAuthenticated, accessToken } = useAuth()
  const [tasks, setTasks] = useState<GetOneTaskInterface[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const state = store.getState()
    console.log("Auth state in tasks page:", {
      isAuthenticated,
      hasAccessToken: !!accessToken,
      accessToken: accessToken?.substring(0, 10) + "...",
      reduxState: {
        user: state.auth.user,
        accessToken: state.auth.accessToken?.substring(0, 10) + "...",
        refreshToken: state.auth.refreshToken?.substring(0, 10) + "...",
      },
    })

    if (!isAuthenticated) {
      console.log("Not authenticated, redirecting to signin")
      router.push("/auth/signin")
      return
    }

    const fetchTasks = async () => {
      try {
        console.log("Fetching tasks with token:", accessToken?.substring(0, 10) + "...")
        const res = await axiosInstance.get<TasksResponse>("/api/v1/tasks")
        if (!res.data?.data) {
          throw new Error("Failed to fetch data")
        }
        setTasks(res.data.data)
      } catch (error: any) {
        console.error("Error fetching tasks:", error)
        if (error.response?.status === 401) {
          console.log("Received 401, redirecting to signin")
          router.push("/auth/signin")
        }
      } finally {
        setIsLoading(false)
      }
    }

    fetchTasks()
  }, [isAuthenticated, router, accessToken])

  if (!isAuthenticated) {
    return null
  }

  if (isLoading) {
    return <div>Loading...</div>
  }

  return <TaskManagerComponent tasks={tasks} />
}
