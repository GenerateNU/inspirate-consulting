package user

import "github.com/sethvargo/go-password/password"

// generateTempPassword creates the temporary password a new user signs up with.
// Need to replace current generation with a proper method that actually follows the constraints we have
func generateTempPassword() (string, error) {
	return password.Generate(16, 4, 4, false, false)

}
