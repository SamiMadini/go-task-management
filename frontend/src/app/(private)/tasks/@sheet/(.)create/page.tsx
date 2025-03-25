import RenderTaskFormSheet from "@/app/_components/task/sheets/render-task-form-sheet.component"

export default function TaskCreateSheetPage() {
  return (
    <RenderTaskFormSheet
      title="Create a Task"
      description="Fill all fields to create the task"
      overlay={false}
      isModal={true}
      open={true}
      task={null}
    />
  )
}
