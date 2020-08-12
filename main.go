package main

import (
	"github.com/duo-labs/webauthn.io/session"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

var webAuthn *webauthn.WebAuthn
var sessionStore *session.Store

func main() {
	var err error
	webAuthn, err = webauthn.New(&webauthn.Config{
		RPDisplayName: "Fido Demo",
		RPID:          "localhost",
		RPOrigin:      "http://localhost:8080",
	})

	if err != nil {
		log.Fatal("failed to create WebAuthn from config:", err)
	}

	sessionStore, err = session.NewStore()
	if err != nil {
		log.Fatal("failed to create session store:", err)
	}

	r := mux.NewRouter()
	auth := r.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register/start/", StartRegistration).Methods("POST")
	auth.HandleFunc("/register/finish/{username}", FinishRegistration).Methods("POST")
	auth.HandleFunc("/login/start/", StartLogin).Methods("POST")
	auth.HandleFunc("/login/finish/{username}", FinishLogin).Methods("POST")
	//Todo replace with SPA frontend
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	serverAddress := ":8080"
	log.Println("starting server at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))
}
