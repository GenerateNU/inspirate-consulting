package globalCollegeRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GlobalCollegeRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlGlobalCollegeFiles embed.FS

func NewGlobalCollegeRepository(db *pgxpool.Pool) *GlobalCollegeRepository {
	return &GlobalCollegeRepository{db: db}
}
