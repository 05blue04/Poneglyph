package data

import (
	"context"
	"database/sql"
	"time"
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
)

type DevilFruit struct {
	ID             int64     `json:"id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	Character_id   int64     `json:"character_id"`
	PreviousOwners string    `json:"previous_owners"`
	Episode        int       `json:"episode"`
}

type DevilFruitModel struct {
	DB *sql.DB
}

var validTypes = map[string]struct{}{
	"zoan":      {},
	"paramecia": {},
	"logia":     {},
}

func ValidateDevilFruit(v *validator.Validator, devilFruit *DevilFruit) {

	//name validation
	v.Check(devilFruit.Name != "", "name", "must be provided")
	v.Check(len(devilFruit.Name) < 300, "name", "must not be more than 300 bytes long")
	v.Check(utf8.ValidString(devilFruit.Name), "name", "must be valid UTF-8")

	//description validation
	v.Check(devilFruit.Description != "", "description", "must be provided")
	v.Check(len(devilFruit.Description) >= 10, "description", "must be at least 10 characters long")
	v.Check(len(devilFruit.Description) <= 2000, "description", "must not be more than 2000 characters long")
	v.Check(utf8.ValidString(devilFruit.Description), "description", "must be valid UTF-8")

	//devilfruit type validation
	v.Check(devilFruit.Type != "", "type", "must be provided")
	v.Check(IsValidType(devilFruit.Type), "type", "must be a valid devil fruit type")

	//add previous owners validation

	//character_id validation
	v.Check(devilFruit.Character_id >= 0, "character_id", "must be a positive integer")

	//episode validation
	v.Check(devilFruit.Episode != 0, "episode", "must be provided")
	v.Check(devilFruit.Episode <= 1200, "episode", "must not be greater than 1200")
	v.Check(devilFruit.Episode > 0, "episode", "must not be negative")
}

func (m DevilFruitModel) Insert(devilFruit *DevilFruit) error {
	return nil
}

func (m DevilFruitModel) Get(id int64) (*DevilFruit, error) {
	return nil, nil
}

func (m DevilFruitModel) Update(devilFruit *DevilFruit) error {
	return nil
}

func (m DevilFruitModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM devilfruits
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

func (m DevilFruitModel) GetAll(args ...any) ([]*DevilFruit, Metadata, error)

func IsValidType(devilFruitType string) bool {
	if devilFruitType == "" {
		return false
	}

	_, exists := validTypes[devilFruitType]

	return exists
}
