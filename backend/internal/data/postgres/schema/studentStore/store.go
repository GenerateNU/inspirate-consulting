package studentstore

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlStudentFiles embed.FS

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
	return &StudentRepository{db: db}
}
