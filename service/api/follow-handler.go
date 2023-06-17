package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

/*
Handler of the PUT operation for the path /users/:username/follows/:followname
*/
func (rt *_router) followUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username from the parameters
	username := ps.ByName("username")

	// Get the token (id) of the user
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the name of the user to follow from the parameters
	followname := ps.ByName("followname")

	// Get the token (id) of the user
	fid, err := rt.db.GetToken(followname)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user to follow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is trying to follow themselves
	if uid == fid {
		rt.baseLogger.Error("An user can't follow themselves")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is banned by the user they want to follow
	isBanned, err := rt.db.CheckBan(uid, fid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user has banned ", fid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user you're trying to follow is banned")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user has banned the user they want to follow
	isBanned, err = rt.db.CheckBan(fid, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if you're banned by ", fid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("You have been banned by ", fid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is already following the user they want to follow
	isFollowing, err := rt.db.CheckFollow(uid, fid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is already following " + fid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isFollowing {
		rt.baseLogger.Error("The user is already following " + fid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	StartTransaction(rt, w)

	// Follow the user
	err = rt.db.FollowUser(uid, fid)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the follow into the database. Maybe the user id or the follow id is not found.")
		w.WriteHeader(http.StatusBadRequest)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the follow into the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	// No content (204) -> the operation was successful
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("User followed successfully")

}

func (rt *_router) unfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username from the parameters
	username := ps.ByName("username")

	// Get the token (id) of the user
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the name of the user to follow from the parameters
	followname := ps.ByName("followname")

	// Get the token (id) of the user
	fid, err := rt.db.GetToken(followname)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user to unfollow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is trying to follow themselves
	if uid == fid {
		rt.baseLogger.Error("An user can't unfollow themselves")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is not following the user they want to unfollow
	isFollowing, err := rt.db.CheckFollow(uid, fid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is already following " + fid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !isFollowing {
		rt.baseLogger.Error("The user is already unfollowed " + fid)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	StartTransaction(rt, w)

	// Follow the user
	err = rt.db.UnfollowUser(uid, fid)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the follow into the database. Maybe the user id or the follow id is not found.")
		w.WriteHeader(http.StatusBadRequest)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the follow into the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	// No content (204) -> the operation was successful
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("User followed successfully")
}

func (rt *_router) checkFollow(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := ps.ByName("username")
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the name of the user to follow from the parameters
	followname := ps.ByName("followname")

	// Get the token (id) of the user
	fid, err := rt.db.GetToken(followname)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with getting the token of the user to unfollow")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isFollowing, err := rt.db.CheckFollow(uid, fid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is already following " + fid)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isFollowing {
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(Status{State: true})
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error encoding")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(Status{State: false})
	if err != nil {
		rt.baseLogger.WithError(err).Error("Error encoding")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
