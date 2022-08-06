package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io"
	"log"
	"net/http"
	"strings"
)

func jsonError(w io.Writer, msg string) {
	resp, errMrsh := json.Marshal(map[string]interface{}{
		"message": msg,
	})
	if errMrsh != nil {
		log.Println("Error of Marshal", errMrsh)
	}
	_, err := w.Write(resp)
	if err != nil {
		log.Println("Error of write", err)
	}
}

var (
	ExampleTokenSecret = []byte("супер секретный ключ")
)

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inToken := r.Header.Get("authorization")
		if inToken == "" {
			next.ServeHTTP(w, r)
			return
		}
		_, err := GetInfoFromToken(inToken)
		if err == nil {
			next.ServeHTTP(w, r)
			return
		}
		fmt.Println(err)
		fmt.Println("no auth")
		w.WriteHeader(422)
		jsonError(w, "no auth")
	})
}
func GetInfoFromToken(token string) (map[string]interface{}, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return ExampleTokenSecret, nil
	}
	token = strings.Split(token, " ")[1]
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, hashSecretGetter)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	user := claims["user"].(map[string]interface{})
	data := map[string]interface{}{
		"Login": user["username"].(string),
		"ID":    user["id"],
		"iat":   claims["iat"],
		"exp":   claims["exp"],
	}
	return data, err
}
