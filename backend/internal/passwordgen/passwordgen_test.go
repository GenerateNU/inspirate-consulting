package passwordgen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	hasUpper   = regexp.MustCompile(`[A-Z]`)
	hasLower   = regexp.MustCompile(`[a-z]`)
	hasDigit   = regexp.MustCompile(`[0-9]`)
	hasSpecial = regexp.MustCompile(`[!@#~$%^&*()+|_.,;<>?/{}\-]`)
)

func TestGenerateTempPassword(t *testing.T) {
	t.Parallel()

	allowed := lowerChars + upperChars + digitChars + specialChars
	seen := map[string]bool{}

	// many runs, since a generator that only sometimes misses a class would pass a single check
	for range 1000 {
		password, err := GenerateTempPassword()
		require.NoError(t, err)

		assert.Len(t, password, tempPasswordLength)
		assert.GreaterOrEqual(t, len(password), 8, "signup requires at least 8 characters")
		assert.True(t, hasUpper.MatchString(password), "missing uppercase: %q", password)
		assert.True(t, hasLower.MatchString(password), "missing lowercase: %q", password)
		assert.True(t, hasDigit.MatchString(password), "missing digit: %q", password)
		assert.True(t, hasSpecial.MatchString(password), "missing special character: %q", password)
		for _, c := range password {
			assert.True(t, strings.ContainsRune(allowed, c), "unexpected character %q in %q", c, password)
		}

		assert.False(t, seen[password], "duplicate password %q", password)
		seen[password] = true
	}
}
