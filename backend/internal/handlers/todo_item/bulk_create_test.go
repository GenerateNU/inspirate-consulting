package todoitem

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGraduationYearsToStudentYears(t *testing.T) {
	t.Parallel()

	// November 2026 is in the 2026–27 school year, so seniors graduate in 2027
	now := time.Date(2026, time.November, 21, 0, 0, 0, 0, time.UTC)
	got := graduationYearsToStudentYears([]int{2027, 2029, 2035}, now)
	assert.Equal(t, []string{"senior", "sophomore"}, got)
}
