package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wralith/first-class-func-structure/todo"
)

type UpsertTodoFn func(context.Context, *todo.Todo) (*todo.Todo, error)
type ListTodosFn func(context.Context, todo.TodoQuery) ([]todo.Todo, error)

func UpsertTodo(upsert UpsertTodoFn) func(context.Context, *todo.UpsertTodoInput) (*todo.Todo, error) {
	return func(ctx context.Context, input *todo.UpsertTodoInput) (*todo.Todo, error) {
		t := todo.NewTodo(*input)
		log.Info().Str("title", t.Title).Msg("Creating new todo")
		return upsert(ctx, t)
	}
}

func ToggleTodoCompleted(list ListTodosFn, upsert UpsertTodoFn) func(context.Context, uuid.UUID) (*todo.Todo, error) {
	return func(ctx context.Context, id uuid.UUID) (*todo.Todo, error) {
		query := todo.NewTodoQuery()
		query.ID = id
		l, err := list(ctx, query)
		if err != nil {
			return nil, err
		}
		if len(l) == 0 {
			return nil, todo.ErrTodoNotFound
		}
		t := l[0]
		t.ToggleCompleted()
		return upsert(ctx, &t)
	}
}

func ListTodos(list ListTodosFn) func(context.Context, todo.TodoQuery) ([]todo.Todo, error) {
	return func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
		return list(ctx, input)
	}
}
