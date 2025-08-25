package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Characters CharacterModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Characters: CharacterModel{DB: db},
	}
}
