package db

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/killfeed/esi"
	"github.com/rusher2004/nerdb/killfeed/zkill"
)

// toAnySlice converts a slice of any type to a slice of any type https://go.dev/doc/faq#convert_slice_of_interface
func toAnySlice[T any](i []T) []any {
	out := make([]any, len(i))
	for i, v := range i {
		out[i] = v
	}

	return out
}

func insertAttackers(ctx context.Context, p DBPool, killmailID int, attackers []esi.KillmailAttacker) error {
	query := `
	INSERT INTO killmail.attacker(esi_killmail_id, esi_character_id, esi_corporation_id, esi_alliance_id, esi_faction_id, damage_done, final_blow, security_status, ship_type_id, weapon_type_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	for _, a := range attackers {
		tag, err := p.Exec(ctx, query, killmailID, a.CharacterID, a.CorporationID, a.AllianceID, a.FactionID, a.DamageDone, a.FinalBlow, a.SecurityStatus, a.ShipTypeID, a.WeaponTypeID)
		if err != nil {
			return fmt.Errorf("error inserting attacker: %w", err)
		}

		if tag.RowsAffected() == 0 {
			log.Printf("no rows affected: insert attacker character_id: %d, killmail_id: %d\n", a.CharacterID.Int32, killmailID)
		}
	}

	return nil
}

func insertVictim(ctx context.Context, p DBPool, killmailID int, v esi.KillmailVictim) error {
	query := `
	INSERT INTO killmail.victim(esi_killmail_id, esi_character_id, esi_corporation_id, esi_alliance_id, esi_faction_id, ship_type_id, damage_taken, position_x, position_y, position_z)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`
	tag, err := p.Exec(ctx, query, killmailID, v.CharacterID, v.CorporationID, v.AllianceID, v.FactionID, v.ShipTypeID, v.DamageTaken, v.Position.X, v.Position.Y, v.Position.Z)
	if err != nil {
		return fmt.Errorf("error inserting victim: %w", err)
	}

	if tag.RowsAffected() == 0 {
		log.Printf("no rows affected: insert victimg character_id: %d, killmail_id: %d\n", v.CharacterID.Int32, killmailID)
		return nil
	}

	if err := insertVictimItems(ctx, p, killmailID, "victim", v.Items); err != nil {
		return fmt.Errorf("error inserting victim items: %w", err)
	}

	return nil
}

func insertVictimItems(ctx context.Context, p DBPool, parentID int, parentType string, items []esi.KillmailVictimItem) error {
	query := `
	INSERT INTO killmail.victim_item(parent_id, parent_type, flag, item_type_id, quantity_destroyed, quantity_dropped, singleton)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING victim_item_id;
	`

	for _, i := range items {
		row := p.QueryRow(ctx, query, parentID, parentType, i.Flag, i.ItemTypeID, i.QuantityDestroyed, i.QuantityDropped, i.Singleton)

		var victimItemID int
		if err := row.Scan(&victimItemID); err != nil {
			return fmt.Errorf("error inserting victim item: %w", err)
		}

		if i.Items != nil {
			if err := insertVictimItems(ctx, p, victimItemID, "item", *i.Items); err != nil {
				return fmt.Errorf("error inserting victim item details: %w", err)
			}
		}
	}

	return nil
}

func upsertEntities(ctx context.Context, p DBPool, entities []int, entityType string) error {
	if len(entities) == 0 {
		return nil
	}

	query := fmt.Sprintf("INSERT INTO player.%s(esi_%s_id) VALUES ", entityType, entityType)

	var placeholders []string
	for i := range entities {
		placeholders = append(placeholders, fmt.Sprintf("($%d)", i+1))
	}

	query += strings.Join(placeholders, ", ") + " ON CONFLICT (esi_" + entityType + "_id) DO NOTHING;"

	if _, err := p.Exec(ctx, query, toAnySlice(entities)...); err != nil {
		return fmt.Errorf("execution error: %w", err)
	}

	return nil
}

func upsertParticipants(ctx context.Context, p DBPool, participants []esi.KillMailParticipant) error {
	var allianceIDs, charIDs, corpIDs []int
	for _, p := range participants {
		if p.AllianceID.Valid && p.AllianceID.Int32 != 0 && !slices.Contains(allianceIDs, int(p.AllianceID.Int32)) {
			allianceIDs = append(allianceIDs, int(p.AllianceID.Int32))
		}

		if !slices.Contains(charIDs, int(p.CharacterID.Int32)) {
			charIDs = append(charIDs, int(p.CharacterID.Int32))
		}

		if p.CorporationID.Valid && p.CorporationID.Int32 != 0 && !slices.Contains(corpIDs, int(p.CorporationID.Int32)) {
			corpIDs = append(corpIDs, int(p.CorporationID.Int32))
		}
	}

	if err := upsertEntities(ctx, p, charIDs, "character"); err != nil {
		return fmt.Errorf("error upserting esi character: %w", err)
	}

	if err := upsertEntities(ctx, p, corpIDs, "corporation"); err != nil {
		return fmt.Errorf("error upserting esi corporation: %w", err)
	}

	if err := upsertEntities(ctx, p, allianceIDs, "alliance"); err != nil {
		return fmt.Errorf("error upserting esi alliance: %w", err)
	}

	return nil
}

func upsertKillmail(ctx context.Context, p DBPool, k esi.Killmail) error {
	query := `
	INSERT INTO killmail.esi_killmail(esi_killmail_id, time, moon_id, solar_system_id, war_id)
	VALUES ($1, $2, $3, $4, $5)
	ON CONFLICT (esi_killmail_id) DO UPDATE SET moon_id = EXCLUDED.moon_id, war_id = EXCLUDED.war_id;
	`
	tag, err := p.Exec(ctx, query, k.KillmailID, k.KillmailTime, k.MoonID, k.SolarSystemID, k.WarID)
	if err != nil {
		return fmt.Errorf("error inserting killmail: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func upsertZkillInfo(ctx context.Context, p DBPool, killmailID int, z zkill.ZKKillInfo) error {
	query := `
	INSERT INTO killmail.zkill_info(esi_killmail_id, awox, destroyed_value, dropped_value, fitted_value, hash, location_id, npc, points, solo, total_value)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	ON CONFLICT (esi_killmail_id) DO UPDATE SET awox = EXCLUDED.awox, destroyed_value = EXCLUDED.destroyed_value, dropped_value = EXCLUDED.dropped_value, fitted_value = EXCLUDED.fitted_value, hash = EXCLUDED.hash, location_id = EXCLUDED.location_id, npc = EXCLUDED.npc, points = EXCLUDED.points, solo = EXCLUDED.solo, total_value = EXCLUDED.total_value;
	`

	tag, err := p.Exec(
		ctx,
		query,
		killmailID,
		z.Awox,
		z.DestroyedValue,
		z.DroppedValue,
		z.FittedValue,
		z.Hash,
		z.LocationID,
		z.NPC,
		z.Points,
		z.Solo,
		z.TotalValue,
	)

	if err != nil {
		return fmt.Errorf("error inserting zkill info: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (c *Client) UpsertESIKillmail(ctx context.Context, k esi.Killmail) error {
	participants := k.UniqueParticipants()
	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		if err := upsertParticipants(ctx, c.pool, participants); err != nil {
			return fmt.Errorf("error upserting participants: %w", err)
		}

		if err := upsertKillmail(ctx, c.pool, k); err != nil {
			return fmt.Errorf("error upserting killmail: %w", err)
		}

		if err := insertVictim(ctx, c.pool, int(k.KillmailID), k.Victim); err != nil {
			return fmt.Errorf("error inserting victim: %w", err)
		}

		if err := insertAttackers(ctx, c.pool, int(k.KillmailID), k.Attackers); err != nil {
			return fmt.Errorf("error inserting attackers: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

func (c *Client) UpsertZkillKillmail(ctx context.Context, z zkill.RedisQPackage) error {
	if err := c.UpsertESIKillmail(ctx, z.ESIKill); err != nil {
		return fmt.Errorf("error upserting ESI killmail: %w", err)
	}

	if err := upsertZkillInfo(ctx, c.pool, int(z.ESIKill.KillmailID), z.ZKKill); err != nil {
		return fmt.Errorf("error inserting zkill info: %w", err)
	}

	return nil
}
