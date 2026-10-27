package passwordgen

import "github.com/sethvargo/go-password/password"

// Generate creates a temporary password for a new or reset account.
// Need to replace current generation with a proper method that actually follows the constraints we have
func Generate() (string, error) {
	return password.Generate(16, 4, 4, false, false)
}
