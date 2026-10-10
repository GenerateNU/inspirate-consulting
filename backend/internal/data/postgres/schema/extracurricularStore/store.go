package extracurricularRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ExtracurricularRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlExtracurricularFiles embed.FS

func NewExtracurricularRepository(db *pgxpool.Pool) *ExtracurricularRepository {
	return &ExtracurricularRepository{db: db} // Constructor present
}
