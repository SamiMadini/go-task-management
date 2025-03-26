package task

import (
	"context"
	"time"

	"sama/go-task-management/commons"

	"github.com/google/uuid"
)

type Repository interface {
	GetByID(id string) (commons.Task, error)
	GetAll() ([]commons.Task, error)
	Create(task commons.Task) (commons.Task, error)
	Update(task commons.Task) error
	Delete(id string) error
}

type UserRepository interface {
	GetByID(id string) (commons.User, error)
}

type Service struct {
	logger   commons.Logger
	taskRepo Repository
	userRepo UserRepository
}

func NewService(logger commons.Logger, taskRepo Repository, userRepo UserRepository) *Service {
	return &Service{
		logger:   logger,
		taskRepo: taskRepo,
		userRepo: userRepo,
	}
}

func (s *Service) GetTask(ctx context.Context, taskID string, userID string) (*commons.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if userID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != userID) {
		return nil, commons.ErrForbidden
	}

	return &task, nil
}

func (s *Service) GetAllTasks(ctx context.Context, userID string) ([]commons.Task, error) {
	tasks, err := s.taskRepo.GetAll()
	if err != nil {
		return nil, err
	}

	userTasks := make([]commons.Task, 0)
	for _, task := range tasks {
		if task.CreatorID == userID || (task.AssigneeID != nil && *task.AssigneeID == userID) {
			userTasks = append(userTasks, task)
		}
	}

	return userTasks, nil
}

func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (*commons.Task, error) {
	now := time.Now()
	task := commons.Task{
		ID:          uuid.New().String(),
		Title:       input.Title,
		Description: input.Description,
		Status:      input.Status,
		Priority:    input.Priority,
		DueDate:     input.DueDate,
		CreatorID:   input.CreatorID,
		AssigneeID:  input.AssigneeID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	createdTask, err := s.taskRepo.Create(task)
	if err != nil {
		return nil, err
	}

	return &createdTask, nil
}

func (s *Service) UpdateTask(ctx context.Context, taskID string, input UpdateTaskInput) (*commons.Task, error) {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	if input.UserID != task.CreatorID && (task.AssigneeID == nil || *task.AssigneeID != input.UserID) {
		return nil, commons.ErrForbidden
	}

	if input.Title != "" {
		task.Title = input.Title
	}
	if input.Description != "" {
		task.Description = input.Description
	}
	if input.Status != "" {
		task.Status = input.Status
	}
	if input.Priority != 0 {
		task.Priority = input.Priority
	}
	if !input.DueDate.IsZero() {
		task.DueDate = input.DueDate
	}
	if input.AssigneeID != nil {
		task.AssigneeID = input.AssigneeID
	}

	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(task); err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Service) DeleteTask(ctx context.Context, taskID string, userID string) error {
	task, err := s.taskRepo.GetByID(taskID)
	if err != nil {
		return err
	}

	if userID != task.CreatorID {
		return commons.ErrForbidden
	}

	return s.taskRepo.Delete(taskID)
}
