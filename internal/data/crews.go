package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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
	MemberCount int       `json:"member_count"`
}

type CrewMember struct {
	ID     int64   `json:"id"`
	Name   string  `json:"name"`
	Bounty Berries `json:"bounty,omitempty"`
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

	countQuery := `
	    SELECT COUNT(*)
	    FROM crew_members
	    WHERE crew_id = $1
	`

	err = m.DB.QueryRowContext(ctx, countQuery, id).Scan(&crew.MemberCount)
	if err != nil {
		return nil, err
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

func (m CrewModel) AddMember(crewID, characterID int64) error {
	query := `
        INSERT INTO crew_members (character_id, crew_id)
        VALUES ($1, $2)
    `

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, characterID, crewID)
	if err != nil {
		return err
	}

	bountyQuery := `
        UPDATE crews
        SET total_bounty = total_bounty + COALESCE((SELECT bounty FROM characters WHERE id = $1), 0)
        WHERE id = $2
    `

	_, err = m.DB.ExecContext(ctx, bountyQuery, characterID, crewID)
	if err != nil {
		return err
	}

	return nil
}

func (m CrewModel) DeleteMember(crewID, characterID int64) error {
	query := `
		DELETE FROM crew_members
		WHERE crew_id = $1 AND character_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, crewID, characterID)
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

	bountyQuery := `
        UPDATE crews
        SET total_bounty = total_bounty - COALESCE((SELECT bounty FROM characters WHERE id = $1), 0)
        WHERE id = $2
    `
	_, err = m.DB.ExecContext(ctx, bountyQuery, crewID, characterID)
	if err != nil {
		return err
	}

	return nil
}

func (m CrewModel) GetMembers(crewID int64, bounty Berries, filters Filters) ([]*CrewMember, Metadata, error) {

	bountyCondition := "(c.bounty >= $2 OR $2 = 0)"

	if strings.Contains(filters.Sort, "bounty") {
		bountyCondition = "(c.bounty >= $2 AND c.bounty > 0)"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), c.id, c.name, c.bounty
		FROM characters c
		INNER JOIN crew_members cm ON c.id = cm.character_id
		WHERE cm.crew_id = $1
		AND %s
		ORDER BY %s %s, c.id ASC
		LIMIT $3 OFFSET $4`, bountyCondition, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	args := []any{crewID, bounty, filters.limit(), filters.offset()}

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var members []*CrewMember

	for rows.Next() {
		var member CrewMember
		var bounty *Berries

		err := rows.Scan(&totalRecords, &member.ID, &member.Name, &bounty)
		if err != nil {
			return nil, Metadata{}, err
		}

		if bounty != nil {
			member.Bounty = *bounty
		}

		members = append(members, &member)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return members, metadata, nil
}

func (m CrewModel) GetAll(search string, shipName string, totalBounty Berries, filters Filters) ([]*Crew, Metadata, error) {
	bountyCondition := "(total_bounty >= $3 OR $3 = 0)"
	if strings.Contains(filters.Sort, "total_bounty") {
		bountyCondition = "(total_bounty >= $3 AND total_bounty > 0)"
	}

	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), 
			c.id, c.created_at, c.updated_at, c.name, c.description, 
			c.ship_name, c.captain_id, c.captain_name, c.total_bounty
		FROM crews c
		WHERE (to_tsvector('english', c.name || ' ' || c.description || ' ' || COALESCE(c.ship_name, '') || ' ' || COALESCE(c.captain_name, '')) @@ plainto_tsquery('english', $1) OR $1 = '')
		AND (LOWER(c.ship_name) = LOWER($2) OR $2 = '' OR c.ship_name IS NULL)
		AND %s
		ORDER BY %s %s, c.id ASC
		LIMIT $4 OFFSET $5`, bountyCondition, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, search, shipName, totalBounty, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	crews := []*Crew{}
	totalRecords := 0

	for rows.Next() {
		var crew Crew
		err := rows.Scan(
			&totalRecords,
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
			return nil, Metadata{}, err
		}
		crews = append(crews, &crew)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return crews, metadata, nil
}
