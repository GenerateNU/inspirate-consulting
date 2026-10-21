package counselorRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CounselorRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlCounselorFiles embed.FS

func NewCounselorRepository(db *pgxpool.Pool) *CounselorRepository {
	return &CounselorRepository{db: db}
}
