package globalCollegeRepository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"inspirate-consulting/internal/errs"
)

// Postgres SQLSTATE codes relevant to global college writes.
const (
	pgUniqueViolationCode     = "23505"
	pgForeignKeyViolationCode = "23503"
	pgCheckViolationCode      = "23514"
)

// mapCycleError turns constraint violations from college_application_cycles writes into 400s.
// The handler validates cycles first, so these are a backstop.
func mapCycleError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}
	switch pgErr.Code {
	case pgCheckViolationCode:
		return errs.BadRequest("invalid application cycle: " + pgErr.ConstraintName)
	case pgUniqueViolationCode:
		return errs.BadRequest("each application cycle can only be listed once")
	}
	return err
}
