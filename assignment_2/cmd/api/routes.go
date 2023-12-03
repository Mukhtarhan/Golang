package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/videos", app.requirePermission("videos:read", app.listVideosHandler))
	router.HandlerFunc(http.MethodPost, "/v1/videos", app.requirePermission("videos:write", app.createVideoHandler))
	router.HandlerFunc(http.MethodGet, "/v1/videos/:id", app.requirePermission("videos:read", app.showVideoHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/videos/:id", app.requirePermission("videos:write", app.updateVideoHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/videos/:id", app.requirePermission("videos:write", app.deleteVideoHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
