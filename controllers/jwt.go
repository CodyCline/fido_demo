package controllers

import (
	"context"
	"encoding/json"
	"fido_demo/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"strings"
)

type JWTResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

var JwtAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{
			"/auth/register/start",
			"/auth/register/finish/{username}",
			"/auth/login/start",
			"/auth/login/finish",
			"/",
		} //List of endpoints that doesn't require auth
		requestPath := r.URL.Path //current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {

			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			response := JWTResponse{
				Status:  true,
				Message: "Missing authentication token",
			}
			JSONResponse(w, response, http.StatusForbidden)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			JSONResponse(w, fmt.Errorf("Invalid/Malformed auth token"), http.StatusForbidden)
			return
		}

		tokenPart := splitted[1] //Grab the token part
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			JSONResponse(w, JWTResponse{false, "Malformed authentication token"}, http.StatusForbidden)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			JSONResponse(w, JWTResponse{false, "Invalid authentication token"}, http.StatusBadRequest)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		ctx := context.WithValue(r.Context(), "account", tk.ID)
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
