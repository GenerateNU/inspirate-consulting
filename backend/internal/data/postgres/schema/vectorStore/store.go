package vectorRepository

import (
	"embed"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Vectors are written by the embed-vectors edge function; this store only reads them for search.
type VectorRepository struct {
	db *pgxpool.Pool
}

//go:embed sql/*.sql
var SqlVectorFiles embed.FS

func NewVectorRepository(db *pgxpool.Pool) *VectorRepository {
	return &VectorRepository{db: db}
}

// converts an embedding to pgvector's text format, e.g. "[0.1,0.2,0.3]".
func formatVector(embedding []float32) string {
	parts := make([]string, len(embedding))
	for i, v := range embedding {
		parts[i] = strconv.FormatFloat(float64(v), 'f', -1, 32)
	}
	return "[" + strings.Join(parts, ",") + "]"
}
