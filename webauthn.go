package main

import (
	"encoding/json"
	"fido_demo/models"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/gorilla/mux"
	"net/http"
)

//StartRegistration checks if username is taken then creates an account
//and sends registration data back to the client
func StartRegistration(w http.ResponseWriter, r *http.Request) {
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

//FinishRegistration ...
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

//StartLogin gets user by username, checks if it exists and sends data to the client
func StartLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a models.Account

	err := decoder.Decode(&a)
	if err != nil {
		jsonResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
	}

	user, err := models.GetUser(a.Username)

	// user doesn't exist
	if err != nil {
		jsonResponse(w, "Error: cannot find username", http.StatusBadRequest)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := webAuthn.BeginLogin(user)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store session data as marshaled JSON
	err = sessionStore.SaveWebauthnSession("authentication", sessionData, r, w)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, options, http.StatusOK)
}

//FinishLogin Get user sign off token, increment counter update last used, issue jwt
func FinishLogin(w http.ResponseWriter, r *http.Request) {

	// get username
	vars := mux.Vars(r)
	username := vars["username"]

	// get user
	user, err := models.GetUser(username)

	// user doesn't exist
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// load the session data
	sessionData, err := sessionStore.GetWebauthnSession("authentication", r)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	cred, err := webAuthn.FinishLogin(user, sessionData, r)
	if err != nil {
		jsonResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonResponse(w, cred, http.StatusOK)
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
