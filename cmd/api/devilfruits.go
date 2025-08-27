package main

import "net/http"

func (app *application) createDevilFruitHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		Type           string   `json:"type"`
		Character_id   int      `json:"character_id"`
		PreviousOwners []string `json:"previous_owners"`
	}

	err := app.readJSON(w, r, input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

}

func (app *application) showDevilFruitHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) updateDevilFruitHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteDevilFruitHandler(w http.ResponseWriter, r *http.Request) {

}
