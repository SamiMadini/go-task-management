import { TaskPriority, TaskStatus } from "@/app/domain/task/enums.task"

export interface GetOneTaskInterface {
  id: string
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  in_app_sent: boolean
  email_sent: boolean
  due_date: string
  created_at: string
  updated_at: string
}
