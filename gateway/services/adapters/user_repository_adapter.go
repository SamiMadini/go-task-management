package adapters

import (
	"sama/go-task-management/commons"
)

type UserRepositoryAdapter struct {
	commons.UserRepositoryInterface
}

func (a *UserRepositoryAdapter) Update(user commons.User) error {
	_, err := a.UserRepositoryInterface.UpdatePassword(user.ID, user.HashedPassword, user.Salt)
	return err
}

func (a *UserRepositoryAdapter) UpdatePassword(userID string, hashedPassword string) error {
	_, err := a.UserRepositoryInterface.UpdatePassword(userID, hashedPassword, "")
	return err
}
