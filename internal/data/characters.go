package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
)

type Character struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"`
	Fruit       string    `json:"devil_fruit"`
	Bounty      *Berries  `json:"bounty,omitempty"`
	Debut       string    `json:"debut"`
	//look to add a pre and post time skip field!
}

type CharacterModel struct {
	DB *sql.DB
}

func ValidateCharacter(v *validator.Validator, character *Character) {
	//Name validation
	v.Check(character.Name != "", "name", "must be provided")
	v.Check(len(character.Name) < 300, "name", "must not be more than 300 bytes long")
	v.Check(utf8.ValidString(character.Name), "name", "must be valid UTF-8")

	//Age validation
	v.Check(character.Age > 0, "age", "must be a positive integer")
	v.Check(character.Age != 0, "age", "must be provided")

	// Description validation
	v.Check(character.Description != "", "description", "must be provided")
	v.Check(len(character.Description) >= 10, "description", "must be at least 10 characters long")
	v.Check(len(character.Description) <= 2000, "description", "must not be more than 2000 characters long")
	v.Check(utf8.ValidString(character.Description), "description", "must be valid UTF-8")

	// Origin validation
	v.Check(character.Origin != "", "origin", "must be provided")
	v.Check(len(character.Origin) <= 200, "origin", "must not be more than 200 characters long")
	v.Check(utf8.ValidString(character.Origin), "origin", "must be valid UTF-8")

	//bounty validation
	if character.Bounty != nil {
		v.Check(*character.Bounty >= 0, "bounty", "must not be negative")
		v.Check(*character.Bounty <= 10000000000, "bounty", "must not exceed 10B berries")
		v.Check(*character.Bounty >= 1000, "bounty", "active bounties should be at least 1000 berries")
	}

}

func (m CharacterModel) Insert(character *Character) error {
	query := `
		INSERT INTO characters (name, age, description, origin, fruit, bounty, debut)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	var bounty sql.NullInt64
	if character.Bounty != nil {
		bounty = sql.NullInt64{Int64: int64(*character.Bounty), Valid: true}
	}

	args := []any{character.Name, character.Age, character.Description, character.Origin, character.Fruit, bounty, character.Debut}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m CharacterModel) Get(id int64) (*Character, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT * FROM characters
		WHERE id = $1
	`

	var character Character

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&character.ID,
		&character.CreatedAt,
		&character.Age,
		&character.Description,
		&character.Origin,
		&character.Fruit,
		&character.Bounty,
		&character.Debut,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &character, nil
}

func (m CharacterModel) Update(character *Character) error {
	return nil
}

func (m CharacterModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM characters
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)

	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
