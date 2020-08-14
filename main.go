package main

import (
	"fido_demo/controllers"
	"fmt"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var webAuthn *webauthn.WebAuthn

func main() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Fido Demo",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:3000",
	})

	if err != nil {
		log.Fatal("failed to create WebAuthn from config:", err)
	}

	router := mux.NewRouter()
	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register/start", StartRegistration).Methods("POST")
	auth.HandleFunc("/register/finish/{username}/{session}", FinishRegistration).Methods("POST")
	auth.HandleFunc("/login/start", StartLogin).Methods("POST")
	auth.HandleFunc("/login/finish/{username}", FinishLogin).Methods("POST")
	router.HandleFunc("/todos", FakeData).Methods("GET")
	router.Use(controllers.EnforceJWTAuth)
	//Todo replace with SPA frontend
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	serverAddress := ":8080"
	log.Println("starting server at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, handlers.CORS(header, methods, origins)(router)))
}

type FakeResponse struct {
	Message string `json:"message"`
}

func FakeData(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("account")
	fmt.Println("the user", user)
	resp := FakeResponse{
		Message: "Hello",
	}
	controllers.JSONResponse(w, resp, http.StatusOK)
}
