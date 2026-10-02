package essayreview

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"

	dbinterface "inspirate-consulting/internal/data/db-interface"
	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/models"
)

// expectTx makes the mock run the callback the handler passes to WithTx, so
// the orchestration inside it is exercised. The mocked queries ignore db.
func expectTx(mockRepo *mocks.EssayReviewRepository) {
	mockRepo.On("WithTx", mock.Anything, mock.Anything).
		Return(func(ctx context.Context, fn func(db dbinterface.QueryInterface) error) error {
			return fn(nil)
		})
}

var fixedTime = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

// transactionRow builds an essay_review_transaction, defaulting to an open spend.
// Mutators override individual fields.
func transactionRow(mutate ...func(*models.EssayReviewTransaction)) *models.EssayReviewTransaction {
	essayID := uuid.New()
	row := &models.EssayReviewTransaction{
		ID:        uuid.New(),
		StudentID: uuid.New(),
		EssayID:   &essayID,
		Subtotal:  -1,
		EntryType: models.ReviewEntrySpend,
		CreatedAt: fixedTime,
		UpdatedAt: fixedTime,
	}

	for _, m := range mutate {
		m(row)
	}

	deriveStatus(row)
	return row
}

func deriveStatus(row *models.EssayReviewTransaction) {
	if row.EntryType != models.ReviewEntrySpend {
		row.Status = nil
		return
	}

	status := models.ReviewStatusOpen
	switch {
	case row.Refund != nil:
		status = models.ReviewStatusRefunded
	case row.CompletedAt != nil:
		status = models.ReviewStatusCompleted
	}
	row.Status = &status
}
