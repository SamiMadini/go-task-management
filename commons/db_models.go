package commons

import (
	"time"
)

// DBTask represents the database model for tasks
type DBTask struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	Priority    int       `db:"priority" json:"priority"`
	DueDate     time.Time `db:"due_date" json:"due_date"`
	CreatorID   string    `db:"creator_id" json:"creator_id"`
	AssigneeID  *string   `db:"assignee_id" json:"assignee_id,omitempty"`
	EmailSent   bool      `db:"email_sent" json:"email_sent"`
	InAppSent   bool      `db:"in_app_sent" json:"in_app_sent"`
	Deleted     bool      `db:"deleted" json:"deleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
	Events      []TaskSystemEvent `db:"events" json:"events,omitempty"`
}

// DBUser represents the database model for users
type DBUser struct {
	ID             string    `db:"id" json:"id"`
	Handle         string    `db:"handle" json:"handle"`
	Email          string    `db:"email" json:"email"`
	HashedPassword string    `db:"password_hash" json:"-"`
	Salt           string    `db:"salt" json:"-"`
	Status         string    `db:"status" json:"status"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// DBTaskSystemEvent represents the database model for task system events
type DBTaskSystemEvent struct {
	ID            string    `db:"id" json:"id"`
	TaskId        string    `db:"task_id" json:"task_id"`
	CorrelationId string    `db:"correlation_id" json:"correlation_id"`
	Origin        string    `db:"origin" json:"origin"`
	Action        string    `db:"action" json:"action"`
	Message       string    `db:"message" json:"message"`
	JsonData      string    `db:"json_data" json:"json_data"`
	EmitAt        time.Time `db:"emit_at" json:"emit_at"`
	CreatedAt     time.Time `db:"created_at" json:"created_at"`
}

