package db

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/null"
)

type Alliance struct {
	AllianceID    int32              `db:"alliance_id"`
	CreatorCorpID int32              `db:"creator_corporation_id"`
	CreatorID     int32              `db:"creator_id"`
	DateFounded   time.Time          `db:"date_founded"`
	ExecCorpID    null.JSONNullInt32 `db:"executor_corporation_id"`
	FactionID     null.JSONNullInt32 `db:"faction_id"`
	Name          string             `db:"name"`
	Ticker        string             `db:"ticker"`
}

// CopyAlliances upserts the given allianeces into the database using the copyAny function.
func (c *Client) CopyAlliances(ctx context.Context, alliances []Alliance) error {
	cols := []string{
		"alliance_id",
		"creator_corporation_id",
		"creator_id",
		"date_founded",
		"executor_corporation_id",
		"faction_id",
		"name",
		"ticker",
	}
	var anyAlliances [][]any

	setClauses := []string{
		"alliance_id",
		"creator_corporation_id = EXCLUDED.creator_corporation_id",
		"creator_id = EXCLUDED.creator_id",
		"date_founded = EXCLUDED.date_founded",
		"executor_corporation_id = EXCLUDED.executor_corporation_id",
		"faction_id = EXCLUDED.faction_id",
		"name = EXCLUDED.name",
		"ticker = EXCLUDED.ticker",
	}

	for _, alliance := range alliances {
		anyAlliances = append(anyAlliances, []any{
			alliance.AllianceID,
			alliance.CreatorCorpID,
			alliance.CreatorID,
			alliance.DateFounded,
			alliance.ExecCorpID,
			alliance.FactionID,
			alliance.Name,
			alliance.Ticker,
		})
	}

	today := time.Now().Format("20060102")

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		return copyAny(ctx, tx, "player", "alliance", today, cols, anyAlliances, setClauses...)
	}); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}

// GetMissingAllianceIDs returns a list of alliance IDs that are not in the database.
func (c *Client) GetMissingAllianceIDs(ctx context.Context, ids []int32) ([]int32, error) {
	queryVals := make([]string, 0, len(ids))

	for _, id := range ids {
		queryVals = append(queryVals, "("+strconv.Itoa(int(id))+")")
	}

	query := `
		SELECT
			missing
		FROM
			(
			VALUES
	` + strings.Join(queryVals, ",") + `
			) AS t(missing)
		WHERE
			missing NOT IN
			(SELECT
				alliance_id
			FROM
				player.alliance
			);
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var missing []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		missing = append(missing, id)
	}

	return missing, nil
}

// GetUnnamedAllianceIDs returns a list of alliance IDs that have no name.
func (c *Client) GetUnnamedAllianceIDs(ctx context.Context, count int) ([]int32, error) {
	query := `
		SELECT
			alliance_id
		FROM
			player.alliance
		WHERE
			name IS NULL OR name = ''
		LIMIT $1;
	`

	rows, err := c.pool.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	defer rows.Close()

	var unnamed []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		unnamed = append(unnamed, id)
	}

	return unnamed, nil
}

// InsertAlliance inserts the given alliance into the database. No ON CONFLICT clause is used.
func (c *Client) InsertAlliance(ctx context.Context, alliance Alliance) error {
	query := `
		INSERT INTO player.alliance
		(
			alliance_id,
			creator_corporation_id,
			creator_id,
			date_founded,
			executor_corporation_id,
			faction_id,
			name,
			ticker
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`

	if _, err := c.pool.Exec(ctx, query,
		alliance.AllianceID,
		alliance.CreatorCorpID,
		alliance.CreatorID,
		alliance.DateFounded,
		alliance.ExecCorpID,
		alliance.FactionID,
		alliance.Name,
		alliance.Ticker,
	); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}

// UpdateAlliance updates the given alliance in the database, using the alliance_id as the key.
func (c *Client) UpdateAlliance(ctx context.Context, alliance Alliance) error {
	query := `
		UPDATE player.alliance
		SET
			creator_corporation_id = $2,
			creator_id = $3,
			date_founded = $4,
			executor_corporation_id = $5,
			faction_id = $6,
			name = $7,
			ticker = $8
		WHERE
			alliance_id = $1;
	`

	if _, err := c.pool.Exec(ctx, query,
		alliance.AllianceID,
		alliance.CreatorCorpID,
		alliance.CreatorID,
		alliance.DateFounded,
		alliance.ExecCorpID,
		alliance.FactionID,
		alliance.Name,
		alliance.Ticker,
	); err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
