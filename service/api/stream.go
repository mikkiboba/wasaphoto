package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var followingList []string

	username := ps.ByName("username")

	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	followingList, err = rt.db.GetFollowingList(uid)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the following list")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbposts, err := rt.db.GetStream(followingList)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the stream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts []Post
	for index := range dbposts {
		posts = append(posts, FromDatabase(dbposts[index]))
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to encode the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getPosts(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := ps.ByName("username")

	uid, err := rt.db.GetToken(username)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the token")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dbposts, err := rt.db.GetPosts(uid)
	if errors.Is(err, sql.ErrNoRows) {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to get the stream")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var posts []Post
	for index := range dbposts {
		posts = append(posts, FromDatabase(dbposts[index]))
	}

	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while trying to encode the json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


}