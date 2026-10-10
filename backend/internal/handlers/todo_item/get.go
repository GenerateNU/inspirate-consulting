package todoitem

import (
	"context"
	"strings"

	"inspirate-consulting/internal/auth"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
	"inspirate-consulting/internal/pagination"

	"github.com/google/uuid"
)

func (h *Handler) ListTodoItems(ctx context.Context, input *models.ListTodoItemsInput) (*models.TodoItemPage, error) {
	query, err := buildTodoItemQuery(ctx, input.TodoItemFilters, input.Cursor, input.Limit)
	if err != nil {
		return nil, err
	}
	return h.listTodoItemPage(ctx, query)
}

func (h *Handler) SearchTodoItems(ctx context.Context, input *models.SearchTodoItemsInput) (*models.TodoItemPage, error) {
	search := strings.TrimSpace(input.Q)
	if search == "" {
		return nil, errs.BadRequest("search query must not be blank")
	}

	query, err := buildTodoItemQuery(ctx, input.TodoItemFilters, input.Cursor, input.Limit)
	if err != nil {
		return nil, err
	}
	query.Search = &search
	return h.listTodoItemPage(ctx, query)
}

func (h *Handler) listTodoItemPage(ctx context.Context, query models.TodoItemQuery) (*models.TodoItemPage, error) {
	limit := query.Limit
	query.Limit = limit + 1

	rows, err := h.TodoItemRepository.ListTodoItems(ctx, query)
	if err != nil {
		return nil, err
	}

	items, nextCursor, err := pagination.NextPage(rows, limit, func(item models.TodoItem) models.TodoItemCursor {
		return models.TodoItemCursor{
			SortBy:    query.SortBy,
			SortOrder: query.SortOrder,
			SortValue: query.SortValue(item),
			ID:        item.ID,
		}
	})
	if err != nil {
		return nil, err
	}

	return &models.TodoItemPage{Items: items, NextCursor: nextCursor}, nil
}

func buildTodoItemQuery(ctx context.Context, filters models.TodoItemFilters, cursor string, limit int) (models.TodoItemQuery, error) {
	query := models.TodoItemQuery{
		StudentID: filters.StudentID,
		Status:    filters.Status,
		EssayID:   optionalString(filters.EssayID),
		MediaID:   optionalString(filters.MediaID),
		LinkedTo:  optionalString(filters.LinkedTo),
		SortBy:    filters.SortBy,
		SortOrder: filters.SortOrder,
		Limit:     limit,
	}
	// TODO: restrict student_id so students only see their own todos once auth is implemented.
	if query.StudentID == "" {
		query.StudentID = auth.GetStudentID(ctx)
	}
	if filters.GlobalCollegeID != 0 {
		query.GlobalCollegeID = &filters.GlobalCollegeID
	}
	query.ApplyDefaults()

	if cursor != "" {
		after, err := pagination.Decode[models.TodoItemCursor](cursor)
		if err != nil {
			return query, err
		}
		if _, err := uuid.Parse(after.ID); err != nil {
			return query, errs.BadRequest("invalid cursor")
		}
		// A cursor's sort value is only meaningful under the sort it was taken from.
		if after.SortBy != query.SortBy || after.SortOrder != query.SortOrder {
			return query, errs.BadRequest("cursor does not match sort_by and sort_order")
		}
		query.After = &after
	}

	return query, nil
}

func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
