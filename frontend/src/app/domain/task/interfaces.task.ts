import { TaskPriority, TaskStatus } from "@/app/domain/task/enums.task"

export interface GetOneTaskInterface {
  id: string
  correlation_id: string
  title: string
  description: string
  status: TaskStatus
  priority: TaskPriority
  in_app_sent: boolean
  email_sent: boolean
  due_date: string
  created_at: string
  updated_at: string
  events: GetTaskSystemEventInterface[]
}

export interface GetTaskSystemEventInterface {
  id: string
  task_id: string
  correlation_id: string
  origin: string
  action: string
  message: string
  json_data: string
  emit_at: string
  created_at: string
}

export interface TaskSystemEventGroupInterface {
  task_id: string
  correlation_id: string
  events: GetTaskSystemEventInterface[]
}
