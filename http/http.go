package http

import (
	"context"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/wralith/first-class-func-structure/todo"
)

// TODO: All unhandled errors returns internal with below message,
// Check if any error type should be handled differently to create meaningful error messages

type UpsertTodoFn func(context.Context, *todo.UpsertTodoInput) (*todo.Todo, error)
type ListTodoFn func(context.Context, todo.TodoQuery) ([]todo.Todo, error)

// CreateTodo godoc
//
//	@Summary	Create New Todo
//	@Tags		Todos
//	@Produce	json
//	@Param		input	body		todo.UpsertTodoInput	true	"Input"
//	@Success	200		{object}	todo.Todo
//	@Failure	400		{object}	ErrorResponse
//	@Failure	500		{object}	ErrorResponse
//	@Router		/todos [post]
func CreateTodo(upsert UpsertTodoFn) fiber.Handler {
	return func(c fiber.Ctx) error {
		var input todo.UpsertTodoInput
		err := c.Bind().WithAutoHandling().Body(&input)
		if err != nil {
			return err
		}

		t, err := upsert(c.Context(), &input)
		if err != nil {
			return ErrInternalServer
		}

		return c.JSON(t)
	}
}

// ListTodos godoc
//
//	@Summary	List Todos
//	@Tags		Todos
//	@Produce	json
//	@Param		title						query		string	false	"Title"
//	@Param		completed					query		boolean	false	"Completed"
//	@Param		category					query		string	false	"Category"
//	@Param		completed_after_deadline	query		boolean	false	"Completed After Deadline"
//	@Param		size						query		integer	false	"Size"
//	@Param		page						query		integer	false	"Page"
//	@Success	200							{object}	[]todo.Todo
//	@Failure	400							{object}	ErrorResponse
//	@Failure	500							{object}	ErrorResponse
//	@Router		/todos [get]
func ListTodos(list ListTodoFn) fiber.Handler {
	return func(c fiber.Ctx) error {
		input := todo.NewTodoQuery()
		err := c.Bind().WithAutoHandling().Query(&input)
		if err != nil {
			return err
		}

		todos, err := list(c.Context(), input)
		if err != nil {
			return ErrInternalServer
		}

		return c.JSON(todos)
	}
}

// GetSingleTodo godoc
//
//	@Summary	Get Single Todo By ID
//	@Tags		Todos
//	@Produce	json
//	@Param		id	path		string	true	"Todo ID"
//	@Success	200	{object}	todo.Todo
//	@Failure	400	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/todos/{id} [get]
func GetSingleTodo(list ListTodoFn) fiber.Handler {
	return func(c fiber.Ctx) error {
		input := todo.NewTodoQuery()
		err := c.Bind().WithAutoHandling().URI(&input)
		if err != nil {
			return err
		}

		log.Info().Str("id", input.ID.String()).Msg("Getting single todo")
		log.Info().Uint64("size", input.Size).Msg("Getting single todo")

		todos, err := list(c.Context(), input)
		if err != nil {
			return ErrInternalServer
		}

		if len(todos) == 0 {
			return fiber.NewError(fiber.StatusNotFound, "todo not found with id: "+input.ID.String())
		}

		return c.JSON(todos[0])
	}
}

type ToggleTodoCompletedFn func(context.Context, uuid.UUID) (*todo.Todo, error)

// ToggleTodoCompleted godoc
//
//	@Summary	Toggle Todo By ID
//	@Tags		Todos
//	@Produce	json
//	@Param		id	path		string	true	"Todo ID"
//	@Success	200	{object}	todo.Todo
//	@Failure	400	{object}	ErrorResponse
//	@Failure	404	{object}	ErrorResponse
//	@Failure	500	{object}	ErrorResponse
//	@Router		/todos/{id} [patch]
func ToggleTodoCompleted(toggle ToggleTodoCompletedFn) fiber.Handler {
	return func(c fiber.Ctx) error {
		input := todo.NewTodoQuery()
		err := c.Bind().WithAutoHandling().URI(&input)
		if err != nil {
			return err
		}

		t, err := toggle(c.Context(), input.ID)
		if err != nil {
			return err
		}

		return c.JSON(t)
	}
}
