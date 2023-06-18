package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Session routes
	rt.router.POST("/session", rt.wrap(rt.login))

	// User routes
	rt.router.PUT("/users/:username", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.setMyUsername))))
	rt.router.GET("/users/:username", rt.wrap(rt.checkAuthorization(rt.getUserProfile)))

	// Stream routes
	rt.router.GET("/users/:username/stream", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.getMyStream))))

	// Posts routes
	rt.router.POST("/users/:username/posts", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.uploadPhoto))))
	rt.router.DELETE("/users/:username/posts/:postid", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.deletePhoto))))

	// Follow routes
	rt.router.PUT("/users/:username/follows/:followname", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.followUser))))
	rt.router.DELETE("/users/:username/follows/:followname", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.unfollowUser))))

	// Ban routes
	rt.router.PUT("/users/:username/bans/:banname", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.banUser))))
	rt.router.DELETE("/users/:username/bans/:banname", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.unbanUser))))

	// Like routes
	rt.router.PUT("/users/:username/posts/:postid/likes", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.likePhoto))))
	rt.router.DELETE("/users/:username/posts/:postid/likes", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.unlikePhoto))))

	// Comment routes
	rt.router.POST("/users/:username/posts/:postid/comments", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.comment))))
	rt.router.DELETE("/users/:username/posts/:postid/comments/:commentid", rt.wrap(rt.checkAuthorization(rt.checkUser(rt.uncomment))))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// EXTRA
	rt.router.GET("/posts/:postid", rt.wrap(rt.checkAuthorization(rt.getPhoto)))
	rt.router.GET("/users", rt.wrap(rt.checkAuthorization(rt.getUsers)))
	rt.router.GET("/users/:username/bans/:banname", rt.wrap(rt.checkAuthorization(rt.getBan)))
	rt.router.GET("/users/:username/follows/:followname", rt.wrap(rt.checkAuthorization(rt.checkFollow)))
	rt.router.GET("/users/:username/posts", rt.wrap(rt.checkAuthorization(rt.getPosts)))
	rt.router.GET("/users/:username/posts/:postid/likes", rt.wrap(rt.checkAuthorization(rt.getLike)))
	rt.router.GET("/posts/:postid/comments", rt.wrap(rt.getComments))

	return rt.router
}
