package api

import (
	"database/sql"
	"errors"
	"net/http"
	"encoding/json"


	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

/*
Handler of the operation PUT for the route /users/:username/posts/:postid/likes
*/
func (rt *_router) likePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username of the user liking the photo from the parameters
	username := ps.ByName("username")

	// Get the post id from the parameters
	postID := ps.ByName("postid")

	// Get the id of the user liking the photo
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the owner (id) of the post
	oid, err := rt.db.GetPostOwner(postID)
	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the owner of the post. Maybe the post doesn't exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the owner of the post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is trying to like a post they own
	if uid == oid {
		rt.baseLogger.Error("A user can't like a photo they posted")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the owner banned the user
	isBanned, err := rt.db.CheckBan(oid, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to check if the post's owner banned the user liking")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user is banned by the post's owner")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the user banned the post's owner
	isBanned, err = rt.db.CheckBan(uid, oid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to check if the user liking banned the post's owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The post's owner is banned by the user liking")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	StartTransaction(rt, w)

	// Like the photo
	err = rt.db.LikePhoto(uid, postID)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the like into the databse. Maybe the user or the post don't exist")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the like into the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	// Everything OK
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("Photo liked succesfully")
}

/*
Handler of the operation DELETE for the route /users/:username/posts/:postid/likes
*/
func (rt *_router) unlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username of the user unliking the photo from the parameters
	username := ps.ByName("username")

	// Get the post id from the parameters
	postID := ps.ByName("postid")

	// Get the id of the user unliking the photo
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the owner (id) of the post
	oid, err := rt.db.GetPostOwner(postID)
	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the owner of the post. Maybe the post doesn't exist")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the owner of the post")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user is trying to unlike a post they own
	if uid == oid {
		rt.baseLogger.Error("A user can't unlike a photo they posted")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the owner banned the user
	isBanned, err := rt.db.CheckBan(oid, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to check if the post's owner banned the user unliking")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user is banned by the post's owner")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// Check if the user banned the post's owner
	isBanned, err = rt.db.CheckBan(uid, oid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to check if the user unliking banned the post's owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The post's owner is banned by the user unliking")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	StartTransaction(rt, w)

	// Like the photo
	err = rt.db.UnlikePhoto(uid, postID)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the like into the databse. Maybe the user or the post don't exist")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the like into the database.")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	// Everything OK
	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("Photo unliked succesfully")
}

func (rt *_router) getLike(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username of the user unliking the photo from the parameters
	username := ps.ByName("username")

	// Get the post id from the parameters
	postID := ps.ByName("postid")

	// Get the id of the user liking the photo
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	isLiked, err := rt.db.CheckLike(uid, postID)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking the like")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isLiked {
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
	}
	return



}