// DBPasswordResetToken represents the database model for password reset tokens
type DBPasswordResetToken struct {
	ID        string    `db:"id" json:"id"`
	UserID    string    `db:"user_id" json:"user_id"`
	Token     string    `db:"token" json:"token"`
	ExpiresAt time.Time `db:"expires_at" json:"expires_at"`
	Used      bool      `db:"used" json:"used"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

// DBInAppNotification represents the database model for in-app notifications
type DBInAppNotification struct {
	ID          string     `db:"id" json:"id"`
	UserID      string     `db:"user_id" json:"user_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	IsRead      bool       `db:"is_read" json:"is_read"`
	ReadAt      *time.Time `db:"read_at" json:"read_at,omitempty"`
	Deleted     bool       `db:"deleted" json:"deleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
}

// ToTask converts a DBTask to a domain Task
func (dt *DBTask) ToTask() Task {
	task := Task{
		ID:          dt.ID,
		Title:       dt.Title,
		Description: dt.Description,
		Status:      dt.Status,
		Priority:    dt.Priority,
		DueDate:     dt.DueDate,
		CreatorID:   dt.CreatorID,
		AssigneeID:  dt.AssigneeID,
		EmailSent:   dt.EmailSent,
		InAppSent:   dt.InAppSent,
		Deleted:     dt.Deleted,
		DeletedAt:   dt.DeletedAt,
		CreatedAt:   dt.CreatedAt,
		UpdatedAt:   dt.UpdatedAt,
		Events:      dt.Events,
	}
	return task
}

// FromTask converts a domain Task to a DBTask
func (dt *DBTask) FromTask(t Task) {
	dt.ID = t.ID
	dt.Title = t.Title
	dt.Description = t.Description
	dt.Status = t.Status
	dt.Priority = t.Priority
	dt.DueDate = t.DueDate
	dt.CreatorID = t.CreatorID
	dt.AssigneeID = t.AssigneeID
	dt.EmailSent = t.EmailSent
	dt.InAppSent = t.InAppSent
	dt.Deleted = t.Deleted
	dt.DeletedAt = t.DeletedAt
	dt.CreatedAt = t.CreatedAt
	dt.UpdatedAt = t.UpdatedAt
	dt.Events = t.Events
}

// ToUser converts a DBUser to a domain User
func (du *DBUser) ToUser() User {
	return User{
		ID:             du.ID,
		Handle:         du.Handle,
		Email:          du.Email,
		HashedPassword: du.HashedPassword,
		Salt:           du.Salt,
		Status:         du.Status,
		CreatedAt:      du.CreatedAt,
		UpdatedAt:      du.UpdatedAt,
	}
}

// FromUser converts a domain User to a DBUser
func (du *DBUser) FromUser(u User) {
	du.ID = u.ID
	du.Handle = u.Handle
	du.Email = u.Email
	du.HashedPassword = u.HashedPassword
	du.Salt = u.Salt
	du.Status = u.Status
	du.CreatedAt = u.CreatedAt
	du.UpdatedAt = u.UpdatedAt
}

// ToTaskSystemEvent converts a DBTaskSystemEvent to a domain TaskSystemEvent
func (de *DBTaskSystemEvent) ToTaskSystemEvent() TaskSystemEvent {
	return TaskSystemEvent{
		ID:            de.ID,
		TaskId:        de.TaskId,
		CorrelationId: de.CorrelationId,
		Origin:        de.Origin,
		Action:        de.Action,
		Message:       de.Message,
		JsonData:      de.JsonData,
		EmitAt:        de.EmitAt,
		CreatedAt:     de.CreatedAt,
	}
}

// FromTaskSystemEvent converts a domain TaskSystemEvent to a DBTaskSystemEvent
func (de *DBTaskSystemEvent) FromTaskSystemEvent(e TaskSystemEvent) {
	de.ID = e.ID
	de.TaskId = e.TaskId
	de.CorrelationId = e.CorrelationId
	de.Origin = e.Origin
	de.Action = e.Action
	de.Message = e.Message
	de.JsonData = e.JsonData
	de.EmitAt = e.EmitAt
	de.CreatedAt = e.CreatedAt
}

// ToPasswordResetToken converts a DBPasswordResetToken to a domain PasswordResetToken
func (dp *DBPasswordResetToken) ToPasswordResetToken() PasswordResetToken {
	return PasswordResetToken{
		ID:        dp.ID,
		UserID:    dp.UserID,
		Token:     dp.Token,
		ExpiresAt: dp.ExpiresAt,
		Used:      dp.Used,
		CreatedAt: dp.CreatedAt,
	}
}

// FromPasswordResetToken converts a domain PasswordResetToken to a DBPasswordResetToken
func (dp *DBPasswordResetToken) FromPasswordResetToken(p PasswordResetToken) {
	dp.ID = p.ID
	dp.UserID = p.UserID
	dp.Token = p.Token
	dp.ExpiresAt = p.ExpiresAt
	dp.Used = p.Used
	dp.CreatedAt = p.CreatedAt
}

// ToInAppNotification converts a DBInAppNotification to a domain InAppNotification
func (d *DBInAppNotification) ToInAppNotification() InAppNotification {
	return InAppNotification{
		ID:          d.ID,
		UserID:      d.UserID,
		Title:       d.Title,
		Description: d.Description,
		IsRead:      d.IsRead,
		ReadAt:      d.ReadAt,
		Deleted:     d.Deleted,
		DeletedAt:   d.DeletedAt,
		UpdatedAt:   d.UpdatedAt,
		CreatedAt:   d.CreatedAt,
	}
}

// FromInAppNotification converts a domain InAppNotification to a DBInAppNotification
func (d *DBInAppNotification) FromInAppNotification(n InAppNotification) {
	d.ID = n.ID
	d.UserID = n.UserID
	d.Title = n.Title
	d.Description = n.Description
	d.IsRead = n.IsRead
	d.ReadAt = n.ReadAt
	d.Deleted = n.Deleted
	d.DeletedAt = n.DeletedAt
	d.UpdatedAt = n.UpdatedAt
	d.CreatedAt = n.CreatedAt
}
