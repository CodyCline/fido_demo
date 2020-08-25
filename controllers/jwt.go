package controllers

import (
	"context"
	"encoding/json"
	"fido_demo/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var jwtKey = []byte("256-bit-key")

type Claims struct {
	jwt.StandardClaims
	ID       uint
	Username string `json:"username"`
	Name     string `json:"name"`
}

//AuthResponse is a middleware response typically indicating there is a improper or lack of token
type AuthResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

//CreateJWT issues a token with username as claims
func CreateJWT(a *models.Account) string {
	t, err := strconv.ParseUint("60", 10, 32)
	expirationTime := time.Now().Add(time.Duration(t) * time.Minute)

	claims := &Claims{
		Username: a.Username,
		Name:     a.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			Issuer:    "localhost",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return err.Error()
	}
	return tokenString
}

//EnforceJWTAuth is custom middleware to enforce authentication on specified
func EnforceJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response := AuthResponse{
				StatusCode: 401,
				Message:    "Missing authentication token",
			}
			JSONResponse(w, response, http.StatusUnauthorized)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //Split token from `Bearer {token}`
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response := AuthResponse{
				StatusCode: 403,
				Message:    "Invalid authentication token",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		//Error decoding the token
		if err != nil {
			response := AuthResponse{
				StatusCode: 403,
				Message:    "Malformed or expired authentication token",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		//Token is invalid, maybe not signed on this server
		if !token.Valid {
			response := AuthResponse{
				StatusCode: 403,
				Message:    "Token is not valid",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "account", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx)) //proceed in the middleware chain!
		return
	})
}

//JSONResponse returns a marshalled response with a status code
func JSONResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.WriteHeader(c)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", dj)
	return
}
