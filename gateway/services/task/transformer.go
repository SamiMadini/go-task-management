package task

import (
	"sama/go-task-management/commons"
	"sama/go-task-management/gateway/services/auth"
	"time"
)

type TaskResponse struct {
	ID          string                    `json:"id"`
	Title       string                    `json:"title"`
	Description string                    `json:"description"`
	Status      string                    `json:"status"`
	Priority    int                       `json:"priority"`
	DueDate     time.Time                 `json:"due_date"`
	CreatedAt   time.Time                 `json:"created_at"`
	UpdatedAt   time.Time                 `json:"updated_at"`
	Creator     auth.UserResponse         `json:"creator"`
	Assignee    *auth.UserResponse        `json:"assignee,omitempty"`
	Events      []TaskSystemEventResponse `json:"events"`
}

type TaskSystemEventResponse struct {
	ID            string    `json:"id"`
	TaskId        string    `json:"task_id"`
	CorrelationId string    `json:"correlation_id"`
	Origin        string    `json:"origin"`
	Action        string    `json:"action"`
	Message       string    `json:"message"`
	JsonData      string    `json:"json_data"`
	EmitAt        time.Time `json:"emit_at"`
}

func ToTaskResponse(task commons.Task, creator commons.User, assignee *commons.User) TaskResponse {
	response := TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		DueDate:     task.DueDate,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
		Creator:     auth.ToUserResponse(creator),
		Events:      make([]TaskSystemEventResponse, len(task.Events)),
	}

	if assignee != nil {
		assigneeResponse := auth.ToUserResponse(*assignee)
		response.Assignee = &assigneeResponse
	}

	for i, event := range task.Events {
		response.Events[i] = TaskSystemEventResponse{
			ID:            event.ID,
			TaskId:        event.TaskId,
			CorrelationId: event.CorrelationId,
			Origin:        event.Origin,
			Action:        event.Action,
			Message:       event.Message,
			JsonData:      event.JsonData,
			EmitAt:        event.EmitAt,
		}
	}

	return response
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}

func ToTaskListResponse(tasks []commons.Task, users map[string]commons.User) TaskListResponse {
	response := TaskListResponse{
		Tasks: make([]TaskResponse, len(tasks)),
	}

	for i, task := range tasks {
		creator := users[task.CreatorID]
		var assignee *commons.User
		if task.AssigneeID != nil {
			if assigneeUser, ok := users[*task.AssigneeID]; ok {
				assignee = &assigneeUser
			}
		}
		response.Tasks[i] = ToTaskResponse(task, creator, assignee)
	}

	return response
}
