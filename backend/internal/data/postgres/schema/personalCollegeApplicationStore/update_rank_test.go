package personalCollegeApplicationRepository

import (
	"context"
	"fmt"
	"testing"

	testutils "inspirate-consulting/internal/data/postgres/testUtils"
	"inspirate-consulting/internal/models"
)

// rank can be null
func intPtr(i int) *int { return &i }

func TestUpdateApplicationRank_UnrankedToRanked(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeA := createTestGlobalCollege(t, db, "Rank Test Alpha", "Alpha City, AA")
	collegeB := createTestGlobalCollege(t, db, "Rank Test Beta", "Beta City, BB")
	studentID := "00000000-0000-0000-0000-000000000031"

	appA, _ := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeA, ApplicationType: "ED", Category: "reach",
	})
	appB, _ := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeB, ApplicationType: "RD", Category: "safety",
	})

	_, err := repo.UpdateApplicationRank(ctx, studentID, appA.ID, intPtr(1))
	if err != nil {
		t.Fatalf("UpdateApplicationRank(appA, 1) failed: %v", err)
	}

	results, err := repo.UpdateApplicationRank(ctx, studentID, appB.ID, intPtr(1))
	if err != nil {
		t.Fatalf("UpdateApplicationRank(appB, 1) failed: %v", err)
	}

	ranksByID := make(map[int64]*int)
	for _, r := range results {
		ranksByID[r.ID] = r.Rank
	}

	if ranksByID[appB.ID] == nil || *ranksByID[appB.ID] != 1 {
		t.Errorf("expected appB rank 1, got %v", ranksByID[appB.ID])
	}
	if ranksByID[appA.ID] == nil || *ranksByID[appA.ID] != 2 {
		t.Errorf("expected appA shifted to rank 2, got %v", ranksByID[appA.ID])
	}
}

func TestUpdateApplicationRank_MoveUp(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	studentID := "00000000-0000-0000-0000-000000000032"
	ids := make([]int64, 4)
	for i := 0; i < 4; i++ {
		college := createTestGlobalCollege(t, db, fmt.Sprintf("Move Up College %d", i), fmt.Sprintf("City %d, ST", i))
		app, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: college, ApplicationType: "ED", Category: "reach",
		})
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
		ids[i] = app.ID
		if _, err := repo.UpdateApplicationRank(ctx, studentID, app.ID, intPtr(i+1)); err != nil {
			t.Fatalf("setup rank assignment failed: %v", err)
		}
	}

	results, err := repo.UpdateApplicationRank(ctx, studentID, ids[3], intPtr(2))
	if err != nil {
		t.Fatalf("UpdateApplicationRank (move up) failed: %v", err)
	}

	ranksByID := make(map[int64]*int)
	for _, r := range results {
		ranksByID[r.ID] = r.Rank
	}

	assertRank(t, ranksByID, ids[0], 1)
	assertRank(t, ranksByID, ids[3], 2)
	assertRank(t, ranksByID, ids[1], 3)
	assertRank(t, ranksByID, ids[2], 4)
}

func TestUpdateApplicationRank_MoveDown(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	studentID := "00000000-0000-0000-0000-000000000033"
	ids := make([]int64, 4)
	for i := 0; i < 4; i++ {
		college := createTestGlobalCollege(t, db, fmt.Sprintf("Move Down College %d", i), fmt.Sprintf("City %d, ST", i))
		app, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: college, ApplicationType: "ED", Category: "reach",
		})
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
		ids[i] = app.ID
		if _, err := repo.UpdateApplicationRank(ctx, studentID, app.ID, intPtr(i+1)); err != nil {
			t.Fatalf("setup rank assignment failed: %v", err)
		}
	}

	results, err := repo.UpdateApplicationRank(ctx, studentID, ids[0], intPtr(3))
	if err != nil {
		t.Fatalf("UpdateApplicationRank (move down) failed: %v", err)
	}

	ranksByID := make(map[int64]*int)
	for _, r := range results {
		ranksByID[r.ID] = r.Rank
	}

	assertRank(t, ranksByID, ids[1], 1)
	assertRank(t, ranksByID, ids[2], 2)
	assertRank(t, ranksByID, ids[0], 3)
	assertRank(t, ranksByID, ids[3], 4)
}

func TestUpdateApplicationRank_UnrankCompactsRemaining(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	studentID := "00000000-0000-0000-0000-000000000034"
	ids := make([]int64, 4)
	for i := 0; i < 4; i++ {
		college := createTestGlobalCollege(t, db, fmt.Sprintf("Unrank College %d", i), fmt.Sprintf("City %d, ST", i))
		app, err := repo.CreatePersonalCollegeApplication(ctx, studentID, models.CreatePersonalCollegeApplicationRequestBody{
			GlobalCollegeID: college, ApplicationType: "ED", Category: "reach",
		})
		if err != nil {
			t.Fatalf("setup create failed: %v", err)
		}
		ids[i] = app.ID
		if _, err := repo.UpdateApplicationRank(ctx, studentID, app.ID, intPtr(i+1)); err != nil {
			t.Fatalf("setup rank assignment failed: %v", err)
		}
	}

	results, err := repo.UpdateApplicationRank(ctx, studentID, ids[1], nil)
	if err != nil {
		t.Fatalf("UpdateApplicationRank (unrank) failed: %v", err)
	}

	ranksByID := make(map[int64]*int)
	for _, r := range results {
		ranksByID[r.ID] = r.Rank
	}

	assertRank(t, ranksByID, ids[0], 1)
	if ranksByID[ids[1]] != nil {
		t.Errorf("expected ids[1] to be unranked (nil), got %v", ranksByID[ids[1]])
	}
	assertRank(t, ranksByID, ids[2], 2)
	assertRank(t, ranksByID, ids[3], 3)
}

func TestUpdateApplicationRank_NonexistentApplicationFails(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	_, err := repo.UpdateApplicationRank(ctx, "00000000-0000-0000-0000-000000000035", 999999, intPtr(1))
	if err == nil {
		t.Fatal("expected an error ranking a nonexistent application, got nil")
	}
}

func TestUpdateApplicationRank_WrongStudentFails(t *testing.T) {
	t.Parallel()

	db := testutils.SetupTestDB(t)
	repo := NewPersonalCollegeApplicationRepository(db)
	ctx := context.Background()

	collegeID := createTestGlobalCollege(t, db, "Rank Ownership College", "City, ST")
	ownerStudentID := "00000000-0000-0000-0000-000000000036"
	otherStudentID := "00000000-0000-0000-0000-000000000037"

	app, err := repo.CreatePersonalCollegeApplication(ctx, ownerStudentID, models.CreatePersonalCollegeApplicationRequestBody{
		GlobalCollegeID: collegeID, ApplicationType: "ED", Category: "reach",
	})
	if err != nil {
		t.Fatalf("setup create failed: %v", err)
	}

	_, err = repo.UpdateApplicationRank(ctx, otherStudentID, app.ID, intPtr(1))
	if err == nil {
		t.Fatal("expected an error ranking another student's application, got nil")
	}
}

func assertRank(t *testing.T, ranksByID map[int64]*int, id int64, expected int) {
	t.Helper()
	rank := ranksByID[id]
	if rank == nil {
		t.Errorf("expected id=%d to have rank %d, got nil", id, expected)
		return
	}
	if *rank != expected {
		t.Errorf("expected id=%d to have rank %d, got %d", id, expected, *rank)
	}
}