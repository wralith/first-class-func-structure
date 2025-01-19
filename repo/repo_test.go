package repo

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/postgres"
	"github.com/stretchr/testify/require"
	"github.com/wralith/first-class-func-structure/todo"
)

var pool *pgxpool.Pool

const migration = `
CREATE TABLE todos (
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

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		return
	}
	p := postgres.Preset(
		postgres.WithUser("test", "password"),
		postgres.WithDatabase("circle"),
		postgres.WithQueries(migration),
	)

	container, err := gnomock.Start(p)
	if err != nil {
		panic(err)
	}
	defer gnomock.Stop(container)

	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s  dbname=%s sslmode=disable",
		container.Host, container.DefaultPort(), "test", "password", "circle",
	)

	pool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

// You can separate those ofc, lots of works...
func TestCreateUpdateGetFlow(t *testing.T) {
	item := todo.NewTodo(todo.UpsertTodoInput{
		Title:      "Test",
		Category:   "Test",
		DeadlineAt: time.Now().Add(time.Hour),
	})
	ctx := context.Background()
	_, err := InsertTodo(pool)(ctx, item)
	require.NoError(t, err)

	query := todo.NewTodoQuery()
	query.ID = item.ID
	items, err := ListTodos(pool)(ctx, query)
	require.NoError(t, err)
	got := items[0]
	require.Equal(t, item.Title, got.Title)
	require.Equal(t, item.ID, got.ID)

	item.Title = "Updated Title"
	got, err = UpdateTodo(pool)(ctx, item)
	require.NoError(t, err)
	require.Equal(t, item.Title, got.Title)
	require.Equal(t, item.ID, got.ID)

	items, err = ListTodos(pool)(ctx, query)
	require.NoError(t, err)
	got = items[0]
	require.Equal(t, item.Title, got.Title)
	require.Equal(t, item.ID, got.ID)
}
