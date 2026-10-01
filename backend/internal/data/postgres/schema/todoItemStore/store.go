package todoItemRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoItemRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlTodoItemFiles embed.FS

func NewTodoItemRepository(db *pgxpool.Pool) *TodoItemRepository {
	return &TodoItemRepository{db: db}
}
