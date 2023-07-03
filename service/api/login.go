package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

/*
Handler of the route POST for the path /session.
Get the username from the input json and returns the uuid as a json. The uuid is the authorization token.
If the username is not in the database it will create a new user with that username and generate a new uuid.
*/
func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Getting the username from the input json
	var username Username
	err := json.NewDecoder(r.Body).Decode(&username)

	// Checking errors with the json
	if err != nil {
		rt.baseLogger.WithError(err).Warning("Error with the json file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Checking if the username is valid
	if !username.isValid() {
		rt.baseLogger.Warning("The username is not valid this one")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Login
	var uid string
	uid, err = rt.db.Login(username.Username)
	if errors.Is(err, sql.ErrNoRows) {
		// Registration if the username is not in the database
		// New uuid
		uuid, err := uuid.NewV4()
		uid = uuid.String()
		if err != nil {
			rt.baseLogger.WithError(err).Error("There's an error with the uuid generation")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Registration
		err = rt.db.Registration(uid, username.Username)
		if err != nil {
			rt.baseLogger.WithError(err).Error("There's an error with the registration, can't insert the user into the database")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There's an error with the database")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	rt.baseLogger.Info("User logged in with uuid: ", uid)

	// Returning the uuid as a json
	w.Header().Set("Content-Type", "application/json")
	var t Token
	t.Id = uid
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There's an error with the encoding of the uuid")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
