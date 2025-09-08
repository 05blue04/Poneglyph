package main

import (
	"expvar"
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
	router.HandlerFunc(http.MethodGet, "/v1/characters/:id", app.showCharacterHandler)
	router.HandlerFunc(http.MethodGet, "/v1/characters", app.listCharactersHandler)
	router.Handler(http.MethodPost, "/v1/characters", app.requireAuthOptional(http.HandlerFunc(app.createCharacterHandler)))
	router.Handler(http.MethodPatch, "/v1/characters/:id", app.requireAuthOptional(http.HandlerFunc(app.updateCharacterHandler)))
	router.Handler(http.MethodDelete, "/v1/characters/:id", app.requireAuthOptional(http.HandlerFunc(app.deleteCharacterHandler)))

	//devilfruit endpoints
	router.HandlerFunc(http.MethodGet, "/v1/devilfruits/:id", app.showDevilFruitHandler)
	router.HandlerFunc(http.MethodGet, "/v1/devilfruits", app.listDevilFruitsHandler)
	router.Handler(http.MethodPost, "/v1/devilfruits", app.requireAuthOptional(http.HandlerFunc(app.createDevilFruitHandler)))
	router.Handler(http.MethodPatch, "/v1/devilfruits/:id", app.requireAuthOptional(http.HandlerFunc(app.updateDevilFruitHandler)))
	router.Handler(http.MethodDelete, "/v1/devilfruits/:id", app.requireAuthOptional(http.HandlerFunc(app.deleteDevilFruitHandler)))

	//crew endpoints
	router.HandlerFunc(http.MethodGet, "/v1/crews/:id", app.showCrewHandler)
	router.HandlerFunc(http.MethodGet, "/v1/crews", app.listCrewsHandler)
	router.HandlerFunc(http.MethodGet, "/v1/crews/:id/members", app.listCrewMembersHandler)
	router.Handler(http.MethodPost, "/v1/crews", app.requireAuthOptional(http.HandlerFunc(app.createCrewHandler)))
	router.Handler(http.MethodPatch, "/v1/crews/:id", app.requireAuthOptional(http.HandlerFunc(app.updateCrewHandler)))
	router.Handler(http.MethodDelete, "/v1/crews/:id", app.requireAuthOptional(http.HandlerFunc(app.deleteCrewHandler)))
	router.Handler(http.MethodPost, "/v1/crews/:id/members", app.requireAuthOptional(http.HandlerFunc(app.addCrewMemberHandler)))
	router.Handler(http.MethodDelete, "/v1/crews/:id/members/:character_id", app.requireAuthOptional(http.HandlerFunc(app.deleteCrewMemberHandler)))

	//metric endpoint
	router.Handler(http.MethodGet, "/v1/metrics", app.requireAuthOptional(expvar.Handler()))

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.logRequest(router)))))
}
