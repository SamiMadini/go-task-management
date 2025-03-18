package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	commons "sama/go-task-management/commons"
	"time"

	"github.com/google/uuid"
)

const (
    notificationTimeout = 20 * time.Second
)

var (
    ErrInvalidRequest = errors.New("invalid request payload")
    ErrTaskNotFound   = errors.New("task not found")
)

func parseDueDate(dueDateStr string, defaultTime time.Time) time.Time {
    if dueDateStr == "" {
        return defaultTime
    }

    t, err := time.Parse(time.RFC3339, dueDateStr)
    if err != nil {
        log.Printf("Error parsing due date: %v", err)
        return defaultTime
    }

    return t
}

func createTaskSystemEvent(taskId, correlationId, origin, action, message, jsonData string) commons.TaskSystemEvent {
    return commons.TaskSystemEvent{
        ID:            uuid.New().String(),
        TaskId:        taskId,
        CorrelationId: correlationId,
        Origin:        origin,
        Action:        action,
        Message:       message,
        JsonData:      jsonData,
        EmitAt:        time.Now(),
        CreatedAt:     time.Now(),
    }
}

func marshallJson(obj any) string {
    jsonData, err := json.Marshal(obj)
    if err != nil {
        return "{}"
    }
    return string(jsonData)
}

type operationResult struct {
    operation string
    err       error
}

func handleOperationErrors(operations []func() error) error {
    errChan := make(chan operationResult, len(operations))

    for i, op := range operations {
        go func(i int, op func() error) {
            if err := op(); err != nil {
                errChan <- operationResult{fmt.Sprintf("operation_%d", i), err}
            }
        }(i, op)
    }

    close(errChan)
    for result := range errChan {
        if result.err != nil {
            return result.err
        }
    }
    return nil
}
