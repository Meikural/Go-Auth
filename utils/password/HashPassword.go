package password

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain text password and returns a bcrypt hash
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}