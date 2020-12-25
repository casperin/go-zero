package cookie

import (
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type userClaims struct {
	UserId    int
	UserEmail string
	IssuedAt  time.Time
	jwt.StandardClaims
}

var cookieSigningKey = []byte("zero: cookie signing key -- SHOULD BE CHANGED")

func Login(w http.ResponseWriter, userId int, userEmail string, expires time.Time) error {
	claims := userClaims{
		UserId:    userId,
		UserEmail: userEmail,
		IssuedAt:  time.Now(),
	}
	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, e := unsignedToken.SignedString(cookieSigningKey)
	if e != nil {
		return e
	}
	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "jwt",
		Value:   token,
		Expires: expires,
	})
	return nil
}

func GetUser(r *http.Request) (int, string, error) {
	jwtCookie, e := r.Cookie("jwt")
	if e != nil {
		return 0, "", e
	}
	claims := userClaims{}
	token, e := jwt.ParseWithClaims(jwtCookie.Value, &claims, func(token *jwt.Token) (interface{}, error) {
		return cookieSigningKey, nil
	})
	if e != nil {
		return 0, "", e
	}
	if time.Since(claims.IssuedAt) > (time.Hour * 24 * 14) {
		return 0, "", errors.New("token is not valid")
	}
	if token.Valid {
		return claims.UserId, claims.UserEmail, nil
	}
	return 0, "", errors.New("token is not valid")
}

func Logout(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "jwt",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
