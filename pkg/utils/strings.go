package utils

import (
	"math/rand"

	"golang.org/x/crypto/bcrypt"
)

func GenerateOtp() string {
	//cfg := config.GetConfig()

	// min := int(math.Pow(10, float64(cfg.Otp.Digits-1)))   // 10^d-1 100000
	// max := int(math.Pow(10, float64(cfg.Otp.Digits)) - 1) // 999999 = 1000000 - 1 (10^d) -1

	// var num = rand.Intn(max-min) + min
	// return strconv.Itoa(num)
	return ""
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword verifies if the given password matches the stored hash.
func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_@*#$%^&()!<>{}"

	result := ""

	b := make([]byte, length)
	for j := range b {
		b[j] = charset[rand.Intn(len(charset))]
	}
	result = string(b)

	return result
}
