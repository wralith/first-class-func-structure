package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wralith/first-class-func-structure/todo"
)

func TestUpsertTodo(t *testing.T) {
	got, err := UpsertTodo(func(ctx context.Context, item *todo.Todo) (*todo.Todo, error) {
		return item, nil
	})(context.Background(), &todo.UpsertTodoInput{
		Title:      "test",
		Category:   "test",
		DeadlineAt: time.Now().Add(time.Hour),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, "test", got.Title)
	require.Equal(t, "test", got.Category)
}

func TestToggleTodoCompleted(t *testing.T) {
	item := todo.NewTodo(todo.UpsertTodoInput{
		Title:      "test",
		Category:   "test",
		DeadlineAt: time.Now().Add(time.Hour),
	})
	item.Completed = false
	item.CompletedAfterDeadline = false

	listFn := func(ctx context.Context, query todo.TodoQuery) ([]todo.Todo, error) {
		return []todo.Todo{*item}, nil
	}
	upsertFn := func(ctx context.Context, item *todo.Todo) (*todo.Todo, error) {
		return item, nil
	}
	got, err := ToggleTodoCompleted(listFn, upsertFn)(context.Background(), uuid.New())

	require.NoError(t, err)
	require.NotNil(t, got)
	require.True(t, got.Completed)
	require.False(t, got.CompletedAfterDeadline)
}
