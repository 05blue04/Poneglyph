package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/05blue04/Poneglyph/internal/validator"
)

type Crew struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ShipName    string    `json:"ship_name"`
	CaptainID   int64     `json:"captain_id"`
	CaptainName string    `json:"captain_name"`
	TotalBounty Berries   `json:"total_bounty"`
}

type CrewModel struct {
	DB *sql.DB
}

func ValidateCrew(v *validator.Validator, crew *Crew) {
	validateName(v, "name", crew.Name)
	validateDescription(v, crew.Description)
	validateName(v, "ship_name", crew.ShipName)

	v.Check(crew.CaptainID > 0, "captain_id", "must be greater than 0")

}

func (m CrewModel) Insert(crew *Crew) error {
	query := `
		INSERT INTO crews (name, description, ship_name, captain_id, captain_name, total_bounty)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	args := []any{
		crew.Name,
		crew.Description,
		crew.ShipName,
		crew.CaptainID,
		crew.CaptainName,
		crew.TotalBounty,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&crew.ID)
}

func (m CrewModel) Get(id int64) (*Crew, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT * FROM crews
		WHERE id = $1
	`

	var crew Crew

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&crew.ID,
		&crew.CreatedAt,
		&crew.UpdatedAt,
		&crew.Name,
		&crew.Description,
		&crew.ShipName,
		&crew.CaptainID,
		&crew.CaptainName,
		&crew.TotalBounty,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &crew, nil
}

func (m CrewModel) Update(crew *Crew) error {
	query := `
		UPDATE crews
		SET name = $1, description = $2, ship_name = $3, captain_id = $4, captain_name = $5, total_bounty = $6, updated_at = now()
		WHERE id = $7
	`
	args := []any{
		crew.Name,
		crew.Description,
		crew.ShipName,
		crew.CaptainID,
		crew.CaptainName,
		crew.TotalBounty,
		crew.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m CrewModel) Delete(id int64) error {
	return deleteRecord(m.DB, "crews", id)
}
