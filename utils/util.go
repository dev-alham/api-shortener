package utils

import (
	"api-shortener/cache"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
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

	sess_jwt, _ := cache.GetValue("AUTH", claims.Email)
	if sess_jwt != token_string {
		return nil, err
	}

	if token != nil && err == nil {
		return claims, err
	} else {
		return nil, err
	}
}

func GetCurrentTimeString() string {
	currentTime := time.Now()
	result := currentTime.Format("2006-01-02 15:04:05")
	return result
}

func GetCurrentTime() time.Time {
	currentTime := time.Now()
	return currentTime
}

func ValidateBetween(param int, smallest int, biggest int) bool {
	if (param >= smallest) && (param <= biggest) {
		return true
	} else {
		return false
	}
}

func DeletePrefixUrl(str string) string {
	re := regexp.MustCompile(`(?m)(http|https)://|(www.)`)
	return re.ReplaceAllString(str, "")
}

func GoogleAccountLogout(c *gin.Context, access_token string) bool {
	response, err := http.Get(os.Getenv("GOOGLE_LOGOUT") + access_token)
	if err != nil {
		return false
	}
	defer response.Body.Close()
	return true
}
