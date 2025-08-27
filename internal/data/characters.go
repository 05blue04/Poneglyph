package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
)

type Character struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Age         int       `json:"age"`
	Description string    `json:"description"`
	Origin      string    `json:"origin"`
	Bounty      *Berries  `json:"bounty,omitempty"`
	Race        string    `json:"race"`
	Episode     int       `json:"episode"`
	TimeSkip    string    `json:"time_skip"`
}

var validRaces = map[string]struct{}{
	"human":            {},
	"fishman":          {},
	"merman":           {},
	"giant":            {},
	"dwarf":            {},
	"mink":             {},
	"lunarian":         {},
	"buccaneer":        {},
	"long arm tribe":   {},
	"long leg tribe":   {},
	"snake neck tribe": {},
	"three-eye tribe":  {},
	"snakeneck tribe":  {},
	"longarm tribe":    {},
	"longleg tribe":    {},
	"tontatta":         {},
	"kuja":             {},
	"skypiean":         {},
	"shandian":         {},
	"birkan":           {},
	"cyborg":           {},
	"zombie":           {},
	"artificial human": {},
	"reindeer":         {}, // For Chopper
	"skeleton":         {}, // For Brook
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
	// potentially enforce origin creation to ensure the origin is an actual place in OP world

	//bounty validation
	if character.Bounty != nil {
		v.Check(*character.Bounty >= 0, "bounty", "must not be negative")
		v.Check(*character.Bounty <= 10000000000, "bounty", "must not exceed 10B berries")
		v.Check(*character.Bounty >= 100, "bounty", "active bounties should be at least 100 berries")
	}

	//race validation
	v.Check(character.Race != "", "race", "must be provided")
	v.Check(IsValidRace(character.Race), "race", "must be a valid One Piece race")

	//episode validation
	v.Check(character.Episode != 0, "episode", "must be provided")
	v.Check(character.Episode <= 1200, "episode", "must not be greater than 1200")
	v.Check(character.Episode > 0, "episode", "must not be negative")

	//time skip validationa
	v.Check(character.TimeSkip != "", "time_skip", "must be provided")
	v.Check(isValidTimeSkip(character.TimeSkip), "time_skip", "must be either pre or post")

}

func (m CharacterModel) Insert(character *Character) error {
	query := `
		INSERT INTO characters (name, age, description, origin, bounty, race, episode, time_skip)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	var bounty sql.NullInt64
	if character.Bounty != nil {
		bounty = sql.NullInt64{Int64: int64(*character.Bounty), Valid: true}
	}

	args := []any{character.Name, character.Age, character.Description, character.Origin, bounty, character.Race, character.Episode, character.TimeSkip}

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
		&character.UpdatedAt,
		&character.Name,
		&character.Age,
		&character.Description,
		&character.Origin,
		&character.Race,
		&character.Bounty,
		&character.Episode,
		&character.TimeSkip,
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
	query := `
		UPDATE characters
		SET name = $1, age = $2, description = $3, origin = $4, bounty = $5, race = $6,time_skip = $7, updated_at = now() 
		WHERE id = $8
	`
	var bounty sql.NullInt64
	if character.Bounty != nil {
		bounty = sql.NullInt64{Int64: int64(*character.Bounty), Valid: true}
	}

	args := []any{
		character.Name,
		character.Age,
		character.Description,
		character.Origin,
		bounty,
		character.Race,
		character.TimeSkip,
		character.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

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

func (m CharacterModel) GetAll(name string, age int, origin, race string, bounty Berries, timeSkip string, filters Filters) ([]*Character, error) {

	bountyCondition := "(bounty >= $5 OR $5 = 0)"

	if strings.Contains(filters.Sort, "bounty") {
		bountyCondition = "(bounty >= $5 AND bounty > 0)"
	}
	query := fmt.Sprintf(`
		SELECT id, created_at, name, age, description, origin, race, bounty, episode, time_skip
		FROM characters
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (LOWER(race) = LOWER($2) OR $2 = '')
		AND (LOWER(time_skip) = LOWER($3) OR $3 = '')
		AND (age >= $4 OR $4 = 0)
		AND %s
		ORDER BY %s %s, id ASC`, bountyCondition, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, name, race, timeSkip, age, bounty)
	if err != nil {
		return nil, err
	}

	//ensure that the resultset is closed before GetAll() returns.
	defer rows.Close()

	characters := []*Character{}

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&character.ID,
			&character.CreatedAt,
			&character.Name,
			&character.Age,
			&character.Description,
			&character.Origin,
			&character.Race,
			&character.Bounty,
			&character.Episode,
			&character.TimeSkip,
		)
		if err != nil {
			return nil, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return characters, nil

}

// IsValidRace checks if the provided race is a valid One Piece race
func IsValidRace(race string) bool {
	if race == "" {
		return false
	}

	_, exists := validRaces[race]
	return exists
}

// GetValidRaces returns a slice of all valid races (for API documentation, etc.)
func GetValidRaces() []string {
	races := make([]string, 0, len(validRaces))
	for race := range validRaces {
		// Capitalize first letter for display
		races = append(races, race)
	}
	return races
}

func isValidTimeSkip(timeSkip string) bool {
	if timeSkip == "pre" {
		return true
	}

	if timeSkip == "post" {
		return true
	}

	return false
}
