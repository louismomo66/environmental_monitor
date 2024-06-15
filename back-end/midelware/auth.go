package midelware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/louismomo66/logger/utils"
	"golang.org/x/crypto/bcrypt"
)
var (
	Secretkey string = "secretkeyjwt"
)
func IsAuthorized(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            Err := utils.SetError(nil, "No Token Found")
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Err)
            return
        }

        tokenString := strings.Split(authHeader, "Bearer ")[1]
        fmt.Println("Received Token:", tokenString) // Debugging

        mySigningKey := []byte(Secretkey)
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return mySigningKey, nil
        })

        if err != nil {
            Err := utils.SetError(err, "Your Token has expired")
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(Err)
            return
        }

        if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
            role := claims["role"].(string)
            r.Header.Set("Role", role)
            handler.ServeHTTP(w, r)
            return
        }

        Err := utils.SetError(nil, "Not Authorized")
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(Err)
    }
}





func GenerateJWT(email, role string) (string, error) {
    mySigningKey := []byte(Secretkey)
    token := jwt.New(jwt.SigningMethodHS256)
    claims := token.Claims.(jwt.MapClaims)
    claims["authorized"] = true
    claims["email"] = email
    claims["role"] = role
    claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

    tokenString, err := token.SignedString(mySigningKey)
    if err != nil {
        fmt.Printf("Something went wrong: %s", err.Error())
        return "", err
    }
    return tokenString, nil
}

//take password as input and generate new hash password from it
func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
//compare plain password with hash password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("HOME PUBLIC INDEX PAGE"))
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Role") != "admin" {
		w.Write([]byte("Not authorized."))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}