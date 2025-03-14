"use client"

import { CalendarIcon, Filter, Plus, Search } from "lucide-react"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { GetOneTaskInterface } from "@/app/tasks/domain/interfaces.task"

export default function TaskManagerRightColumnComponent({ tasks }: { tasks: GetOneTaskInterface[] }) {
  return (
    <div>
      <Card>
        <CardHeader>
          <CardTitle>Task Summary</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="mr-2 h-4 w-4 rounded-full bg-blue-500" />
                <span>To Do</span>
              </div>
              <span className="font-medium">{tasks.filter((t) => t.status === "todo").length}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="mr-2 h-4 w-4 rounded-full bg-amber-500" />
                <span>In Progress</span>
              </div>
              <span className="font-medium">{tasks.filter((t) => t.status === "in_progress").length}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <div className="mr-2 h-4 w-4 rounded-full bg-green-500" />
                <span>Done</span>
              </div>
              <span className="font-medium">{tasks.filter((t) => t.status === "done").length}</span>
            </div>
          </div>

          <div className="mt-6">
            <h3 className="mb-4 font-medium">Upcoming Deadlines</h3>
            <div className="space-y-3">
              {tasks
                .sort((a, b) => new Date(a.due_date).getTime() - new Date(b.due_date).getTime())
                .slice(0, 3)
                .map((task) => (
                  <div key={task.id} className="flex items-start gap-2">
                    <CalendarIcon className="mt-0.5 h-4 w-4 text-muted-foreground" />
                    <div>
                      <p className="text-sm font-medium line-clamp-1">{task.title}</p>
                      <p className="text-xs text-muted-foreground">
                        Due {new Date(task.due_date).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
                      </p>
                    </div>
                  </div>
                ))}
            </div>
          </div>

          <div className="mt-6">
            <h3 className="mb-4 font-medium">Team Members</h3>
            <div className="flex -space-x-2">
              <Avatar className="border-2 border-background">
                <AvatarImage src="/placeholder.svg?height=32&width=32" />
                <AvatarFallback>JD</AvatarFallback>
              </Avatar>
              <Avatar className="border-2 border-background">
                <AvatarImage src="/placeholder.svg?height=32&width=32" />
                <AvatarFallback>AB</AvatarFallback>
              </Avatar>
              <Avatar className="border-2 border-background">
                <AvatarImage src="/placeholder.svg?height=32&width=32" />
                <AvatarFallback>CK</AvatarFallback>
              </Avatar>
              <Avatar className="border-2 border-background">
                <AvatarImage src="/placeholder.svg?height=32&width=32" />
                <AvatarFallback>DM</AvatarFallback>
              </Avatar>
              <Avatar className="border-2 border-background">
                <AvatarFallback className="bg-primary text-primary-foreground">+3</AvatarFallback>
              </Avatar>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
