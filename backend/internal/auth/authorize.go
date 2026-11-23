package auth

import (
	"context"

	"inspirate-consulting/internal/data/postgres/schema"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetRole(ctx context.Context) string {
	role, _ := ctx.Value("Role").(string)
	return role
}

type Authorizer struct {
	db                     *pgxpool.Pool
	counselorAssignedQuery string
}

func NewAuthorizer(db *pgxpool.Pool) *Authorizer {
	query, err := schema.ReadSQLBaseScript("counselor_assigned_to_student.sql", sqlAuthFiles)
	if err != nil {
		panic("failed to read counselor_assigned_to_student.sql: " + err.Error())
	}
	return &Authorizer{db: db, counselorAssignedQuery: query}
}

func ParseStudentID(raw string) (uuid.UUID, error) {
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, errs.BadRequest("student_id must be a UUID")
	}
	return id, nil
}

func ParseOptionalStudentID(raw string) (*uuid.UUID, error) {
	if raw == "" {
		return nil, nil
	}
	id, err := ParseStudentID(raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func (a *Authorizer) ResolveStudentID(ctx context.Context, requested *uuid.UUID) (uuid.UUID, error) {
	return resolveStudentID(ctx, requested, a.CanAccessStudent)
}

func (a *Authorizer) ResolveViewableStudentID(ctx context.Context, requested *uuid.UUID) (uuid.UUID, error) {
	return resolveStudentID(ctx, requested, a.CanViewStudent)
}

func resolveStudentID(ctx context.Context, requested *uuid.UUID, check func(context.Context, uuid.UUID) error) (uuid.UUID, error) {
	if requested == nil {
		id, err := GetStudentID(ctx)
		if err != nil {
			return uuid.Nil, errs.BadRequest("student_id is required")
		}
		return id, nil
	}
	if err := check(ctx, *requested); err != nil {
		return uuid.Nil, err
	}
	return *requested, nil
}

func (a *Authorizer) CanViewStudent(ctx context.Context, targetStudentID uuid.UUID) error {
	if a == nil {
		return nil
	}

	switch models.Role(GetRole(ctx)) {
	case models.StudentRole:
		return checkOwnStudent(ctx, targetStudentID)

	case models.CounselorRole:
		if _, err := GetCounselorID(ctx); err != nil {
			return errs.Unauthorized()
		}
		return nil

	default:
		return errs.Forbidden("insufficient role")
	}
}

func checkOwnStudent(ctx context.Context, targetStudentID uuid.UUID) error {
	callerStudentID, err := GetStudentID(ctx)
	if err != nil {
		return errs.Unauthorized()
	}
	if callerStudentID != targetStudentID {
		return errs.Forbidden("you can only access your own data")
	}
	return nil
}

func (a *Authorizer) CanAccessStudent(ctx context.Context, targetStudentID uuid.UUID) error {
	if a == nil {
		return nil
	}
	role := GetRole(ctx)

	switch models.Role(role) {
	case models.StudentRole:
		return checkOwnStudent(ctx, targetStudentID)

	case models.CounselorRole:
		counselorID, err := GetCounselorID(ctx)
		if err != nil {
			return errs.Unauthorized()
		}
		var assigned bool
		err = a.db.QueryRow(ctx, a.counselorAssignedQuery, targetStudentID, counselorID).Scan(&assigned)
		if err != nil {
			return errs.InternalServerError("authorization check failed")
		}
		if !assigned {
			return errs.Forbidden("student is not assigned to you")
		}
		return nil

	default:
		return errs.Forbidden("insufficient role")
	}
}
