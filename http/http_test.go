package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wralith/first-class-func-structure/todo"
)

func TestCreateTodo(t *testing.T) {
	t.Parallel()
	t.Run("200", func(t *testing.T) {
		handler := CreateTodo(func(ctx context.Context, input *todo.UpsertTodoInput) (*todo.Todo, error) {
			return todo.NewTodo(*input), nil
		})
		input := todo.UpsertTodoInput{Title: "test", Category: "test", DeadlineAt: time.Now().Add(time.Hour)}
		fiberTestWrapper(t, handler, input, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 200, r.StatusCode)
		})
	})

	t.Run("400", func(t *testing.T) {
		input := todo.UpsertTodoInput{Category: "test", DeadlineAt: time.Now().Add(time.Hour)}
		h := CreateTodo(func(ctx context.Context, input *todo.UpsertTodoInput) (*todo.Todo, error) {
			return todo.NewTodo(*input), nil
		})
		fiberTestWrapper(t, h, input, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 400, r.StatusCode)
		})
	})

	t.Run("500", func(t *testing.T) {
		input := todo.UpsertTodoInput{Title: "test", Category: "test", DeadlineAt: time.Now().Add(time.Hour)}
		h := CreateTodo(func(ctx context.Context, input *todo.UpsertTodoInput) (*todo.Todo, error) {
			return nil, errors.New("test")
		})
		fiberTestWrapper(t, h, input, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 500, r.StatusCode)
		})
	})
}

func TestListTodos(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		item1 := todo.NewTodo(todo.UpsertTodoInput{Title: "test1", Category: "test", DeadlineAt: time.Now().Add(time.Hour)})
		item2 := todo.NewTodo(todo.UpsertTodoInput{Title: "test2", Category: "test", DeadlineAt: time.Now().Add(time.Hour)})
		items := []todo.Todo{*item1, *item2}
		query := todo.TodoQuery{Category: "test"}
		h := ListTodos(func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
			return items, nil
		})
		fiberTestWrapper[[]todo.Todo](t, h, query, func(r *http.Response, b []todo.Todo, err error) {
			require.NoError(t, err)
			require.Equal(t, 200, r.StatusCode)
			require.Len(t, b, 2)
		})
	})

	t.Run("500", func(t *testing.T) {
		query := todo.TodoQuery{Category: "test"}
		h := ListTodos(func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
			return nil, errors.New("test")
		})
		fiberTestWrapper(t, h, query, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 500, r.StatusCode)
		})
	})
}

func TestGetSingleTodo(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		item := todo.NewTodo(todo.UpsertTodoInput{Title: "test1", Category: "test", DeadlineAt: time.Now().Add(time.Hour)})
		query := todo.TodoQuery{ID: item.ID}
		h := GetSingleTodo(func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
			return []todo.Todo{*item}, nil
		})
		fiberTestWrapper[todo.Todo](t, h, query, func(r *http.Response, b todo.Todo, err error) {
			require.NoError(t, err)
			require.Equal(t, 200, r.StatusCode)
			require.Equal(t, item.ID, b.ID)
		})
	})

	t.Run("404", func(t *testing.T) {
		query := todo.TodoQuery{ID: uuid.New()}
		h := GetSingleTodo(func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
			return nil, nil
		})
		fiberTestWrapper(t, h, query, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 404, r.StatusCode)
		})
	})

	t.Run("500", func(t *testing.T) {
		item := todo.NewTodo(todo.UpsertTodoInput{Title: "test1", Category: "test", DeadlineAt: time.Now().Add(time.Hour)})
		query := todo.TodoQuery{ID: item.ID}
		h := GetSingleTodo(func(ctx context.Context, input todo.TodoQuery) ([]todo.Todo, error) {
			return nil, errors.New("test")
		})
		fiberTestWrapper(t, h, query, func(r *http.Response, b any, err error) {
			require.NoError(t, err)
			require.Equal(t, 500, r.StatusCode)
		})
	})
}

func TestToggleTodoCompleted(t *testing.T) {
	t.Run("200", func(t *testing.T) {
		item := todo.NewTodo(todo.UpsertTodoInput{Title: "test1", Category: "test", DeadlineAt: time.Now().Add(time.Hour)})
		h := ToggleTodoCompleted(func(ctx context.Context, u uuid.UUID) (*todo.Todo, error) {
			item.ToggleCompleted()
			return item, nil
		})
		fiberTestWrapper[todo.Todo](t, h, uuid.New(), func(r *http.Response, b todo.Todo, err error) {
			require.NoError(t, err)
			require.Equal(t, 200, r.StatusCode)
			require.True(t, b.Completed)
		})
	})
}

func fiberTestWrapper[T any](t *testing.T, h fiber.Handler, input any, assertions func(*http.Response, T, error)) {
	app := fiber.New(fiber.Config{
		StructValidator: &StructValidator{Validator: validator.New()},
		ErrorHandler:    ErrorHandler,
	})
	app.Post("/", h)
	js, err := json.Marshal(input)
	require.NoError(t, err)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(js))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)

	if res.Body == nil {
		assertions(res, *new(T), err)
		return
	}

	var body T
	b, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	err = json.Unmarshal(b, &body)
	assertions(res, body, err)
}
