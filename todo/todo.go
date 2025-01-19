package todo

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Completed      bool      `json:"completed"`
	Category       string    `json:"category"`
	CreatedAt      time.Time `json:"created"`
	UpdatedAt      time.Time `json:"updated"`
	FinishedAt     time.Time `json:"finished"`
	DeadlineAt     time.Time `json:"deadline"`
	DeadlinePassed bool      `json:"deadlinePassed"`
}

type UpsertTodoInput struct {
	Title    string    `json:"title" query:"title" validate:"required"`
	Category string    `json:"category" query:"category"`
	Deadline time.Time `json:"deadline" query:"deadline"`
}

func NewTodo(input UpsertTodoInput) *Todo {
	return &Todo{
		ID:         uuid.New(),
		Title:      input.Title,
		Category:   input.Category,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeadlineAt: input.Deadline,
	}
}

func (t *Todo) Update(input UpsertTodoInput) {
	t.Title = input.Title
	t.Category = input.Category
	t.UpdatedAt = time.Now()
	t.DeadlineAt = input.Deadline
}

func (t *Todo) Finish() {
	if !t.DeadlineAt.IsZero() && time.Now().After(t.DeadlineAt) {
		t.DeadlinePassed = true
	}
	t.Completed = true
	t.UpdatedAt = time.Now()
	t.FinishedAt = time.Now()
}
