package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Status options for a task
const (
	StatusPending  string = "pending"
	StatusStarted  string = "started"
	StatusCanceled string = "canceled"
	StatusFinished string = "finished"
	StatusDeleted  string = "deleted"
)

const (
	nameMaxLength        int = 50
	descriptionMaxLength int = 500
	categoryMaxLength    int = 20
)

// TaskStatus choices
var TaskStatus = [9]string{
	StatusPending,
	StatusStarted,
	StatusCanceled,
	StatusFinished,
	StatusDeleted,
}

// Task defines an element at the TODO list
type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Category    string     `json:"category,omitempty"`
	Status      string     `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Created     *time.Time `json:"created"`
}

// Validate a Task data
func (t *Task) Validate() error {
	if len(t.Name) == 0 {
		return errors.New("Task needs a Name")
	}

	if len(t.Name) > nameMaxLength {
		return errors.Errorf("Task name must be less than %d characters", nameMaxLength)
	}

	if len(t.Description) > descriptionMaxLength {
		return errors.Errorf("Task description must be less than %d characters", descriptionMaxLength)
	}

	if len(t.Category) > categoryMaxLength {
		return errors.Errorf("Task category must be less than %d characters", categoryMaxLength)
	}

	// TODO category show be taken from a list of existing categories

	if t.Status == "" {
		t.Status = StatusPending
	}

	found := false
	status := strings.ToLower(t.Status)
	for _, s := range TaskStatus {
		if s == status {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("Task status should be one of %v", TaskStatus)
	}

	return nil
}
