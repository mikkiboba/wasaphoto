package api

import (
	"errors"
	"net/http"
	"encoding/json"
	"database/sql"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

/*
Handler of the operation PUT for the route /users/:username/bans/:banname
*/
func (rt *_router) banUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username from the parameters
	username := ps.ByName("username")

	// Get the username's block id
	bid := ps.ByName("banname")

	// Get the tokens
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the token from the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bid, err = rt.db.GetToken(bid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the token from the banned user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is trying to ban themselves
	if uid == bid {
		rt.baseLogger.Error("An user can't ban themselves")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is already banned
	isBanned, err := rt.db.CheckBan(uid, bid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is already banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if isBanned {
		rt.baseLogger.Error("The user ", bid, " is already banned")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	StartTransaction(rt, w)

	// Ban the user
	err = rt.db.BanUser(uid, bid)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the ban into the database. Maybe at lease one of the users doesn't exists.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is and error with the insert of the ban into the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	// Remove the follow
	isFollowing, err := rt.db.CheckFollow(uid, bid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with checking if the banned user is following the banning user")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}
	if isFollowing {
		err = rt.db.UnfollowUser(uid, bid)
		if err != nil {
			rt.baseLogger.WithError(err).Error("There is an error with the delete of the follow 1 from the database")
			w.WriteHeader(http.StatusInternalServerError)
			Rollback(rt, w)
			return
		}
	}

	// Remove the follow from the other sideeee
	isFollowing, err = rt.db.CheckFollow(bid, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with checking if the banned user is following the banning user")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}
	if isFollowing {
		err = rt.db.UnfollowUser(bid, uid)
		if err != nil {
			rt.baseLogger.WithError(err).Error("There is an error with the delete of the follow 2 from the database")
			w.WriteHeader(http.StatusInternalServerError)
			Rollback(rt, w)
			return
		}
	}

	// Delete all comments on the user's posts
	err = rt.db.DeleteCommentsBan(uid, bid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("There is an error while deleting all the comments of the banned user")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	err = rt.db.DeleteCommentsBan(bid, uid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("There is an error while deleting all the comments of the banning user")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}


	Commit(rt, w)

	// No content (204) -> the operation was successful
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("User banned successfully")

}

/*
Handler of the operation DELETE for the route /users/:username/bans/:banname
*/
func (rt *_router) unbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username from the parameters
	username := ps.ByName("username")

	// Get the username's block id
	bid := ps.ByName("banname")

	// Get the tokens
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the token from the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bid, err = rt.db.GetToken(bid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the token from the banned user")
	}

	// Check if the user is trying to ban themselves
	if uid == bid {
		rt.baseLogger.Error("An user can't unban themselves")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the user is already banned
	isBanned, err := rt.db.CheckBan(uid, bid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is already unbanned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !isBanned {
		rt.baseLogger.Error("The user ", bid, " is not banned")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	StartTransaction(rt, w)

	// Unban the user
	err = rt.db.UnbanUser(uid, bid)
	if errors.Is(err, database.ErrElementNotDeleted) {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the ban from the database. Maybe at lease one of the users doesn't exists.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the ban from the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	// No content (204) -> the operation was successful
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("User unbanned successfully")
}


/*
Handler of the operation GET for the route /users/:username/bans/:banname
*/
func (rt *_router) getBan(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	userbanning := ps.ByName("username")
	userbanned := ps.ByName("banname")

	userbanning, err := rt.db.GetToken(userbanning)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the banning user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userbanned, err = rt.db.GetToken(userbanned)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the banned user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} 

	isBanned, err := rt.db.CheckBan(userbanning, userbanned)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is banned")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
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