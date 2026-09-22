package extracurricularRepository

import "github.com/jackc/pgx/v5/pgxpool"

type ExtracurricularRepository struct {
	db *pgxpool.Pool
}

func NewExtracurricularRepository(db *pgxpool.Pool) *ExtracurricularRepository {
	return &ExtracurricularRepository{db: db}
}
