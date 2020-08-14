package main

import (
	"encoding/base64"
	"encoding/json"
	"fido_demo/controllers"
	"fido_demo/models"
	"fmt"
	"github.com/duo-labs/webauthn/protocol"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/gorilla/mux"
	"net/http"
)

type LoginSuccess struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RegisterChallenge struct {
	Options     *protocol.CredentialCreation `json:"options"`
	SessionData *webauthn.SessionData        `json:"session_data"`
	Username    string                       `json:"username"`
}

type LoginChallenge struct {
	Options     *protocol.CredentialAssertion `json:"options"`
	SessionData *webauthn.SessionData         `json"session_data"`
	Username    string                        `json:"username"`
}

//StartRegistration checks if username is taken then creates an account
//and sends registration data back to the client
func StartRegistration(w http.ResponseWriter, r *http.Request) {
	//Decode the request
	decoder := json.NewDecoder(r.Body)
	var a models.Account
	err := decoder.Decode(&a)
	if err != nil {
		controllers.JSONResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
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
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Send the session data down to store as cookie in frontend

	resp := RegisterChallenge{
		Options:     options,
		SessionData: sessionData,
	}

	controllers.JSONResponse(w, resp, http.StatusOK)
	return
}

//FinishRegistration ...
func FinishRegistration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	session := vars["session"]
	// decoder := json.NewDecoder(r.Body)
	raw, e := base64.StdEncoding.DecodeString(session)
	if e != nil {
		fmt.Println(e)
	}
	var sess = webauthn.SessionData{}
	err := json.Unmarshal(raw, &sess)

	// get user
	account, notFound := models.GetUser(username)
	if notFound != nil {
		controllers.JSONResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
		return
	}

	credential, err := webAuthn.FinishRegistration(account, sess, r)
	if err != nil {
		fmt.Println("Finish Register \n", err)
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account.AddCredential(*credential)

	controllers.JSONResponse(w, "SUCCESS!", http.StatusOK)
}

//StartLogin gets user by username, checks if it exists and sends data to the client
func StartLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a models.Account

	err := decoder.Decode(&a)
	if err != nil {
		controllers.JSONResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
	}

	user, err := models.GetUser(a.Username)

	// user doesn't exist
	if err != nil {
		controllers.JSONResponse(w, "Error: cannot find username", http.StatusBadRequest)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := webAuthn.BeginLogin(user)
	if err != nil {
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// store session data as marshaled JSON
	resp := LoginChallenge{
		Options:     options,
		SessionData: sessionData,
	}
	controllers.JSONResponse(w, resp, http.StatusOK)
}

//FinishLogin Get user sign off token, increment counter update last used, issue jwt
func FinishLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	session := vars["session"]
	// decoder := json.NewDecoder(r.Body)
	raw, e := base64.StdEncoding.DecodeString(session)
	if e != nil {
		fmt.Println(e)
	}
	var sess = webauthn.SessionData{}
	err := json.Unmarshal(raw, &sess)

	// get user
	user, err := models.GetUser(username)

	// user doesn't exist
	if err != nil {
		controllers.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	//Todo increment counter
	_, err = webAuthn.FinishLogin(user, sess, r)
	if err != nil {
		controllers.JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	response := LoginSuccess{
		Message: "Login Successful",
		Token:   controllers.CreateJWT(user),
		Success: true,
	}
	controllers.JSONResponse(w, response, http.StatusOK)
}
