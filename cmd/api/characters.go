package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/05blue04/Poneglyph/internal/data"
	"github.com/05blue04/Poneglyph/internal/validator"
)

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string        `json:"name"`
		Age         int           `json:"age"`
		Description string        `json:"description"`
		Origin      string        `json:"origin"`
		Fruit       string        `json:"devil_fruit"`
		Bounty      *data.Berries `json:"bounty,omitempty"` //optional field
		Debut       string        `json:"debut"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character := &data.Character{
		Name:        input.Name,
		Age:         input.Age,
		Description: input.Description,
		Origin:      input.Origin,
		Fruit:       input.Fruit,
		Bounty:      input.Bounty,
		Debut:       input.Debut,
	}

	v := validator.New()

	if data.ValidateCharacter(v, character); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Characters.Insert(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/characters/%d", character.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"character": character}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/characters/%d", character.ID))

	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCharacterHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Characters.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "character successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
