package todoitem

import (
	"fmt"

	"github.com/google/uuid"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// validateLink returns a 400 if a link id is set but link_type doesn't match it.
func validateLink(linkType *models.TodoLinkType, essayID, mediaID, collegeID *uuid.UUID, extracurricularID *int64) error {
	expect := func(want models.TodoLinkType) error {
		if linkType == nil || *linkType != want {
			return errs.BadRequest(fmt.Sprintf("link_type must be %q when its id is set", want))
		}
		return nil
	}

	switch {
	case essayID != nil:
		return expect(models.TodoLinkEssay)
	case mediaID != nil:
		return expect(models.TodoLinkModule)
	case collegeID != nil:
		return expect(models.TodoLinkCollege)
	case extracurricularID != nil:
		return expect(models.TodoLinkExtracurricular)
	}
	return nil
}