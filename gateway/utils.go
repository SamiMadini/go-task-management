package main

import (
	"encoding/json"
	"errors"
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
        log.Printf("Due date is empty, returning default time: %v", defaultTime)
        return defaultTime
    }

    log.Printf("Attempting to parse due date: %s", dueDateStr)

    t, err := time.Parse(time.RFC3339Nano, dueDateStr)
    if err != nil {
        log.Printf("Failed to parse with RFC3339Nano: %v", err)

        t, err = time.Parse(time.RFC3339, dueDateStr)
        if err != nil {
            log.Printf("Failed to parse with RFC3339: %v", err)
            return defaultTime
        }
    }

    log.Printf("Successfully parsed due date: %v", t)
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
