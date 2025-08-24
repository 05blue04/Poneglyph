package data

import (
	"database/sql"
	"time"
)

type Character struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"`
	Fruit       string    `json:"devil_fruit"`
	Bounty      string    `json:"bounty"`
	Debut       string    `json:"debut"`
}

type CharacterModel struct {
	DB *sql.DB
}
