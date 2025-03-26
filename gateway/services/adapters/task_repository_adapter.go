package adapters

import (
	"sama/go-task-management/commons"
)

type TaskRepositoryAdapter struct {
	commons.TaskRepositoryInterface
}

func (a *TaskRepositoryAdapter) GetByID(id string) (commons.Task, error) {
	return a.TaskRepositoryInterface.GetByID(id)
}

func (a *TaskRepositoryAdapter) GetByUserID(userID string) ([]commons.Task, error) {
	return a.TaskRepositoryInterface.GetByUserID(userID)
}

func (a *TaskRepositoryAdapter) Create(task commons.Task) (commons.Task, error) {
	return a.TaskRepositoryInterface.Create(task)
}

func (a *TaskRepositoryAdapter) Update(task commons.Task) error {
	return a.TaskRepositoryInterface.Update(task)
}

func (a *TaskRepositoryAdapter) Delete(id string) error {
	return a.TaskRepositoryInterface.Delete(id)
}
