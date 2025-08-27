package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/05blue04/Poneglyph/internal/data"
	"github.com/05blue04/Poneglyph/internal/validator"
)

func (app *application) createCharacterHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string        `json:"name"`
		Age         int           `json:"age"`
		Description string        `json:"description"`
		Origin      string        `json:"origin"`
		Bounty      *data.Berries `json:"bounty,omitempty"` //optional field
		Race        string        `json:"race"`
		Episode     int           `json:"episode"`
		TimeSkip    string        `json:"time_skip"`
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
		Bounty:      input.Bounty,
		Race:        strings.ToLower(input.Race),
		Episode:     input.Episode,
		TimeSkip:    strings.ToLower(input.TimeSkip),
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

func (app *application) updateCharacterHandler(w http.ResponseWriter, r *http.Request) {

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

	var input struct {
		Name        *string       `json:"name"`
		Age         *int          `json:"age"`
		Description *string       `json:"description"`
		Origin      *string       `json:"origin"`
		Bounty      *data.Berries `json:"bounty,omitempty"`
		Race        *string       `json:"race"`
		Episode     *int          `json:"episode"`
		TimeSkip    *string       `json:"time_skip"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Race != nil {
		race := strings.ToLower(*input.TimeSkip)
		character.Race = race
	}

	if input.TimeSkip != nil {
		timeSkip := strings.ToLower(*input.TimeSkip)
		character.TimeSkip = timeSkip
	}

	updateIfNotNil(&character.Name, input.Name)
	updateIfNotNil(&character.Age, input.Age)
	updateIfNotNil(&character.Description, input.Description)
	updateIfNotNil(&character.Origin, input.Origin)
	updateIfNotNil(&character.Episode, input.Episode)

	if input.Bounty != nil {
		character.Bounty = input.Bounty
	}

	err = app.models.Characters.Update(character)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"character": character}, nil)
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

func (app *application) listCharactersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string
		Age      int
		Origin   string
		Race     string
		Bounty   *data.Berries
		TimeSkip string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Age = app.readInt(qs, "age", 0, v)
	input.Origin = app.readString(qs, "origin", "")
	input.Race = app.readString(qs, "race", "")
	// input.Bounty = app.readBounty

	//pagination shit
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "age", "bounty", "race"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	characters, err := app.models.Characters.GetAll(input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"characters": characters}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
