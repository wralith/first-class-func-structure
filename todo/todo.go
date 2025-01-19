package todo

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID                     uuid.UUID `json:"id"`
	Title                  string    `json:"title"`
	Completed              bool      `json:"completed"`
	Category               string    `json:"category"`
	CreatedAt              time.Time `json:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt"`
	FinishedAt             time.Time `json:"finishedAt"`
	DeadlineAt             time.Time `json:"deadlineAt"`
	CompletedAfterDeadline bool      `json:"completedAfterDeadline"`
}

type UpsertTodoInput struct {
	Title      string    `json:"title" query:"title" validate:"required"`
	Category   string    `json:"category" query:"category"`
	DeadlineAt time.Time `json:"deadlineAt" query:"deadline"`
}

func NewTodo(input UpsertTodoInput) *Todo {
	return &Todo{
		ID:         uuid.New(),
		Title:      input.Title,
		Category:   input.Category,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		DeadlineAt: input.DeadlineAt,
	}
}

func (t *Todo) Update(input UpsertTodoInput) {
	t.Title = input.Title
	t.Category = input.Category
	t.UpdatedAt = time.Now()
	t.DeadlineAt = input.DeadlineAt
}

func (t *Todo) ToggleCompleted() {
	t.UpdatedAt = time.Now()
	if t.Completed {
		t.Completed = false
		t.FinishedAt = time.Time{}
		t.CompletedAfterDeadline = false
	} else {
		t.Completed = true
		t.FinishedAt = time.Now()
		if !t.DeadlineAt.IsZero() && time.Now().After(t.DeadlineAt) {
			t.CompletedAfterDeadline = true
		}
	}
}

// TODO: Range queries is on you to implement haha
type TodoQuery struct {
	ID                     uuid.UUID `json:"id" param:"id"`
	Title                  string    `json:"title" query:"title"`
	Completed              bool      `json:"completed" query:"completed"`
	Category               string    `json:"category" query:"category"`
	CompletedAfterDeadline bool      `json:"completedAfterDeadline" query:"completedAfterDeadline"`

	Size uint64 `json:"size" query:"size"`
	Page uint64 `json:"page" query:"page"`
}

func NewTodoQuery() TodoQuery {
	return TodoQuery{
		Size: 10,
		Page: 1,
	}
}
