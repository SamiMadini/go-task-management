package commons

import (
	"github.com/google/uuid"
)

func IsValidUUID(uuidStr string) bool {
	_, err := uuid.Parse(uuidStr)
	return err == nil
}
