package globalCollegeRepository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/models"
)

// collegeRow is one global_colleges row, before its cycles are attached.
type collegeRow struct {
	ID             int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
	SchoolName     string
	SchoolLocation string
}

func (c collegeRow) toModel() models.GlobalCollege {
	return models.GlobalCollege{
		ID:             c.ID,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
		SchoolName:     c.SchoolName,
		SchoolLocation: c.SchoolLocation,
		Cycles:         []models.CollegeApplicationCycle{},
	}
}

type cycleRow struct {
	GlobalCollegeID int64
	Cycle           string
	Deadline        *time.Time
	RollingStart    *string
	RollingEnd      *string
}

// collectColleges reads global_colleges rows and attaches each college's application cycles.
func collectColleges(ctx context.Context, db dbinterface.QueryInterface, rows pgx.Rows) ([]models.GlobalCollege, error) {
	collegeRows, err := pgx.CollectRows(rows, pgx.RowToStructByPos[collegeRow])
	if err != nil {
		return nil, err
	}

	colleges := make([]models.GlobalCollege, len(collegeRows))
	for i, row := range collegeRows {
		colleges[i] = row.toModel()
	}

	if err := attachCycles(ctx, db, colleges); err != nil {
		return nil, err
	}
	return colleges, nil
}

// attachCycles loads the application cycles for all given colleges in one query.
func attachCycles(ctx context.Context, db dbinterface.QueryInterface, colleges []models.GlobalCollege) error {
	if len(colleges) == 0 {
		return nil
	}

	listCyclesQuery, err := schema.ReadSQLBaseScript("list_college_application_cycles.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return err
	}

	ids := make([]int64, len(colleges))
	indexByID := make(map[int64]int, len(colleges))
	for i, college := range colleges {
		ids[i] = college.ID
		indexByID[college.ID] = i
	}

	rows, err := db.Query(ctx, listCyclesQuery, ids)
	if err != nil {
		return err
	}
	cycleRows, err := pgx.CollectRows(rows, pgx.RowToStructByPos[cycleRow])
	if err != nil {
		return err
	}

	for _, row := range cycleRows {
		i := indexByID[row.GlobalCollegeID]
		colleges[i].Cycles = append(colleges[i].Cycles, models.CollegeApplicationCycle{
			Cycle:        row.Cycle,
			Deadline:     row.Deadline,
			RollingStart: row.RollingStart,
			RollingEnd:   row.RollingEnd,
		})
	}
	return nil
}

// replaceCycles makes the college's offered cycles exactly the given set (upserting each, deleting the rest).
// Must run in the same transaction as the college write.
func replaceCycles(ctx context.Context, db dbinterface.QueryInterface, collegeID int64, cycles []models.CollegeApplicationCycle) error {
	deleteQuery, err := schema.ReadSQLBaseScript("delete_college_application_cycles_not_in.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return err
	}
	upsertQuery, err := schema.ReadSQLBaseScript("upsert_college_application_cycle.sql", SqlGlobalCollegeFiles)
	if err != nil {
		return err
	}

	// must be non-nil: a NULL array would match nothing and delete nothing
	keep := make([]string, 0, len(cycles))
	for _, cycle := range cycles {
		keep = append(keep, cycle.Cycle)
	}

	if _, err := db.Exec(ctx, deleteQuery, collegeID, keep); err != nil {
		return err
	}

	for _, cycle := range cycles {
		if _, err := db.Exec(ctx, upsertQuery, collegeID, cycle.Cycle, cycle.Deadline, cycle.RollingStart, cycle.RollingEnd); err != nil {
			return err
		}
	}
	return nil
}
