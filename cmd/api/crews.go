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
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character, err := app.models.Characters.Get(input.CaptainID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "character_id does not exist")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var total_bounty data.Berries

	if character.Bounty != nil {
		total_bounty = *character.Bounty
	}

	crew := &data.Crew{
		Name:        input.Name,
		Description: input.Description,
		ShipName:    input.ShipName,
		CaptainID:   input.CaptainID,
		CaptainName: character.Name,
		TotalBounty: total_bounty,
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
		Name        *string `json:"name"`
		Description *string `json:"description"`
		ShipName    *string `json:"ship_name"`
		CaptainID   *int64  `json:"captain_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	updateIfNotNil(&crew.Name, input.Name)
	updateIfNotNil(&crew.Description, input.Description)
	updateIfNotNil(&crew.ShipName, input.ShipName)

	if input.CaptainID != nil {
		newCaptain, err := app.models.Characters.Get(*input.CaptainID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.errorResponse(w, r, http.StatusUnprocessableEntity, "character_id does not exist")
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		oldCaptain, err := app.models.Characters.Get(crew.CaptainID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		var oldBounty data.Berries
		var newBounty data.Berries

		if oldCaptain.Bounty != nil {
			oldBounty = *oldCaptain.Bounty
		}

		if newCaptain.Bounty != nil {
			newBounty = *newCaptain.Bounty
		}

		newTotalBounty := (crew.TotalBounty - oldBounty) + newBounty

		crew.CaptainID = newCaptain.ID
		crew.CaptainName = newCaptain.Name
		crew.TotalBounty = newTotalBounty
	}

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

func (app *application) addCrewMemberHandler(w http.ResponseWriter, r *http.Request) {
	crewID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		CharacterID int64 `json:"character_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character, err := app.models.Characters.Get(input.CharacterID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "character_id does not exist")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Crews.AddMember(crewID, character.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (app *application) deleteCrewMemberHandler(w http.ResponseWriter, r *http.Request) {
	crewID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		CharacterID int64 `json:"character_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	character, err := app.models.Characters.Get(input.CharacterID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "character_id does not exist")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Crews.DeleteMember(crewID, character.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, fmt.Sprintf("character %v is not a member of this crew", character.Name))
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": fmt.Sprintf("crew member %v successfully deleted", character.Name)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
