package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
Handler of the operation PUT for the path /users/:username.

It needs the current username as a parameter and the new username as a json in the request body.
*/
func (rt *_router) setMyUsername(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the - current - username from the parameters
	username := ps.ByName("username")
	username = strings.ReplaceAll(username, " ", "")
	
	// Get the - new - username from the request body as a json
	var un Username
	err := json.NewDecoder(r.Body).Decode(&un)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the json in the request body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	un.Username = strings.ReplaceAll(un.Username, " ", "")
	if !un.isValid() {
		rt.baseLogger.WithError(err).Error("The username is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Start the transaction
	StartTransaction(rt, w)

	// Check if the username already exists in the database
	isThere, err := rt.db.CheckUsername(un.Username)
	if isThere {
		// If the username already exists -> error and rollback of the transaction
		rt.baseLogger.WithError(err).Error("The username is already taken")
		w.WriteHeader(http.StatusConflict)
		Rollback(rt, w)
		return
	}

	// Get the token of the user
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the token")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Set the new username
	err = rt.db.SetUsername(uid, un.Username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the change of the username")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Commit of the transaction
	Commit(rt, w)

	// Action successfully done (204 NO CONTENT)
	w.WriteHeader(http.StatusNoContent)
	rt.baseLogger.Info("Username set successfully")
}
