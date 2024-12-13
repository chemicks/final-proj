package domain

import "time"

type Task struct {
	Id          uint64
	UserId      uint64
	Title       string
	Description string
	Status      TaskStatus
	Date        time.Time
	CreatedDate time.Time
	UpdatedDate time.Time
	DeletedDate *time.Time
}

type TaskStatus string

const (
	NewTaskStatus       TaskStatus = "NEW"
	ImportantTaskStatus TaskStatus = "IMPORTANT"
	CompleteTaskStatus  TaskStatus = "DONE"
	ExpiredTaskStatus   TaskStatus = "EXPIRED"
)
