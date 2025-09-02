package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/05blue04/Poneglyph/internal/validator"
	"github.com/lib/pq"
)

type DevilFruit struct {
	ID             int64          `json:"id"`
	CreatedAt      time.Time      `json:"-"`
	UpdatedAt      time.Time      `json:"-"`
	Name           string         `json:"name"`
	Description    string         `json:"description"`
	Type           string         `json:"type"`
	CurrentOwner   sql.NullString `json:"-"`
	Character_id   sql.NullInt64  `json:"-"`
	PreviousOwners []string       `json:"previous_owners"`
	Episode        int            `json:"episode"`
}

type DevilFruitModel struct {
	DB *sql.DB
}

func (df DevilFruit) MarshalJSON() ([]byte, error) {
	var currentOwner *string
	if df.CurrentOwner.Valid {
		currentOwner = &df.CurrentOwner.String
	}

	var characterID *int64
	if df.Character_id.Valid {
		characterID = &df.Character_id.Int64
	}

	type Alias DevilFruit

	return json.Marshal(struct {
		Alias
		CurrentOwner *string `json:"current_owner"`
		CharacterID  *int64  `json:"character_id"`
	}{
		Alias:        Alias(df),
		CurrentOwner: currentOwner,
		CharacterID:  characterID,
	})
}

func ValidateDevilFruit(v *validator.Validator, devilFruit *DevilFruit) {

	validateName(v, "name", devilFruit.Name)
	validateDescription(v, devilFruit.Description)
	v.Check(devilFruit.Type != "", "type", "must be provided")
	v.Check(IsValidType(devilFruit.Type), "type", "must be a valid devil fruit type")

	v.Check(len(devilFruit.PreviousOwners) <= 10, "previous_owners", "must not have more than 10 previous owners")
	for i, owner := range devilFruit.PreviousOwners {
		v.Check(owner != "", "previous_owners", fmt.Sprintf("owner at index %d must not be empty", i))
		v.Check(len(owner) <= 200, "previous_owners", fmt.Sprintf("owner name at index %d must not be more than 200 characters", i))
		v.Check(utf8.ValidString(owner), "previous_owners", fmt.Sprintf("owner name at index %d must be valid UTF-8", i))
	}
	v.Check(validator.Unique(devilFruit.PreviousOwners), "previous_owners", "must not contain duplicates")

	validateEpisode(v, devilFruit.Episode)

}

func (m DevilFruitModel) Insert(devilFruit *DevilFruit) error {
	query := `
		INSERT INTO devilfruits (name, description, type, character_id, current_owner, previousOwners, episode)
    	VALUES ($1, $2, $3, $4, $5, $6, $7)
 		RETURNING id
	`

	args := []any{devilFruit.Name, devilFruit.Description, devilFruit.Type, devilFruit.Character_id, devilFruit.CurrentOwner, pq.Array(devilFruit.PreviousOwners), devilFruit.Episode}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&devilFruit.ID)
}

func (m DevilFruitModel) Get(id int64) (*DevilFruit, error) {
	query := `
		SELECT * FROM devilfruits
		WHERE id = $1
	`

	var devilFruit DevilFruit

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&devilFruit.ID,
		&devilFruit.CreatedAt,
		&devilFruit.UpdatedAt,
		&devilFruit.Name,
		&devilFruit.Description,
		&devilFruit.Type,
		&devilFruit.CurrentOwner,
		&devilFruit.Character_id,
		pq.Array(&devilFruit.PreviousOwners),
		&devilFruit.Episode,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &devilFruit, nil
}

func (m DevilFruitModel) Update(devilFruit *DevilFruit) error {
	query := `
		UPDATE devilfruits
		SET name = $1, description = $2, type = $3, character_id = $4, current_owner = $5, previousOwners = $6, episode = $7, updated_at = now()
		WHERE id = $8
	`

	args := []any{
		devilFruit.Name,
		devilFruit.Description,
		devilFruit.Type,
		devilFruit.Character_id,
		devilFruit.CurrentOwner,
		pq.Array(devilFruit.PreviousOwners),
		devilFruit.Episode,
		devilFruit.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
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

func (m DevilFruitModel) GetAll(search, fruitType string, filters Filters) ([]*DevilFruit, Metadata, error) {

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, created_at, name, description, type, character_id, ,previousOwners, episode 
		FROM devilfruits
		WHERE (to_tsvector('english', name || ' ' || description) @@ plainto_tsquery('english', $1) OR $1 = '')
		AND (LOWER(type) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, search, fruitType, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	//ensure that the resultset is closed before GetAll() returns.
	defer rows.Close()

	devilFruits := []*DevilFruit{}
	totalRecords := 0

	for rows.Next() {
		var devilFruit DevilFruit

		err := rows.Scan(
			&totalRecords,
			&devilFruit.ID,
			&devilFruit.CreatedAt,
			&devilFruit.Name,
			&devilFruit.Description,
			&devilFruit.Type,
			&devilFruit.Character_id,
			&devilFruit.CurrentOwner,
			pq.Array(&devilFruit.PreviousOwners),
			&devilFruit.Episode,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		devilFruits = append(devilFruits, &devilFruit)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return devilFruits, metadata, nil

}
