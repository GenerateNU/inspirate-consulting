package personalCollegeApplicationRepository

import (
	"embed"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PersonalCollegeApplicationRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlPersonalCollegeApplicationFiles embed.FS

func NewPersonalCollegeApplicationRepository(db *pgxpool.Pool) *PersonalCollegeApplicationRepository {
	return &PersonalCollegeApplicationRepository{db: db}
}
