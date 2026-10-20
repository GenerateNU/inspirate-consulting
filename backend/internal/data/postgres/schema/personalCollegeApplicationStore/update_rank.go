package personalCollegeApplicationRepository

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type appRank struct {
	ID   int64
	Rank *int
}

type rankUpdate struct {
	ID   int64
	Rank *int
}

// sets or clears an application's rank
// shifts the ranks of other applications as needed to maintain a consecutive sequence of ranks (no gaps, no duplicates)
// and preserve the relative order of all other applications
func (r *PersonalCollegeApplicationRepository) UpdateApplicationRank(
	ctx context.Context,
	studentID uuid.UUID,
	applicationID int64,
	newRank *int,
) ([]models.PersonalCollegeApplication, error) {
	// Run as a transaction (ensure there are no half-updates if the process fails mid-way)
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	// fetch applications for this student, ordered by rank, and lock them for update
	selectForRankUpdateQuery, err := schema.ReadSQLBaseScript("select_personal_college_application_for_rank_update.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, selectForRankUpdateQuery, studentID)
	if err != nil {
		return nil, err
	}

	applications, err := pgx.CollectRows(rows, pgx.RowToStructByPos[appRank])
	if err != nil {
		return nil, err
	}

	// ensure given application ID exists and get its current rank (if any)
	applicationsByID := make(map[int64]*int, len(applications))
	for _, app := range applications {
		applicationsByID[app.ID] = app.Rank
	}

	oldRank, found := applicationsByID[applicationID]
	if !found {
		return nil, errs.NotFound("personal college application", "id", applicationID)
	}

	// compute the new rank for every application
	// returns: set of (app_id, new rank)
	updates := computeRankShift(applications, applicationID, oldRank, newRank)

	updateRankQuery, err := schema.ReadSQLBaseScript("update_personal_college_application_rank.sql", SqlPersonalCollegeApplicationFiles)
	if err != nil {
		return nil, err
	}
	for _, update := range updates {
		if _, err := tx.Exec(ctx, updateRankQuery, update.Rank, update.ID); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.ListPersonalCollegeApplicationsByStudentID(ctx, studentID)
}

// returns the full set of (app_id, new rank) pairs needed to move applicationID to newRank,
// keeping every other ranked application's relative order intact
// and the overall sequence without gaps or duplicates
func computeRankShift(applications []appRank, applicationID int64, oldRank *int, newRank *int) []rankUpdate {

	// case 1: ranked -> unranked
	if newRank == nil {
		return rankedToUnranked(applications, applicationID, oldRank)
	}

	// case 2: unranked -> ranked
	if oldRank == nil {
		return unrankedToRanked(applications, applicationID, *newRank)
	}

	// case 3: higher rank -> lower rank (moving up)
	if *newRank < *oldRank {
		return moveUp(applications, applicationID, *oldRank, *newRank)
	}

	// case 4: lower rank -> higher rank (moving down)
	return moveDown(applications, applicationID, *oldRank, *newRank)
}

// everything ranked below the old position shifts up by one to close the gap left behind
func rankedToUnranked(applications []appRank, applicationID int64, oldRank *int) []rankUpdate {
	var updates []rankUpdate
	if oldRank != nil {
		old := *oldRank
		for _, app := range applications {
			if app.ID == applicationID {
				continue
			}
			if app.Rank != nil && *app.Rank > old {
				shifted := *app.Rank - 1
				updates = append(updates, rankUpdate{ID: app.ID, Rank: &shifted})
			}
		}
	}
	return append(updates, rankUpdate{ID: applicationID, Rank: nil})
}

// everything at or below the target shifts down by one to make room for the newly-ranked application
func unrankedToRanked(applications []appRank, applicationID int64, target int) []rankUpdate {
	var updates []rankUpdate
	for _, app := range applications {
		if app.ID == applicationID {
			continue
		}
		if app.Rank != nil && *app.Rank >= target {
			shifted := *app.Rank + 1
			updates = append(updates, rankUpdate{ID: app.ID, Rank: &shifted})
		}
	}
	return append(updates, rankUpdate{ID: applicationID, Rank: &target})
}

// moving to a lower numeric rank (target < old) = everything between the new and old position shifts down by one
func moveUp(applications []appRank, applicationID int64, old int, target int) []rankUpdate {
	var updates []rankUpdate
	for _, app := range applications {
		if app.ID == applicationID {
			continue
		}
		if app.Rank != nil && *app.Rank >= target && *app.Rank < old {
			shifted := *app.Rank + 1
			updates = append(updates, rankUpdate{ID: app.ID, Rank: &shifted})
		}
	}
	return append(updates, rankUpdate{ID: applicationID, Rank: &target})
}

// moving to a higher numeric rank (target > old) = everything between the old and new position shifts up by one
func moveDown(applications []appRank, applicationID int64, old int, target int) []rankUpdate {
	var updates []rankUpdate
	for _, app := range applications {
		if app.ID == applicationID {
			continue
		}
		if app.Rank != nil && *app.Rank > old && *app.Rank <= target {
			shifted := *app.Rank - 1
			updates = append(updates, rankUpdate{ID: app.ID, Rank: &shifted})
		}
	}
	return append(updates, rankUpdate{ID: applicationID, Rank: &target})
}
