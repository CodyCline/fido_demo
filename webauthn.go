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
	"strconv"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type RegisterChallenge struct {
	Success     bool                         `json:"success"`
	Options     *protocol.CredentialCreation `json:"options"`
	SessionData *webauthn.SessionData        `json:"session_data"`
	Username    string                       `json:"username"`
}

type LoginChallenge struct {
	Options     *protocol.CredentialAssertion `json:"options"`
	SessionData *webauthn.SessionData         `json:"session_data"`
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
	account := models.GetUser(a.Username)
	if account != nil {
		res := Response{
			Success: false,
			Message: "Username already taken",
		}
		controllers.JSONResponse(w, res, http.StatusOK)
		return
	}
	//Create account issue challenge
	account = models.NewUser(a.Username, a.Name)

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

	//Todo send the session id back down to client or jwt

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
	raw, e := base64.StdEncoding.DecodeString(session)
	if e != nil {
		fmt.Println(e)
	}
	var sess = webauthn.SessionData{}
	err := json.Unmarshal(raw, &sess)

	// get user
	account := models.GetUser(username)
	if account == nil {
		controllers.JSONResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
		return
	}

	credential, err := webAuthn.FinishRegistration(account, sess, r)
	if err != nil {
		fmt.Println("Finish Register \n", err)
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	account.AddCredential(*credential, "Default Credential")

	response := Response{
		Message: "Registration Successful",
		Success: true,
	}

	controllers.JSONResponse(w, response, http.StatusOK)
}

//StartLogin gets user by username, checks if it exists and sends data to the client
func StartLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var a models.Account
	err := decoder.Decode(&a)
	if err != nil {
		controllers.JSONResponse(w, fmt.Errorf("Error: please supply a username"), http.StatusBadRequest)
		return
	}

	account := models.GetUser(a.Username)

	// user doesn't exist
	if account == nil {
		res := Response{
			Success: false,
			Message: "Cannot find username",
		}
		controllers.JSONResponse(w, res, http.StatusNotFound)
		return
	}

	// generate PublicKeyCredentialRequestOptions, session data
	options, sessionData, err := webAuthn.BeginLogin(account)
	if err != nil {
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//TODO Find a way to make not expose this private, probably jwt without using session
	resp := LoginChallenge{
		Options:     options,
		SessionData: sessionData,
	}
	controllers.JSONResponse(w, resp, http.StatusOK)
	return
}

//FinishLogin Get user sign off token, increment counter update last used, issue jwt
func FinishLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	session := vars["session"]
	raw, decodeErr := base64.StdEncoding.DecodeString(session)
	if decodeErr != nil {
		fmt.Println("Error decode", decodeErr)
	}
	var sess = webauthn.SessionData{}
	err := json.Unmarshal(raw, &sess)
	if err != nil {
		fmt.Println("Err", err)
	}

	// get user
	account := models.GetUser(username)

	// user doesn't exist
	if account == nil {
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Todo increment counter update last used property
	credential, err := webAuthn.FinishLogin(account, sess, r)
	if err != nil {
		fmt.Println("Err", err)
		return
	}

	models.UpdateCredential(credential.Authenticator.AAGUID, credential.Authenticator.SignCount)

	response := Response{
		Message: "Login Successful",
		Token:   controllers.CreateJWT(account),
		Success: true,
	}
	controllers.JSONResponse(w, response, http.StatusOK)
	return
}

//Add additional authenticators to your account.

//BeginNewCredential starts the process of creating an account however it is jwt protected
var BeginNewCredential = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("account").(string)
	account := models.GetUser(user)
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

	//Todo send the session id back down to client or jwt

	resp := RegisterChallenge{
		Options:     options,
		SessionData: sessionData,
	}

	controllers.JSONResponse(w, resp, http.StatusOK)
	return
})

var FinishNewCredential = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	nickname := vars["nickname"]
	session := vars["session"]
	raw, e := base64.StdEncoding.DecodeString(session)
	if e != nil {
		res := Response{
			Success: false,
			Message: "Internal server error",
		}
		controllers.JSONResponse(w, res, http.StatusInternalServerError)
		return
	}
	var sess = webauthn.SessionData{}
	err := json.Unmarshal(raw, &sess)

	// get user
	user := r.Context().Value("account").(string)
	account := models.GetUser(user)

	credential, err := webAuthn.FinishRegistration(account, sess, r)
	if err != nil {
		fmt.Println("Finish Register \n", err)
		controllers.JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newCred := account.AddCredential(*credential, nickname)

	controllers.JSONResponse(w, newCred, http.StatusOK)
	return
})

var DeleteUserCredential = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// user := r.Context().Value("account").(string)
	// account := models.GetUser(user)
	vars := mux.Vars(r)
	i := vars["id"]
	id, err := strconv.ParseUint(i, 10, 32)
	if err != nil {
		controllers.JSONResponse(w, "ERROR", http.StatusInternalServerError)
		return
	}
	//TODO CONVERT TO CORRECT TYPE
	models.DeleteCredential(id)
	rs := Response{
		Success: true,
		Message: "Successfully deleted authenticator",
	}
	controllers.JSONResponse(w, rs, http.StatusNoContent)
	return
})
