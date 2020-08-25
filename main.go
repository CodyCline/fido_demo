package main

import (
	"fido_demo/controllers"
	"fido_demo/models"
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
	auth.HandleFunc("/login/finish/{username}/{session}", FinishLogin).Methods("POST")
	router.Handle(
		"/api/credentials",
		controllers.EnforceJWTAuth(GetCredentialsFor),
	).Methods("GET")
	router.Handle(
		"/api/profile",
		controllers.EnforceJWTAuth(GetUserProfile),
	).Methods("GET")

	router.Handle(
		"/api/credential/start",
		controllers.EnforceJWTAuth(BeginNewCredential),
	).Methods("GET")

	router.Handle(
		"/api/credential/finish/{nickname}/{session}",
		controllers.EnforceJWTAuth(FinishNewCredential),
	).Methods("POST")

	router.Handle(
		"/api/credential/{id}",
		controllers.EnforceJWTAuth(DeleteUserCredential),
	).Methods("DELETE")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./")))
	serverAddress := ":8080"
	log.Println("starting server at", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, handlers.CORS(header, methods, origins)(router)))
}

type CredentialsResponse struct {
	Success     bool                 `json:"success"`
	Credentials []*models.Credential `json:"credentials"`
}

//GetCredentialsFor grabs all the credentials associated with the user for the frontend
var GetCredentialsFor = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	account := r.Context().Value("account").(string)
	data := models.GetCredentials(account)
	resp := CredentialsResponse{
		Success:     true,
		Credentials: data,
	}
	controllers.JSONResponse(w, resp, http.StatusOK)
	return
})

var GetUserProfile = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("account").(string)
	account := models.GetUser(user)
	if account == nil {
		res := Response{
			Success: false,
			Message: "Error user not found, something went wrong",
		}
		controllers.JSONResponse(w, res, http.StatusInternalServerError)
		return
	}
	controllers.JSONResponse(w, account, http.StatusOK)
	return
})
