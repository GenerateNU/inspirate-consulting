package todoitem

import (
	"context"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// Invalid links are rejected before the repository is touched.
func TestHandler_UpdateTodoItemLink_Validation(t *testing.T) {
	t.Parallel()

	essayID := uuid.NewString()
	extracurricularID := int64(12)
	essay := models.TodoLinkEssay
	module := models.TodoLinkModule

	tests := []struct {
		name string
		link models.TodoItemLink
	}{
		{
			name: "link_type doesn't match the id",
			link: models.TodoItemLink{LinkType: &module, EssayID: &essayID},
		},
		{
			name: "two ids at once",
			link: models.TodoItemLink{LinkType: &essay, EssayID: &essayID, ExtracurricularID: &extracurricularID},
		},
		{
			name: "id with no link_type",
			link: models.TodoItemLink{EssayID: &essayID},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// No expectations: the repository must not be called
			mockRepo := mocks.NewTodoItemRepository(t)
			handler := NewHandler(mockRepo, nil)

			_, err := handler.UpdateTodoItemLink(context.Background(), uuid.NewString(), tc.link)
			assert.Error(t, err)
		})
	}
}
