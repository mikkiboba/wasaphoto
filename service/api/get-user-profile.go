package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
 */
func (rt *_router) getUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var profile User

	username := ps.ByName("username")

	if len(username) >= 12 || len(username) <= 3 {
		rt.baseLogger.Error("The username is not valid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uid, err := rt.db.GetToken(username)
	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("The user isn't in the database")
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user has blocked the one making the request
	isBanned, err := rt.db.CheckBan(uid, ctx.Token)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user has banned the one making the request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user has banned the one making the request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user making the request has banned the user
	isBanned, err = rt.db.CheckBan(ctx.Token, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user making the request has banned the user")
		w.WriteHeader(http.StatusForbidden)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user making the request has banned the user")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Get the profile
	profile.Username = username

	followNumber, err := rt.db.GetFollowNumber(uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the number of followers")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	profile.FollowCounter = followNumber

	followingNumber, err := rt.db.GetFollowingNumber(uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the number of following")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	profile.FollowingCounter = followingNumber

	postNumber, err := rt.db.GetPostNumber(uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the number of posts")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	profile.PostCounter = postNumber

	// Put the user informations into a json file
	err = json.NewEncoder(w).Encode(&profile)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to encode the profile json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
