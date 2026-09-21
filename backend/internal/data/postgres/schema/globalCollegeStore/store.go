package globalCollegeRepository

import "github.com/jackc/pgx/v5/pgxpool"

type GlobalCollegeRepository struct {
	db *pgxpool.Pool
}

func NewGlobalCollegeRepository(db *pgxpool.Pool) *GlobalCollegeRepository {
	return &GlobalCollegeRepository{db: db}
}
