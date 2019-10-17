package controller

import (
	"api-shortener/cache"
	"api-shortener/utils"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	googleOauthConfig *oauth2.Config
	// TODO: randomize it
	oauthStateString = "alhamsya"
)

const ex_time_jwt = 5

var access_token string

type User struct {
	Id      string `json:"id"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Claims struct {
	Id          string `json:"id"`
	Email       string `json:"email"`
	AccessToken string `json:"access_token"`
	jwt.StandardClaims
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleOauthConfig = &oauth2.Config{
		RedirectURL: fmt.Sprintf("http://%s:%s/oauth/callback",
			os.Getenv("URL_BACKEND"),
			os.Getenv("PORT_BACKEND"),
		),
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	q := c.Request.URL.Query()
	content, err := getUserInfo(q.Get("state"), q.Get("code"))
	if err != nil {
		fmt.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	user := User{}

	json_err := json.Unmarshal(content, &user)
	if json_err != nil {
		log.Fatal(json_err)
	}

	jwt_token, err_jwt := getJwtToken(user)
	if err_jwt != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrMsg{
			Status:  false,
			Message: err_jwt.Error(),
		})
	}

	cache.SetValueWithTTL("AUTH", user.Email, jwt_token, ex_time_jwt*60)

	c.Redirect(http.StatusFound, "/?token="+jwt_token)
	c.Abort()
}

func getUserInfo(state string, code string) ([]byte, error) {
	if state != oauthStateString {
		return nil, fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	access_token = token.AccessToken

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, nil
}

func getJwtToken(user User) (string, error) {
	expirationTime := time.Now().Add(ex_time_jwt * time.Hour)
	claims := &Claims{
		Id:          user.Id,
		Email:       user.Email,
		AccessToken: access_token,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt_key := []byte(os.Getenv("JWT_KEY"))
	tokenString, err := token.SignedString(jwt_key)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func AuthMiddleware(c *gin.Context) *Claims {
	tokenString := c.Request.Header.Get("Authorization")
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, e error) {
		if jwt.SigningMethodHS256 != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		jwt_key := []byte(os.Getenv("JWT_KEY"))
		return jwt_key, nil
	})

	sess_jwt, _ := cache.GetValue("AUTH", claims.Email)
	if sess_jwt != tokenString {
		return nil
	}

	// if token.Valid && err == nil {
	if token != nil && err == nil {
		return claims
	} else {
		return nil
	}
}

func GoogleLogout(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, utils.ErrMsg{
			Status:  false,
			Message: "Not authorization and access",
		})
		return
	}

	user, err := utils.GetSession(tokenString)
	if user == nil {
		c.JSON(http.StatusUnauthorized, utils.ErrMsg{
			Status:  false,
			Message: err.Error(),
		})
		return
	}

	cache_jwt, _ := cache.GetValue("AUTH", user.Email)
	if CheckLogin(fmt.Sprintf("%v", cache_jwt)) == false {
		c.JSON(http.StatusUnauthorized, utils.ErrMsg{
			Status:  false,
			Message: "Token not valid",
		})
		return
	}

	cache.DelKey("AUTH", user.Email)
	response, err := http.Get(os.Getenv("GOOGLE_LOGOUT") + user.AccessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrMsg{
			Status:  false,
			Message: err.Error(),
		})
		return
	}
	defer response.Body.Close()
	c.Redirect(http.StatusFound, "/")
	c.Abort()
}

func CheckLogin(token string) bool {
	if token == "" {
		return false
	}

	user, _ := utils.GetSession(token)
	if user == nil {
		return false
	}

	sess_jwt, _ := cache.GetValue("AUTH", user.Email)
	if sess_jwt != token {
		return false
	}

	return true
}
