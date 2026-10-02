package globalCollegeRepository

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/jackc/pgx/v5"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

func (r *GlobalCollegeRepository) GetGlobalCollege(ctx context.Context, id int64) (*models.GlobalCollege, error) {

	selectQuery, err := schema.ReadSQLBaseScript("get_global_college.sql", SqlGlobalCollegeFiles)
	if err != nil {
		fmt.Println("[db] ReadSQLBaseScript error:", err)
		err := errs.InternalServerError("Failed to read base query: ", err.Error())
		return nil, err
	}

	rows, err := r.db.Query(ctx, selectQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	globalCollege, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[models.GlobalCollege])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NotFound("global college", "id", strconv.FormatInt(id, 10))
		}
		return nil, err
	}

	return &globalCollege, nil
}
