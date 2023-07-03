package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	var users []string

	username := r.URL.Query().Get("username")

	users, err := rt.db.GetUsers(username, ctx.Token)

	if errors.Is(err, sql.ErrNoRows) {
		rt.baseLogger.Error("Can't find any users with this prefix: ", username)
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while getting the users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		rt.baseLogger.WithError(err).Error("There is an error while encoding the users")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
