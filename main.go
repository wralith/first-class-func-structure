package main

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/wralith/first-class-func-structure/docs"
	"github.com/wralith/first-class-func-structure/http"
	"github.com/wralith/first-class-func-structure/repo"
	"github.com/wralith/first-class-func-structure/service"
)

// TODO: Use migration tool
const migration = `
CREATE TABLE IF NOT EXISTS todos (
	id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	completed BOOLEAN NOT NULL,
	category TEXT,
	created_at TIMESTAMPTZ NOT NULL,
	updated_at TIMESTAMPTZ NOT NULL,
	finished_at TIMESTAMPTZ,
	deadline_at TIMESTAMPTZ,
	completed_after_deadline BOOLEAN
);
`

//	@title			First Class Fuc Structure
//	@version		1.0
//	@description	Swagger Documentation.

// @host	localhost:8000
func main() {
	// TODO: From config, read it from file or env
	pool, err := pgxpool.New(context.Background(), "postgres://test:password@localhost:5432/func")
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(context.Background(), migration)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		ErrorHandler:    http.ErrorHandler,
		StructValidator: &http.StructValidator{Validator: validator.New()},
	})

	app.Use(logger.New())

	app.Use(cors.New())
	app.Get("/docs/spec", static.New("docs/swagger.yaml"))
	app.Get("/docs", static.New("docs/docs.html"))

	app.Group("/todos").
		Post("/", http.CreateTodo(service.UpsertTodo(repo.InsertTodo(pool)))).
		Get("/:id", http.GetSingleTodo(service.ListTodos(repo.ListTodos(pool)))).
		Get("/", http.ListTodos(service.ListTodos(repo.ListTodos(pool)))).
		Patch("/:id", http.ToggleTodoCompleted(service.ToggleTodoCompleted(repo.ListTodos(pool), repo.UpdateTodo(pool))))

	app.Listen(":8000") // TODO: Add graceful shutdown
}
