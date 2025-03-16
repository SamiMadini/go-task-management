import { TaskPriority, TaskPriorityLabel, TaskStatus, TaskStatusLabel } from "@/app/domain/task/enums.task"

export const getPriorityLabel = (priority: number): string => {
  switch (priority) {
    case TaskPriority.HIGH:
      return TaskPriorityLabel.HIGH
    case TaskPriority.MEDIUM:
      return TaskPriorityLabel.MEDIUM
    case TaskPriority.LOW:
      return TaskPriorityLabel.LOW
    default:
      return TaskPriorityLabel.MEDIUM
  }
}

export const getStatusLabel = (status: TaskStatus) => {
  switch (status) {
    case TaskStatus.DONE:
      return TaskStatusLabel.DONE
    case TaskStatus.IN_PROGRESS:
      return TaskStatusLabel.IN_PROGRESS
    case TaskStatus.TODO:
      return TaskStatusLabel.TODO
    default:
      return TaskStatusLabel.TODO
  }
}
