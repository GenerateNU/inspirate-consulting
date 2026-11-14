package extracurricularRepository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// setupTestRepo creates a fresh, isolated test database (via testcontainers)
// and returns a repository connected to it. Safe for t.Parallel().
func setupTestRepo(t *testing.T) (*ExtracurricularRepository, *pgxpool.Pool) {
	t.Helper()

	db := testutils.SetupTestDB(t)
	return NewExtracurricularRepository(db), db
}

// createTestUser inserts a user directly, since extracurriculars.user_id
// references public.users, and returns its id.
func createTestUser(t *testing.T, db *pgxpool.Pool) string {
	t.Helper()

	var userID uuid.UUID
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.users (name, email) VALUES ($1, $2) RETURNING id`,
		"extracurricular store test user", uuid.NewString()+"@test.com").Scan(&userID); err != nil {
		t.Fatalf("setup createTestUser failed: %v", err)
	}
	return userID.String()
}

// minimalCreateInput returns a create payload with only the required fields set.
func minimalCreateInput(name, startDate string) models.CreateExtracurricularInput {
	return models.CreateExtracurricularInput{
		Body: models.CreateExtracurricularRequest{
			Name:         name,
			Status:       models.ExtracurricularStatusDoing,
			Type:         models.ExtracurricularTypeInvestment,
			StartDate:    startDate,
			Organization: "Robotics Club",
		},
	}
}

// createTestExtracurricular inserts an extracurricular through the repository and returns it.
func createTestExtracurricular(t *testing.T, repo *ExtracurricularRepository, studentID, userID string, input models.CreateExtracurricularInput) models.Extracurricular {
	t.Helper()

	created, err := repo.CreateExtracurricular(context.Background(), studentID, userID, input)
	if err != nil {
		t.Fatalf("setup CreateExtracurricular failed: %v", err)
	}
	return created.Body
}

func ptr[T any](v T) *T { return &v }

// insertAchievement adds an achievement row directly, since the store has no
// create method for achievements, and returns its id.
func insertAchievement(t *testing.T, db *pgxpool.Pool, extracurricularID int64, entryDate, body string) int64 {
	t.Helper()

	var id int64
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.extracurricular_achievements (extracurricular_id, entry_date, body)
		 VALUES ($1, $2, $3) RETURNING id`,
		extracurricularID, entryDate, body).Scan(&id); err != nil {
		t.Fatalf("setup insertAchievement failed: %v", err)
	}
	return id
}

// insertNote adds a note row directly with an explicit created_at so ordering
// is deterministic, and returns its id.
func insertNote(t *testing.T, db *pgxpool.Pool, extracurricularID int64, createdAt, body string) int64 {
	t.Helper()

	var id int64
	if err := db.QueryRow(context.Background(),
		`INSERT INTO public.extracurricular_notes (extracurricular_id, body, created_at)
		 VALUES ($1, $2, $3) RETURNING id`,
		extracurricularID, body, createdAt).Scan(&id); err != nil {
		t.Fatalf("setup insertNote failed: %v", err)
	}
	return id
}

// countRows returns how many rows in table reference the given extracurricular.
func countRows(t *testing.T, db *pgxpool.Pool, table string, extracurricularID int64) int {
	t.Helper()

	column := "extracurricular_id"
	if table == "extracurriculars" {
		column = "id"
	}
	var n int
	if err := db.QueryRow(context.Background(),
		`SELECT count(*) FROM public.`+table+` WHERE `+column+` = $1`, extracurricularID).Scan(&n); err != nil {
		t.Fatalf("countRows(%s) failed: %v", table, err)
	}
	return n
}
