package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict") //this error handles the case where a race condition occurs
)

type Models struct {
	Characters CharacterModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Characters: CharacterModel{DB: db},
	}
}
