package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/05blue04/Poneglyph/internal/data"
	"github.com/05blue04/Poneglyph/internal/validator"
)

func (app *application) createCrewHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		ShipName    string `json:"ship_name"`
		CaptainID   int64  `json:"captain_id"`
		Episode     int    `json:"episode"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	crew := &data.Crew{
		Name:        input.Name,
		Description: input.Description,
		ShipName:    input.ShipName,
		CaptainID:   input.CaptainID,
		// TotalBounty: data.Berries(0),
	}

	v := validator.New()

	if data.ValidateCrew(v, crew); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Crews.Insert(crew)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/crews/%d", crew.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"crew": crew}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showCrewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	crew, err := app.models.Crews.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"crew": crew}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateCrewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	crew, err := app.models.Crews.Get(id)
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
		Description *string       `json:"description"`
		ShipName    *string       `json:"ship_name"`
		CaptainID   *int64        `json:"captain_id"`
		TotalBounty *data.Berries `json:"total_bounty"`
		Episode     *int          `json:"episode"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateIfNotNil(&crew.Name, input.Name)
	updateIfNotNil(&crew.Description, input.Description)
	updateIfNotNil(&crew.ShipName, input.ShipName)
	updateIfNotNil(&crew.CaptainID, input.CaptainID)

	// if input.TotalBounty != nil {
	// 	crew.TotalBounty = *input.TotalBounty
	// }

	v := validator.New()

	if data.ValidateCrew(v, crew); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Crews.Update(crew)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"crew": crew}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteCrewHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Crews.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "crew successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
