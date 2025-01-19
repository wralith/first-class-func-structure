package repo

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wralith/first-class-func-structure/todo"
)

var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
var todosTable = "todos"

// If we want to use UUID generator from database, we could add returning id prefix to query, scan the row and return the id.
// Since we have all the data in the struct, we can just return it, why not
func InsertTodo(pool *pgxpool.Pool) func(context.Context, *todo.Todo) (*todo.Todo, error) {
	return func(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
		sql, args, err := sq.
			Insert(todosTable).
			Columns("id", "title", "completed", "category", "created_at", "updated_at", "finished_at", "deadline_at", "completed_after_deadline").
			Values(t.ID, t.Title, t.Completed, t.Category, t.CreatedAt, t.UpdatedAt, t.FinishedAt, t.DeadlineAt, t.CompletedAfterDeadline).
			ToSql()
		if err != nil {
			return t, err
		}

		_, err = pool.Exec(ctx, sql, args...)
		if err != nil {
			return t, err
		}

		return t, nil
	}
}

func UpdateTodo(pool *pgxpool.Pool) func(context.Context, *todo.Todo) (*todo.Todo, error) {
	return func(ctx context.Context, t *todo.Todo) (*todo.Todo, error) {
		sql, args, err := sq.
			Update(todosTable).
			Set("title", t.Title).
			Set("completed", t.Completed).
			Set("category", t.Category).
			Set("created_at", t.CreatedAt).
			Set("updated_at", t.UpdatedAt).
			Set("finished_at", t.FinishedAt).
			Set("deadline_at", t.DeadlineAt).
			Set("completed_after_deadline", t.CompletedAfterDeadline).
			Where(squirrel.Eq{"id": t.ID}).
			ToSql()
		if err != nil {
			return nil, err
		}

		_, err = pool.Exec(ctx, sql, args...)
		if err != nil {
			return nil, err
		}

		return t, nil
	}
}

type TodoQuery todo.TodoQuery

func (t *TodoQuery) Where() squirrel.Eq {
	where := squirrel.Eq{}
	if t.ID != uuid.Nil {
		where["id"] = t.ID
	}
	if t.Title != "" {
		where["title"] = t.Title
	}
	if t.Completed {
		where["completed"] = t.Completed
	}
	if t.Category != "" {
		where["category"] = t.Category
	}
	if t.CompletedAfterDeadline {
		where["completed_after_deadline"] = t.CompletedAfterDeadline
	}

	return where
}

func ListTodos(pool *pgxpool.Pool) func(context.Context, todo.TodoQuery) ([]todo.Todo, error) {
	return func(ctx context.Context, query todo.TodoQuery) ([]todo.Todo, error) {
		q := TodoQuery(query)
		sql, args, err := sq.
			Select("id", "title", "completed", "category", "created_at", "updated_at", "finished_at", "deadline_at", "completed_after_deadline").
			From(todosTable).
			Where(q.Where()).
			Limit(q.Size).
			Offset(q.Size * (query.Page - 1)).
			ToSql()
		if err != nil {
			return nil, err
		}

		rows, err := pool.Query(ctx, sql, args...)
		if err != nil {
			return nil, err
		}

		return pgx.CollectRows[todo.Todo](rows, pgx.RowToStructByNameLax)
	}
}
