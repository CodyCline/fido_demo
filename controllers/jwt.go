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
	ID       uint
	Username string `json:"username"`
	Name     string `json:"name"`
	jwt.StandardClaims
}

type JWTResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

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

//EnforceJWTAuth is custom middleware to enforce authentication on all routes
//except the ones in the exclusion list
func EnforceJWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//List of endpoints that doesn't require auth
		authRequired := []string{
			"/todos",
		}
		requestPath := r.URL.Path //current request path
		fmt.Println("Req\n", requestPath)

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range authRequired {

			if value != requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response := JWTResponse{
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
			response := JWTResponse{
				StatusCode: 403,
				Message:    "Invalid authentication token",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token

		fmt.Println(tokenPart)
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenPart, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})

		//Error decoding the token
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			fmt.Println(err)
			response := JWTResponse{
				StatusCode: 403,
				Message:    "Malformed authentication token",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		//Token is invalid, maybe not signed on this server
		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response := JWTResponse{
				StatusCode: 403,
				Message:    "Token is not valid",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		fmt.Println(claims)

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "account", claims.ID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

//JSONResponse returns a serialized response
var JSONResponse = func(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", dj)
}
