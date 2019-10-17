package utils

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

type SuccessMsg struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Meta    Meta   `json:"meta"`
}

type ErrMsg struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type Meta struct {
	LongUrl  string `json:"long_url"`
	ShortUrl string `json:"short_url"`
}

type Claims struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	jwt.StandardClaims
}

func RandStr(length int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHJKLMNPQRSTWXYZ" +
		"123456789")
	var temp_str strings.Builder

	for a := 0; a < length; a++ {
		temp_str.WriteRune(chars[rand.Intn(len(chars))])
	}
	str := temp_str.String()
	return str
}

func CheckStrUrl(str string) bool {
	pattern := `(?m)^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\(\)\*\+,;=.]+$`
	match, _ := regexp.MatchString(pattern, str)
	return match
}

func GetSession(token_string string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(token_string, claims, func(token *jwt.Token) (i interface{}, e error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		jwt_key := []byte(os.Getenv("JWT_KEY"))
		return jwt_key, nil
	})

	if token != nil && err == nil {
		return claims, err
	} else {
		return nil, err
	}
}

func GetCurrentTime() string {
	currentTime := time.Now()
	result := currentTime.Format("2006-01-02 15:04:05")
	return result
}
