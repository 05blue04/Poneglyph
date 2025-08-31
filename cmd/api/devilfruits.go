package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/05blue04/Poneglyph/internal/data"
	"github.com/05blue04/Poneglyph/internal/validator"
)

func (app *application) createDevilFruitHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name           string   `json:"name"`
		Description    string   `json:"description"`
		Type           string   `json:"type"`
		Character_id   int64    `json:"character_id"`
		PreviousOwners []string `json:"previous_owners"`
		Episode        int      `json:"episode"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	devilFruit := &data.DevilFruit{
		Name:           input.Name,
		Description:    input.Description,
		Type:           strings.ToLower(input.Type),
		Character_id:   input.Character_id,
		PreviousOwners: input.PreviousOwners,
		Episode:        input.Episode,
	}

	v := validator.New()

	if data.ValidateDevilFruit(v, devilFruit); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DevilFruit.Insert(devilFruit)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/devilfruits/%d", devilFruit.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"devilfruit": devilFruit}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showDevilFruitHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	devilFruit, err := app.models.DevilFruit.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"devilfruit": devilFruit}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateDevilFruitHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	devilFruit, err := app.models.DevilFruit.Get(id)
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
		Name           *string  `json:"name"`
		Description    *string  `json:"description"`
		Type           *string  `json:"type"`
		Character_id   *int64   `json:"character_id"`
		PreviousOwners []string `json:"previous_owners"`
		Episode        *int     `json:"episode"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateIfNotNil(&devilFruit.Name, input.Name)
	updateIfNotNil(&devilFruit.Description, input.Description)
	updateIfNotNil(&devilFruit.Type, input.Type)
	updateIfNotNil(&devilFruit.Character_id, input.Character_id)
	updateIfNotNil(&devilFruit.Episode, input.Episode)

	if input.PreviousOwners != nil {
		devilFruit.PreviousOwners = input.PreviousOwners
	}

	v := validator.New()

	if data.ValidateDevilFruit(v, devilFruit); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.DevilFruit.Update(devilFruit)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"devilfruit": devilFruit}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteDevilFruitHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.DevilFruit.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "devilfruit successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listDevilFruitsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Search string
		Type   string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Search = app.readString(qs, "search", "")
	input.Type = app.readString(qs, "type", "")
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "-id", "-name"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	devilFruits, metadata, err := app.models.DevilFruit.GetAll(input.Search, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"devil_fruits": devilFruits, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
