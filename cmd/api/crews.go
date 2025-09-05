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
		MemberCount: 1,
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

	err = app.models.Crews.AddMember(crew.ID, crew.CaptainID)
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

	characterID, err := app.readCharacterIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	character, err := app.models.Characters.Get(characterID)
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

func (app *application) listCrewMembersHandler(w http.ResponseWriter, r *http.Request) {
	crewID, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Crews.Get(crewID)
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
		Bounty data.Berries
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Bounty = app.readBounty(qs, "bounty", data.Berries(0), v)
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "bounty", "-id", "-name", "-bounty"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	members, metadata, err := app.models.Crews.GetMembers(crewID, input.Bounty, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"crew_members": members, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listCrewsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Search      string
		ShipName    string
		TotalBounty data.Berries
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Search = app.readString(qs, "search", "")
	input.TotalBounty = app.readBounty(qs, "total_bounty", data.Berries(0), v)
	input.ShipName = app.readString(qs, "ship_name", "")
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "bounty", "-id", "-name", "-total_bounty"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	crews, metadata, err := app.models.Crews.GetAll(input.Search, input.ShipName, input.TotalBounty, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"crews": crews, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
