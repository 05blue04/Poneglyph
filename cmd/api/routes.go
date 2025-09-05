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

	//devilfruit endpoints
	router.HandlerFunc(http.MethodPost, "/v1/devilfruits", app.createDevilFruitHandler)
	router.HandlerFunc(http.MethodGet, "/v1/devilfruits/:id", app.showDevilFruitHandler)
	router.HandlerFunc(http.MethodGet, "/v1/devilfruits", app.listDevilFruitsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/devilfruits/:id", app.updateDevilFruitHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/devilfruits/:id", app.deleteDevilFruitHandler)

	//crew endpoints
	router.HandlerFunc(http.MethodPost, "/v1/crews", app.createCrewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/crews/:id", app.showCrewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/crews", app.listCrewsHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/crews/:id", app.updateCrewHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/crews/:id", app.deleteCrewHandler)
	router.HandlerFunc(http.MethodPost, "/v1/crews/:id/members", app.addCrewMemberHandler)
	router.HandlerFunc(http.MethodGet, "/v1/crews/:id/members", app.listCrewMembersHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/crews/:id/members/:character_id", app.deleteCrewMemberHandler)

	return app.recoverPanic(app.rateLimit(app.logRequest(router)))
}
