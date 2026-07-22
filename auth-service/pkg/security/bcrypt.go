package security

import "golang.org/x/crypto/bcrypt"

// HashPassword gelen yalın şifreyi bcrypt ile hashler (Cost: 10 ideal seviyedir)
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

// CheckPasswordHash hashlenmiş şifre ile gelen şifreyi kıyaslar
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}