"use client"

import { Plus, Search } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { useState } from "react"
import TaskListComponent from "@/app/_components/task/task-list.component"
import Link from "next/link"
import { GetOneTaskInterface } from "@/app/domain/task/interfaces.task"
import { getPriorityLabel } from "@/app/_components/task/task-helper"

export default function TaskManagerLeftColumnComponent({
  tasks,
  selectedTaskId,
  handleShowEvents,
}: {
  tasks: GetOneTaskInterface[]
  selectedTaskId: string
  handleShowEvents: (taskId: string) => void
}) {
  const [activeTab, setActiveTab] = useState("all")
  const [searchTerm, setSearchTerm] = useState("")

  const filteredTasks =
    activeTab === "all"
      ? tasks
      : tasks.filter((task) => {
          if (activeTab === "todo") return task.status === "todo"
          if (activeTab === "inprogress") return task.status === "in_progress"
          if (activeTab === "done") return task.status === "done"
          return true
        })

  const searchFilteredTasks = searchTerm
    ? filteredTasks.filter((task) => {
        const searchLower = searchTerm.toLowerCase()
        return (
          (task.title && task.title.toLowerCase().includes(searchLower)) ||
          (task.description && task.description.toLowerCase().includes(searchLower)) ||
          (task.status && task.status.toLowerCase().includes(searchLower)) ||
          (task.priority && getPriorityLabel(task.priority).toLowerCase().includes(searchLower))
        )
      })
    : filteredTasks

  return (
    <div className="">
      <div className="mb-6 flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Tasks</h1>
          <p className="text-muted-foreground">Manage and track your tasks</p>
        </div>
        <div className="flex items-center gap-2">
          <Link href={`/tasks/create`}>
            <Button size="sm">
              <Plus className="h-4 w-4" />
              New Task
            </Button>
          </Link>
        </div>
      </div>

      <div className="mb-6 flex flex-col gap-4 sm:flex-row">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-muted-foreground" />
          <Input
            type="search"
            placeholder="Search tasks..."
            className="w-full pl-8"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
      </div>

      <Tabs defaultValue="all" className="mb-6" onValueChange={setActiveTab}>
        <TabsList>
          <TabsTrigger value="all">All Tasks</TabsTrigger>
          <TabsTrigger value="todo">To Do</TabsTrigger>
          <TabsTrigger value="inprogress">In Progress</TabsTrigger>
          <TabsTrigger value="done">Done</TabsTrigger>
        </TabsList>
        <TabsContent value="all" className="mt-0" />
        <TabsContent value="todo" className="mt-0" />
        <TabsContent value="inprogress" className="mt-0" />
        <TabsContent value="done" className="mt-0" />
      </Tabs>

      <TaskListComponent tasks={searchFilteredTasks} selectedTaskId={selectedTaskId} onShowEvents={handleShowEvents} />
    </div>
  )
}
