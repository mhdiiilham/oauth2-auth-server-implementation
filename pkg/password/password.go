package password

import "golang.org/x/crypto/bcrypt"

// Hash Password function
func Hash(password string) string {
	hb, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(hb)
}

// Compare Password function
func Compare(plainPassword, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))
	return err == nil
}
