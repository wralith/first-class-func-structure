package todo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestNewTodo(t *testing.T) {
	todo := NewTodo(UpsertTodoInput{
		Title:    "Test",
		Category: "Test",
		Deadline: time.Now().Add(time.Hour),
	})

	require.Equal(t, "Test", todo.Title)
	require.Equal(t, "Test", todo.Category)
	require.False(t, todo.Completed)
	require.False(t, todo.Completed)
	require.WithinDuration(t, time.Now(), todo.CreatedAt, time.Second)
	require.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(time.Hour), todo.DeadlineAt, time.Second)
	require.False(t, todo.DeadlinePassed)
}

func TestUpdate(t *testing.T) {
	todo := NewTodo(UpsertTodoInput{
		Title:    "Test",
		Category: "Test",
		Deadline: time.Now().Add(time.Hour),
	})

	todo.Update(UpsertTodoInput{
		Title:    "Test2",
		Category: "Test2",
		Deadline: time.Now().Add(time.Hour * 2),
	})

	require.Equal(t, "Test2", todo.Title)
	require.Equal(t, "Test2", todo.Category)
	require.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
	require.WithinDuration(t, time.Now().Add(time.Hour*2), todo.DeadlineAt, time.Second)
}

func TestFinish(t *testing.T) {
	t.Parallel()
	t.Run("Deadline passed", func(t *testing.T) {
		todo := NewTodo(UpsertTodoInput{
			Title:    "Test",
			Category: "Test",
			Deadline: time.Now().Add(-time.Hour),
		})

		todo.Finish()

		require.True(t, todo.Completed)
		require.True(t, todo.DeadlinePassed)
		require.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
		require.WithinDuration(t, time.Now(), todo.FinishedAt, time.Second)
	})

	t.Run("Deadline not passed", func(t *testing.T) {
		todo := NewTodo(UpsertTodoInput{
			Title:    "Test",
			Category: "Test",
			Deadline: time.Now().Add(time.Hour),
		})

		todo.Finish()

		require.True(t, todo.Completed)
		require.False(t, todo.DeadlinePassed)
		require.WithinDuration(t, time.Now(), todo.UpdatedAt, time.Second)
		require.WithinDuration(t, time.Now(), todo.FinishedAt, time.Second)
	})
}
