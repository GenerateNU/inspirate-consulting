package extracurricular

import (
	"context"
	"net/http"
	"testing"

	"inspirate-consulting/internal/auth"
	extracurricularRepository "inspirate-consulting/internal/data/postgres/schema/extracurricularStore"
	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ownershipFixture is a real database with two counselors, each assigned one
// student, and an extracurricular owned by the first student.
type ownershipFixture struct {
	handler             *Handler
	db                  *pgxpool.Pool
	extracurricularID   int64
	owner, otherStudent uuid.UUID
	assignedCounselor   uuid.UUID
	otherCounselor      uuid.UUID
}

func setupOwnershipFixture(t *testing.T) ownershipFixture {
	t.Helper()
	ctx := context.Background()
	db := testutils.SetupTestDB(t)

	insertUser := func() uuid.UUID {
		var id uuid.UUID
		require.NoError(t, db.QueryRow(ctx,
			`INSERT INTO users (name, email) VALUES ('ownership test', $1) RETURNING id`,
			uuid.NewString()+"@test.com").Scan(&id))
		return id
	}
	insertCounselor := func() uuid.UUID {
		var id uuid.UUID
		require.NoError(t, db.QueryRow(ctx, `INSERT INTO counselor (user_id) VALUES ($1) RETURNING id`, insertUser()).Scan(&id))
		return id
	}
	insertStudent := func(counselorID uuid.UUID) uuid.UUID {
		var id uuid.UUID
		require.NoError(t, db.QueryRow(ctx,
			`INSERT INTO student (user_id, year, gpa, review_balance, counselor_id) VALUES ($1, 'senior', 4, 0, $2) RETURNING id`,
			insertUser(), counselorID).Scan(&id))
		return id
	}

	f := ownershipFixture{db: db}
	f.assignedCounselor = insertCounselor()
	f.otherCounselor = insertCounselor()
	f.owner = insertStudent(f.assignedCounselor)
	f.otherStudent = insertStudent(f.otherCounselor)

	repo := extracurricularRepository.NewExtracurricularRepository(db)
	created, err := repo.CreateExtracurricular(ctx, f.owner.String(), insertUser().String(), models.CreateExtracurricularInput{
		Body: models.CreateExtracurricularRequest{
			Name: "Debate", Status: models.ExtracurricularStatusDoing, Type: models.ExtracurricularTypeInvestment,
			StartDate: "2025-01-01", Organization: "Lincoln High",
		},
	})
	require.NoError(t, err)
	f.extracurricularID = created.Body.ID

	f.handler = NewHandler(repo, auth.NewAuthorizer(db))
	return f
}

func asStudent(id uuid.UUID) context.Context { return studentRoleCtx(id) }

func asCounselor(id uuid.UUID) context.Context {
	return context.WithValue(auth.WithCounselorID(context.Background(), id), "Role", "counselor") //nolint:staticcheck // matches how auth middleware stores the role
}

func (f ownershipFixture) count(t *testing.T, table string) int {
	t.Helper()
	var n int
	require.NoError(t, f.db.QueryRow(context.Background(),
		`SELECT count(*) FROM public.`+table+` WHERE extracurricular_id = $1`, f.extracurricularID).Scan(&n))
	return n
}

func requireForbidden(t *testing.T, err error) {
	t.Helper()
	var httpErr errs.HTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, http.StatusForbidden, httpErr.Code)
}

// exerciseChildren creates, edits and deletes a note and an achievement as the
// caller in ctx, returning the first error.
func (f ownershipFixture) exerciseChildren(ctx context.Context) error {
	id := f.extracurricularID

	note, err := f.handler.CreateExtracurricularNote(ctx, &models.CreateExtracurricularNoteInput{
		ExtracurricularID: id, Body: models.ExtracurricularNoteRequest{Body: "note"},
	})
	if err != nil {
		return err
	}
	if _, err := f.handler.UpdateExtracurricularNote(ctx, &models.UpdateExtracurricularNoteInput{
		ExtracurricularID: id, NoteID: note.ID, Body: models.ExtracurricularNoteRequest{Body: "edited"},
	}); err != nil {
		return err
	}
	if err := f.handler.DeleteExtracurricularNote(ctx, &models.DeleteExtracurricularNoteInput{ExtracurricularID: id, NoteID: note.ID}); err != nil {
		return err
	}

	achievement, err := f.handler.CreateExtracurricularAchievement(ctx, &models.CreateExtracurricularAchievementInput{
		ExtracurricularID: id, Body: models.CreateExtracurricularAchievementRequest{EntryDate: "2025-03-01", Body: "won"},
	})
	if err != nil {
		return err
	}
	edited := "won again"
	if _, err := f.handler.UpdateExtracurricularAchievement(ctx, &models.UpdateExtracurricularAchievementInput{
		ExtracurricularID: id, AchievementID: achievement.ID, Body: models.UpdateExtracurricularAchievementRequest{Body: &edited},
	}); err != nil {
		return err
	}
	return f.handler.DeleteExtracurricularAchievement(ctx, &models.DeleteExtracurricularAchievementInput{ExtracurricularID: id, AchievementID: achievement.ID})
}

