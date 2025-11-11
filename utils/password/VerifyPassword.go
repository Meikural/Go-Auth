package password

import "golang.org/x/crypto/bcrypt"

// VerifyPassword compares a plain text password with a bcrypt hash
// Returns true if the password matches, false otherwise
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}