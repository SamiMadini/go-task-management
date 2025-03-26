package auth

import (
	"sama/go-task-management/commons"
)

func ToUserResponse(user commons.User) UserResponse {
	return UserResponse{
		ID:     user.ID,
		Handle: user.Handle,
		Email:  user.Email,
		Status: user.Status,
	}
}