// End to end against a real database: only the owning student and their
// assigned counselor may create, edit or delete notes and achievements.
func TestOwnership_ExtracurricularChildren(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	t.Parallel()

	f := setupOwnershipFixture(t)

	t.Run("owning student can create, edit and delete", func(t *testing.T) {
		require.NoError(t, f.exerciseChildren(asStudent(f.owner)))
	})

	t.Run("assigned counselor can create, edit and delete", func(t *testing.T) {
		require.NoError(t, f.exerciseChildren(asCounselor(f.assignedCounselor)))
	})

	// Seed one of each so edit/delete attempts have a real target.
	owner := asStudent(f.owner)
	note, err := f.handler.CreateExtracurricularNote(owner, &models.CreateExtracurricularNoteInput{
		ExtracurricularID: f.extracurricularID, Body: models.ExtracurricularNoteRequest{Body: "keep me"},
	})
	require.NoError(t, err)
	achievement, err := f.handler.CreateExtracurricularAchievement(owner, &models.CreateExtracurricularAchievementInput{
		ExtracurricularID: f.extracurricularID, Body: models.CreateExtracurricularAchievementRequest{EntryDate: "2025-04-01", Body: "keep me"},
	})
	require.NoError(t, err)

	outsiders := map[string]context.Context{
		"another student":      asStudent(f.otherStudent),
		"unassigned counselor": asCounselor(f.otherCounselor),
	}
	for name, ctx := range outsiders {
		t.Run(name+" is forbidden on every operation", func(t *testing.T) {
			id := f.extracurricularID
			edited := "hijacked"

			_, err := f.handler.CreateExtracurricularNote(ctx, &models.CreateExtracurricularNoteInput{
				ExtracurricularID: id, Body: models.ExtracurricularNoteRequest{Body: "x"},
			})
			requireForbidden(t, err)
			_, err = f.handler.UpdateExtracurricularNote(ctx, &models.UpdateExtracurricularNoteInput{
				ExtracurricularID: id, NoteID: note.ID, Body: models.ExtracurricularNoteRequest{Body: edited},
			})
			requireForbidden(t, err)
			requireForbidden(t, f.handler.DeleteExtracurricularNote(ctx, &models.DeleteExtracurricularNoteInput{ExtracurricularID: id, NoteID: note.ID}))

			_, err = f.handler.CreateExtracurricularAchievement(ctx, &models.CreateExtracurricularAchievementInput{
				ExtracurricularID: id, Body: models.CreateExtracurricularAchievementRequest{EntryDate: "2025-01-01", Body: "x"},
			})
			requireForbidden(t, err)
			_, err = f.handler.UpdateExtracurricularAchievement(ctx, &models.UpdateExtracurricularAchievementInput{
				ExtracurricularID: id, AchievementID: achievement.ID, Body: models.UpdateExtracurricularAchievementRequest{Body: &edited},
			})
			requireForbidden(t, err)
			requireForbidden(t, f.handler.DeleteExtracurricularAchievement(ctx, &models.DeleteExtracurricularAchievementInput{ExtracurricularID: id, AchievementID: achievement.ID}))
		})
	}

	t.Run("forbidden attempts left the data untouched", func(t *testing.T) {
		assert.Equal(t, 1, f.count(t, "extracurricular_notes"))
		assert.Equal(t, 1, f.count(t, "extracurricular_achievements"))

		detail, err := f.handler.GetExtracurricular(owner, f.extracurricularID)
		require.NoError(t, err)
		assert.Equal(t, "keep me", detail.Notes[0].Body)
		assert.Equal(t, "keep me", detail.Achievements[0].Body)
	})
}
