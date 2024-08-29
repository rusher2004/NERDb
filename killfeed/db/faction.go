package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/null"
)

type Faction struct {
	CorporationID      null.JSONNullInt32 `db:"corporation_id"`
	Description        string             `db:"description"`
	FactionID          int32              `db:"faction_id"`
	IsUnique           bool               `db:"is_unique"`
	MilitiaCorpID      null.JSONNullInt32 `db:"militia_corporation_id"`
	Name               string             `db:"name"`
	SizeFactor         float32            `db:"size_factor"`
	SolarSystemID      null.JSONNullInt32 `db:"solar_system_id"`
	StationCount       int32              `db:"station_count"`
	StationSystemCount int32              `db:"station_system_count"`
}

func (c *Client) CopyFactions(ctx context.Context, factions []Faction) error {
	cols := []string{
		"corporation_id",
		"description",
		"faction_id",
		"is_unique",
		"militia_corporation_id",
		"name",
		"size_factor",
		"solar_system_id",
		"station_system_count",
		"station_count",
	}

	var anyFactions [][]any

	setClauses := []string{
		"faction_id",
		"corporation_id = EXCLUDED.corporation_id",
		"description = EXCLUDED.description",
		"is_unique = EXCLUDED.is_unique",
		"militia_corporation_id = EXCLUDED.militia_corporation_id",
		"name = EXCLUDED.name",
		"size_factor = EXCLUDED.size_factor",
		"solar_system_id = EXCLUDED.solar_system_id",
		"station_count = EXCLUDED.station_count",
		"station_system_count = EXCLUDED.station_system_count",
	}

	for _, faction := range factions {
		anyFactions = append(anyFactions, []any{
			faction.CorporationID,
			faction.Description,
			faction.FactionID,
			faction.IsUnique,
			faction.MilitiaCorpID,
			faction.Name,
			faction.SizeFactor,
			faction.SolarSystemID,
			faction.StationSystemCount,
			faction.StationCount,
		})
	}

	today := time.Now().Format("20060102")

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		return copyAny(ctx, tx, "universe", "faction", today, cols, anyFactions, setClauses...)
	}); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}
