package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Characters  CharacterModel
	DevilFruits DevilFruitModel
	Crews       CrewModel
	APIKeys     APIKeyModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Characters:  CharacterModel{DB: db},
		DevilFruits: DevilFruitModel{DB: db},
		Crews:       CrewModel{DB: db},
		APIKeys:     APIKeyModel{DB: db},
	}
}
