package helpers

import (
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateJWT(claims jwt.MapClaims, secretKey string, isAccess bool) (string, error) {

	if isAccess {
		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		accessTokenString, err := accessToken.SignedString([]byte(secretKey))

		return accessTokenString, err
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))

	return refreshTokenString, err
}

func AccessTokenGenerate(email, secretKey string) (string, error) {
	return CreateJWT(jwt.MapClaims{
		"exp":   time.Now().Add(time.Minute * time.Duration(15)).Unix(),
		"email": email,
	}, secretKey, true)
}

func RefreshTokenGenerate(email, secretKey string) (string, error) {
	return CreateJWT(jwt.MapClaims{
		"exp":   time.Now().Add(time.Hour * time.Duration(48)).Unix(),
		"email": email,
	}, secretKey, false)
}

func Slugify(title string) string {

	dateTime := time.Now().Format("2006-01-02-15-04-05")

	slugTitle := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	slugTitle = reg.ReplaceAllString(slugTitle, "-")
	slugTitle = strings.Trim(slugTitle, "-")

	return slugTitle + "-" + dateTime
}
