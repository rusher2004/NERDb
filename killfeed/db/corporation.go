package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/null"
)

type Corporation struct {
	AllianceID    null.JSONNullInt32  `db:"esi_alliance_id"`
	CEOID         null.JSONNullInt32  `db:"ceo_id"`
	CorporationID int32               `db:"esi_corporation_id"`
	CreatorID     null.JSONNullInt32  `db:"creator_id"`
	DateFounded   time.Time           `db:"date_founded"`
	Description   string              `db:"description"`
	FactionID     null.JSONNullInt32  `db:"faction_id"`
	HomeStationID null.JSONNullInt32  `db:"home_station_id"`
	MemberCount   int32               `db:"member_count"`
	Name          string              `db:"name"`
	Shares        int64               `db:"shares"`
	TaxRate       float32             `db:"tax_rate"`
	Ticker        string              `db:"ticker"`
	UpdatedAt     time.Time           `db:"updated_at"`
	URL           null.JSONNullString `db:"url"`
	WarEligible   bool                `db:"war_eligible"`
}

func (c *Client) CopyCorporations(ctx context.Context, corps []Corporation) error {
	cols := []string{
		"esi_corporation_id",
		"esi_alliance_id",
		"ceo_id",
		"creator_id",
		"date_founded",
		"description",
		"faction_id",
		"home_station_id",
		"member_count",
		"name",
		"shares",
		"tax_rate",
		"ticker",
		"url",
		"war_eligible",
	}
	var anyCorps [][]any

	setClauses := []string{
		"esi_corporation_id",
		"esi_alliance_id = EXCLUDED.esi_alliance_id",
		"ceo_id = EXCLUDED.ceo_id",
		"creator_id = EXCLUDED.creator_id",
		"date_founded = EXCLUDED.date_founded",
		"description = EXCLUDED.description",
		"faction_id = EXCLUDED.faction_id",
		"home_station_id = EXCLUDED.home_station_id",
		"member_count = EXCLUDED.member_count",
		"name = EXCLUDED.name",
		"shares = EXCLUDED.shares",
		"tax_rate = EXCLUDED.tax_rate",
		"ticker = EXCLUDED.ticker",
		"url = EXCLUDED.url",
		"war_eligible = EXCLUDED.war_eligible",
	}

	for _, corp := range corps {
		anyCorps = append(anyCorps, []any{
			corp.CorporationID,
			corp.AllianceID,
			corp.CEOID,
			corp.CreatorID,
			corp.DateFounded,
			corp.Description,
			corp.FactionID,
			corp.HomeStationID,
			corp.MemberCount,
			corp.Name,
			corp.Shares,
			corp.TaxRate,
			corp.Ticker,
			corp.URL,
			corp.WarEligible,
		})
	}

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		return copyAny(ctx, tx, "player", "corporation", "temp", cols, anyCorps, setClauses...)
	}); err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

func (c *Client) GetUnnamedCorporationIDs(ctx context.Context, count int) ([]int32, error) {
	query := `
		SELECT esi_corporation_id
		FROM player.corporation
		WHERE name IS NULL AND esi_corporation_id <> 0
		LIMIT $1;
	`

	rows, err := c.pool.Query(ctx, query, count)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}

	var corpIDs []int32
	for rows.Next() {
		var corpID int32
		if err := rows.Scan(&corpID); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		corpIDs = append(corpIDs, corpID)
	}

	return corpIDs, nil
}

func (c *Client) UpdateCorporation(ctx context.Context, corp Corporation) error {
	query := `
		UPDATE player.corporation
		SET
			esi_alliance_id = $2,
			ceo_id = $3,
			creator_id = $4,
			date_founded = $5,
			description = $6,
			faction_id = $7,
			home_station_id = $8,
			member_count = $9,
			name = $10,
			shares = $11,
			tax_rate = $12,
			ticker = $13,
			url = $14,
			war_eligible = $15
		WHERE esi_corporation_id = $1;
	`

	_, err := c.pool.Exec(ctx, query,
		corp.CorporationID,
		corp.AllianceID,
		corp.CEOID,
		corp.CreatorID,
		corp.DateFounded,
		corp.Description,
		corp.FactionID,
		corp.HomeStationID,
		corp.MemberCount,
		corp.Name,
		corp.Shares,
		corp.TaxRate,
		corp.Ticker,
		corp.URL,
		corp.WarEligible,
	)
	if err != nil {
		return fmt.Errorf("exec error: %w", err)
	}

	return nil
}
