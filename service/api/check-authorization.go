package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
Check if the user is using the right token to authorize their actions.
Use the action's function as a parameter for this one.
*/
func (rt *_router) checkAuthorization(routerHandler httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {

	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		// Get the token from "Authorization" in the header
		token := r.Header.Get("Authorization")

		// Checks if the token is not empty
		if token == "" {
			rt.baseLogger.Error("The authorization token is empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Split the header to get only the token
		token = strings.Split(token, "Bearer ")[1]
		valid, err := rt.db.CheckToken(token)
		// Check if there's any other error execpt ErrNoRows (<- that one is handled differently)
		if err != nil && !(errors.Is(err, sql.ErrNoRows)) {
			rt.baseLogger.WithError(err).Error("There is an error with the database while checking if the token is valid")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check if the token doesn't exists (this is where ErrNoRows is handled)
		if !valid {
			rt.baseLogger.Error("There is an error with the token: token invalid")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Saving the token into the context
		ctx.Token = token

		// Calling the next routerHandler
		routerHandler(w, r, ps, ctx)
	}

}
