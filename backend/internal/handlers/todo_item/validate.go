package todoitem

import (
	"fmt"

	"inspirate-consulting/internal/errs"
	"inspirate-consulting/internal/models"
)

// validateLink returns a 400 if a link id is set but link_type doesn't match it.
func validateLink(linkType *models.TodoLinkType, hasEssay, hasMedia, hasCollege, hasExtracurricular bool) error {
	expect := func(want models.TodoLinkType) error {
		if linkType == nil || *linkType != want {
			return errs.BadRequest(fmt.Sprintf("link_type must be %q when its id is set", want))
		}
		return nil
	}

	switch {
	case hasEssay:
		return expect(models.TodoLinkEssay)
	case hasMedia:
		return expect(models.TodoLinkModule)
	case hasCollege:
		return expect(models.TodoLinkCollege)
	case hasExtracurricular:
		return expect(models.TodoLinkExtracurricular)
	}
	return nil
}
