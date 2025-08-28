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

type CharacterModel struct {
	DB *sql.DB
}

func ValidateCharacter(v *validator.Validator, character *Character) {

	validateName(v, "name", character.Name)

	v.Check(character.Age > 0, "age", "must be a positive integer")
	v.Check(character.Age != 0, "age", "must be provided")

	validateDescription(v, character.Description)

	v.Check(character.Origin != "", "origin", "must be provided")
	v.Check(len(character.Origin) <= 200, "origin", "must not be more than 200 characters long")
	v.Check(utf8.ValidString(character.Origin), "origin", "must be valid UTF-8")

	//bounty validation
	if character.Bounty != nil {
		validateBounty(v, *character.Bounty)
	}

	//race validation
	v.Check(character.Race != "", "race", "must be provided")
	v.Check(IsValidRace(character.Race), "race", "must be a valid One Piece race")

	validateEpisode(v, character.Episode)
	validateTimeSkip(v, character.TimeSkip)

}

func (m CharacterModel) Insert(character *Character) error {
	query := `
		INSERT INTO characters (name, age, description, origin, bounty, race, episode, time_skip)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	var bounty sql.NullInt64
	if character.Bounty != nil {
		bounty = sql.NullInt64{Int64: int64(*character.Bounty), Valid: true}
	}

	args := []any{character.Name, character.Age, character.Description, character.Origin, bounty, character.Race, character.Episode, character.TimeSkip}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)

	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&character.ID)
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
	return deleteRecord(m.DB, "characters", id)
}

func (m CharacterModel) GetAll(search string, age int, origin, race string, bounty Berries, timeSkip string, filters Filters) ([]*Character, Metadata, error) {

	bountyCondition := "(bounty >= $5 OR $5 = 0)"

	if strings.Contains(filters.Sort, "bounty") {
		bountyCondition = "(bounty >= $5 AND bounty > 0)"
	}
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, name, age, description, origin, race, bounty, episode, time_skip
		FROM characters
		WHERE (to_tsvector('english', name || ' ' || description) @@ plainto_tsquery('english', $1) OR $1 = '')
		AND (LOWER(race) = LOWER($2) OR $2 = '')
		AND (LOWER(time_skip) = LOWER($3) OR $3 = '')
		AND (age >= $4 OR $4 = 0)
		AND %s
		ORDER BY %s %s, id ASC
		LIMIT $6 OFFSET $7`, bountyCondition, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, search, race, timeSkip, age, bounty, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	//ensure that the resultset is closed before GetAll() returns.
	defer rows.Close()

	characters := []*Character{}
	totalRecords := 0

	for rows.Next() {
		var character Character

		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}

		characters = append(characters, &character)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return characters, metadata, nil

}
