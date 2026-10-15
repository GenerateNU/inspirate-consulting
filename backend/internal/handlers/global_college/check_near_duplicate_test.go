package globalcollege

import (
	"context"
	"errors"
	"net/http"
	"testing"

	mocks "inspirate-consulting/internal/data/repo-mocks"
	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandler_CheckNearDuplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	name, location := "Northeastern Boston Campus", "Boston MA"
	selfID := int64(1)

	tests := []struct {
		name        string
		matches     []models.GlobalCollege
		lookupErr   error
		excludeID   *int64
		wantConflict bool
		wantErr     error
	}{
		{name: "no matches: allowed", matches: []models.GlobalCollege{}},
		{
			name:         "a match on create: 409",
			matches:      []models.GlobalCollege{{ID: 7, SchoolName: "Northeastern University - Boston"}},
			wantConflict: true,
		},
		{
			name:      "only match is the college being updated: allowed",
			matches:   []models.GlobalCollege{{ID: selfID, SchoolName: name}},
			excludeID: &selfID,
		},
		{
			name: "self plus another college: 409",
			matches: []models.GlobalCollege{
				{ID: selfID, SchoolName: name},
				{ID: 7, SchoolName: "Northeastern University - Boston"},
			},
			excludeID:    &selfID,
			wantConflict: true,
		},
		{name: "lookup error propagates", lookupErr: errors.New("database error"), wantErr: errors.New("database error")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockRepo := mocks.NewGlobalCollegeRepository(t)
			mockRepo.On("FindNearDuplicateGlobalColleges", mock.Anything, name, location, nearDuplicateThreshold).
				Return(tc.matches, tc.lookupErr)

			err := NewHandler(mockRepo).checkNearDuplicate(ctx, name, location, tc.excludeID)

			switch {
			case tc.wantErr != nil:
				assert.EqualError(t, err, tc.wantErr.Error())
			case tc.wantConflict:
				var httpErr errs.HTTPError
				require.ErrorAs(t, err, &httpErr)
				assert.Equal(t, http.StatusConflict, httpErr.Code)
				// The message should name the existing college the user collided with.
				assert.Contains(t, httpErr.Message, "Northeastern University - Boston")
			default:
				assert.NoError(t, err)
			}
		})
	}
}

func TestFirstOtherCollege(t *testing.T) {
	t.Parallel()

	colleges := []models.GlobalCollege{{ID: 1}, {ID: 2}}
	one, two, three := int64(1), int64(2), int64(3)

	tests := []struct {
		name      string
		colleges  []models.GlobalCollege
		excludeID *int64
		wantID    *int64
	}{
		{name: "no exclusion returns the first", colleges: colleges, wantID: &one},
		{name: "excluded first returns the second", colleges: colleges, excludeID: &one, wantID: &two},
		{name: "exclusion not present returns the first", colleges: colleges, excludeID: &three, wantID: &one},
		{name: "only the excluded college returns nil", colleges: []models.GlobalCollege{{ID: 1}}, excludeID: &one},
		{name: "empty returns nil"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := firstOtherCollege(tc.colleges, tc.excludeID)
			if tc.wantID == nil {
				assert.Nil(t, got)
				return
			}
			if assert.NotNil(t, got) {
				assert.Equal(t, *tc.wantID, got.ID)
			}
		})
	}
}