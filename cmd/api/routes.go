package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//character endpoints
	router.HandlerFunc(http.MethodPost, "/v1/characters", app.createCharacterHandler)
	router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.showCharacterHandler)
	router.HandlerFunc(http.MethodGet, "/v1/characters", app.listCharactersHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/characters/:id", app.updateCharacterHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/characters/:id", app.deleteCharacterHandler)

	return app.recoverPanic(app.logRequest(router))
}
