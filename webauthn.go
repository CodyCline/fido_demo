package main

import (
	"encoding/json"
	"fido_demo/models"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/gorilla/mux"
	"net/http"
)

//BeginRegistration checks if username is taken then creates an account
//and sends registration data back to the client
func BeginRegistration(w http.ResponseWriter, r *http.Request) {
	//Decode the request
	decoder := json.NewDecoder(r.Body)
	var a models.Account
	err := decoder.Decode(&a)
	if err != nil {
		jsonResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
		return
	}

	//Check if account exists
	account, err := models.GetUser(a.Username)
	if err != nil {
		account = models.NewUser(a.Username, "Test Name")
	}

	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = account.CredentialExcludeList()
	}

	//Generate PublicKeyCredentialCreationOptions, session data
	options, sessionData, err := webAuthn.BeginRegistration(
		account,
		registerOptions,
	)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = sessionStore.SaveWebauthnSession("registration", sessionData, r, w)
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, options, http.StatusOK)

}

func FinishRegistration(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	username := vars["username"]

	// get user
	account, err := models.GetUser(username)
	if err != nil {
		jsonResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
	}

	//User doesn't exist
	if err != nil {
		fmt.Println(err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// load the session data
	sessionData, err := sessionStore.GetWebauthnSession("registration", r)
	if err != nil {
		fmt.Println("Session \n", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	credential, err := webAuthn.FinishRegistration(account, sessionData, r)
	if err != nil {
		fmt.Println("Finish Register \n", err)
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	account.AddCredential(*credential)

	jsonResponse(w, "Registration Success", http.StatusOK)
}

func jsonResponse(w http.ResponseWriter, d interface{}, c int) {
	dj, err := json.Marshal(d)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	fmt.Fprintf(w, "%s", dj)
}
