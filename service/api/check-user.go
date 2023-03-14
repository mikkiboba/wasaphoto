package api

import (
	"net/http"

	"git.sapienzaapps.it/fantasticcoffee/fantastic-coffee-decaffeinated/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

/*
Check if the action of a user is done by the user themselves.
Use the action's function as a parameter for this one.
*/
func (rt *_router) checkUser(routerHandler httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

		// Get the username from the parameters
		username := ps.ByName("username")

		// Obtain the token from the username
		uid, err := rt.db.GetToken(username)
		if err != nil {
			rt.baseLogger.WithError(err).Error("There is an error with getting the authorization token")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check if the token is the same of the one in the context -> that means we're checking if the user is acting on their account
		if ctx.Token != uid {
			rt.baseLogger.Error("The user is trying to act on somebody else's profile")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Calling the next routerHandler
		routerHandler(w, r, ps, ctx)

	}
}
