package personalCollegeApplicationRepository

import "github.com/jackc/pgx/v5/pgxpool"

type PersonalCollegeApplicationRepository struct {
	db *pgxpool.Pool
}

func NewPersonalCollegeApplicationRepository(db *pgxpool.Pool) *PersonalCollegeApplicationRepository {
	return &PersonalCollegeApplicationRepository{db: db}
}
