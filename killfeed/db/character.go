package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/null"
)

type Gender string

const (
	GenderFemale = "female"
	GenderMale   = "male"
)

type Character struct {
	AllianceID     null.JSONNullInt32  `db:"esi_alliance_id"`
	Birthday       time.Time           `db:"birthday"`
	CreatedAt      time.Time           `db:"created_at"`
	BloodlineID    int32               `db:"bloodline_id"`
	CharacterID    int32               `db:"esi_character_id"`
	CorporationID  int32               `db:"esi_corporation_id"`
	Description    null.JSONNullString `db:"description"`
	Gender         Gender              `db:"gender"`
	FactionID      null.JSONNullInt32  `db:"faction_id"`
	Name           string              `db:"name"`
	RaceID         int32               `db:"race_id"`
	SecurityStatus float32             `db:"security_status"`
	Title          null.JSONNullString `db:"title"`
	UpdatedAt      time.Time           `db:"updated_at"`
}

// CopyCharacters uses Postgres CopyFrom to copy characters into the database, copying into a temp
// table, then inserting into the persistent table.
func (c *Client) CopyCharacters(ctx context.Context, characters []Character) error {
	cols := []string{
		"esi_character_id",
		"esi_alliance_id",
		"birthday",
		"bloodline_id",
		"esi_corporation_id",
		"description",
		"faction_id",
		"gender",
		"name",
		"race_id",
		"security_status",
	}
	var anyChars [][]any

	setClauses := []string{
		"esi_character_id",
		"esi_alliance_id = EXCLUDED.esi_alliance_id",
		"birthday = EXCLUDED.birthday",
		"bloodline_id = EXCLUDED.bloodline_id",
		"esi_corporation_id = EXCLUDED.esi_corporation_id",
		"description = EXCLUDED.description",
		"faction_id = EXCLUDED.faction_id",
		"gender = EXCLUDED.gender",
		"name = EXCLUDED.name",
		"race_id = EXCLUDED.race_id",
		"security_status = EXCLUDED.security_status",
	}

	for _, c := range characters {
		anyChars = append(anyChars, []any{
			c.CharacterID,
			c.AllianceID,
			c.Birthday,
			c.BloodlineID,
			c.CorporationID,
			c.Description,
			c.FactionID,
			c.Gender,
			c.Name,
			c.RaceID,
			c.SecurityStatus,
		})
	}

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		return copyAny(ctx, tx, "player", "character", "temp", cols, anyChars, setClauses...)
	}); err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

func (c *Client) CopyDeletedCharacters(ctx context.Context, charIDs []int32) error {
	cols := []string{"esi_character_id"}
	var anyChars [][]any

	setClauses := []string{
		"esi_character_id",
		"esi_deleted = true",
	}

	for _, id := range charIDs {
		anyChars = append(anyChars, []any{id})
	}

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		return copyAny(ctx, tx, "player", "character", "temp", cols, anyChars, setClauses...)
	}); err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

// GetUnnamedCharacterIDs returns a list of character IDs where name and esi_deleted are null
func (c *Client) GetUnnamedCharacterIDs(ctx context.Context, count int) ([]int32, error) {
	query := `
		SELECT esi_character_id
		FROM player.character
		// -- WHERE (name IS NULL OR LOWER(name) = 'unknown') AND esi_deleted IS NOT TRUE
		WHERE name IS NULL AND esi_deleted IS NOT TRUE
		LIMIT $1;
	`

	rows, err := c.pool.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var charIDs []int32
	for rows.Next() {
		var charID int32
		if err := rows.Scan(&charID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		charIDs = append(charIDs, charID)
	}

	return charIDs, nil
}

// SetCharacterDeleted sets the esi_deleted flag to true for a character. This is the case when ESI
// has returned a message that the character has been deleted.
func (c *Client) SetCharacterDeleted(ctx context.Context, charID int32) error {
	query := `
		UPDATE player.character
		SET esi_deleted = true
		WHERE esi_character_id = $1;
	`

	if _, err := c.pool.Exec(ctx, query, charID); err != nil {
		return err
	}

	return nil
}

// UpdateCharacter updates the character. No checks on the values or conflicts are done, so the
// values in char will all be updated.
func (c *Client) UpdateCharacter(ctx context.Context, char Character) error {
	query := `
		UPDATE player.character
		SET
			birthday = $1,
			bloodline_id = $2,
			esi_alliance_id = $3,
			esi_corporation_id = $4,
			description = $5,
			faction_id = $6,
			gender = $7,
			name = $8,
			race_id = $9,
			security_status = $10,
			title = $11
		WHERE esi_character_id = $12;
	`

	if _, err := c.pool.Exec(ctx, query,
		char.Birthday,
		char.BloodlineID,
		char.AllianceID,
		char.CorporationID,
		char.Description,
		char.FactionID,
		char.Gender,
		char.Name,
		char.RaceID,
		char.SecurityStatus,
		char.Title,
		char.CharacterID,
	); err != nil {
		return err
	}

	return nil
}
