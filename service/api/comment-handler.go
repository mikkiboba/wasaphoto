package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) comment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the username
	username := ps.ByName("username")

	// Get the id from the username
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the post id
	pid := ps.ByName("postid")

	// Get the owner of the post
	oid, err := rt.db.GetPostOwner(pid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the post owner's id")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the owner of the post banned the user who is commenting
	isBanned, err := rt.db.CheckBan(oid, uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the user is banned by the post's owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The user is banned by the post's owner")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Check if the user who is commenting has banned the post's owner
	isBanned, err = rt.db.CheckBan(uid, oid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while checking if the post's owner is banned by the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if isBanned {
		rt.baseLogger.Error("The post's owner is banned by the user")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the comment's content (the text)
	var text Comment
	err = json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		rt.baseLogger.Error("There is an error while decoding the input json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Comment the photo
	StartTransaction(rt, w)

	err = rt.db.Comment(uid, pid, text.Text)
	if errors.Is(err, database.ErrElementNotAdded) {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the comment into the database. Maybe the user or the post doesn't exist")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the insert of the comment into the database")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	w.WriteHeader(http.StatusNoContent)

	rt.baseLogger.Info("Photo commented succesfully")
}

func (rt *_router) uncomment(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	rt.baseLogger.Info("Uncommenting")

	// Get the username
	username := ps.ByName("username")

	// Get the user's id
	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the user's token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get the post id
	pid := ps.ByName("postid")

	// Get the comment id
	cid := ps.ByName("commentid")

	// Uncomment the photo
	StartTransaction(rt, w)

	err = rt.db.Uncomment(uid, pid, cid)
	if errors.Is(err, database.ErrElementNotDeleted) {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the comment from the database. Maybe the comment doesn't exist")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error with the delete of the comment from the database")
		w.WriteHeader(http.StatusInternalServerError)
		Rollback(rt, w)
		return
	}

	Commit(rt, w)

	rt.baseLogger.Info("Post uncommented succesfully")
	w.WriteHeader(http.StatusNoContent)
}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the post id
	pid := ps.ByName("postid")

	comments, err := rt.db.GetComments(pid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the comments")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var tComments []Comment

	for _, comment := range comments {
		rComment := FromDatabaseComment(comment)
		rComment.User, err = rt.db.GetUsername(rComment.User)
		if err != nil {
			rt.baseLogger.WithError(err).Error("Error while trying to get the username")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		tComments = append(tComments, rComment)

	}

	err = json.NewEncoder(w).Encode(tComments)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to encode the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
