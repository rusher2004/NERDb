package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/killfeed/esi"
	"github.com/rusher2004/nerdb/killfeed/zkill"
)

type attackerWithKillmailID struct {
	killmailID int
	esi.KillmailAttacker
}

type victimWithKillmailID struct {
	killmailID int
	esi.KillmailVictim
}

// commenting out but leaving in, may figure out soon how to efficiently store this.
// type victimItemWithParentID struct {
// 	esi.KillmailVictimItem
// 	parentType string
// 	parentID   int
// }

// toAnySlice converts a slice of any type to a slice of any type https://go.dev/doc/faq#convert_slice_of_interface
func toAnySlice[T any](i []T) []any {
	out := make([]any, len(i))
	for i, v := range i {
		out[i] = v
	}

	return out
}

// copyAny uses Postgres CopyFrom to copy data into table `schema.table`. The values in cols and rows
// must be ordered correctly. As a hacky temporary solution to allow for handling ON CONFLICT,
// setClauses can be provided. If used, seClauses[0] must be the column(s) to conflict on, and the
// rest of the setClauses must be the columns to update, formated as `key = value`.
func copyAny(ctx context.Context, tx pgx.Tx, schema, table, date string, cols []string, rows [][]any, setClauses ...string) error {
	log.Printf("copying %d %s(s)\n", len(rows), table)
	if len(rows) == 0 {
		return nil
	}

	b, err := tx.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	if err := createTempTable(ctx, b, schema, table, date); err != nil {
		return fmt.Errorf("error creating temp table: %w", err)
	}

	count, err := copyToTempTable(ctx, b, schema, table+"_"+date, cols,
		pgx.CopyFromSlice(
			len(rows),
			func(i int) ([]any, error) {
				return rows[i], nil
			},
		),
	)
	if err != nil {
		return fmt.Errorf("error copying %s: %w", table, err)
	}

	toTable := schema + "." + table
	fromTable := schema + "_" + table + "_" + date
	var insertErr error
	switch {
	case len(setClauses) == 0:
		if err := insertFromTempTable(ctx, b, toTable, fromTable); err != nil {
			insertErr = fmt.Errorf("error inserting %s: %w", table, err)
		}
	case len(setClauses) >= 1:
		if err := inserFromTempTableWithConflict(ctx, b, toTable, fromTable, setClauses[0], setClauses[1:]...); err != nil {
			insertErr = fmt.Errorf("error inserting %s: %w", table, err)
		}
	default:
		return errors.New("invalid setClauses")
	}

	if insertErr != nil {
		return insertErr
	}

	log.Printf("copied %d %s(s)\n", count, table)

	return nil
}

// copyAttackers uses copyAny to copy attackers into the database.
func copyAttackers(ctx context.Context, tx pgx.Tx, date string, attackers []attackerWithKillmailID) error {
	var anyAttackers [][]any
	for _, a := range attackers {
		anyAttackers = append(anyAttackers, []any{a.killmailID, a.CharacterID, a.CorporationID, a.AllianceID, a.FactionID, a.DamageDone, a.FinalBlow, a.SecurityStatus, a.ShipTypeID, a.WeaponTypeID})
	}

	cols := []string{"esi_killmail_id", "esi_character_id", "esi_corporation_id", "esi_alliance_id", "esi_faction_id", "damage_done", "final_blow", "security_status", "ship_type_id", "weapon_type_id"}

	return copyAny(ctx, tx, "killmail", "attacker", date, cols, anyAttackers)
}

// copyKillmails uses copyAny to copy killmails into the database.
func copyKillmails(ctx context.Context, tx pgx.Tx, date string, kms []esi.Killmail) error {
	var anyKms [][]any
	for _, k := range kms {
		anyKms = append(anyKms, []any{k.KillmailID, k.KillmailTime, k.MoonID, k.SolarSystemID, k.WarID})
	}

	cols := []string{"esi_killmail_id", "time", "moon_id", "solar_system_id", "war_id"}

	return copyAny(ctx, tx, "killmail", "esi_killmail", date, cols, anyKms)
}

// copyParticipants uses copyAny to copy participants into the database.
func copyParticipants(ctx context.Context, tx pgx.Tx, date string, kms []esi.Killmail) error {
	var participants []esi.KillMailParticipant
	for _, k := range kms {
		participants = append(participants, k.UniqueParticipants()...)
	}

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

	var anyChars [][]any
	for _, c := range charIDs {
		anyChars = append(anyChars, []any{c})
	}
	if err := copyAny(ctx, tx, "player", "character", date, []string{"esi_character_id"}, anyChars); err != nil {
		return fmt.Errorf("error copying characters: %w", err)
	}

	var anyCorps [][]any
	for _, c := range corpIDs {
		anyCorps = append(anyCorps, []any{c})
	}
	if err := copyAny(ctx, tx, "player", "corporation", date, []string{"esi_corporation_id"}, anyCorps); err != nil {
		return fmt.Errorf("error copying corporations: %w", err)
	}

	var anyAlliances [][]any
	for _, a := range allianceIDs {
		anyAlliances = append(anyAlliances, []any{a})
	}
	if err := copyAny(ctx, tx, "player", "alliance", date, []string{"esi_alliance_id"}, anyAlliances); err != nil {
		return fmt.Errorf("error copying alliances: %w", err)
	}

	return nil
}

// copyToTempTable uses Postgres CopyFrom to copy data into a temp table.
func copyToTempTable(ctx context.Context, p pgx.Tx, schema, table string, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error) {
	return p.CopyFrom(ctx, pgx.Identifier{schema + "_" + table}, columnNames, rowSrc)
}

// copyVictims uses copyAny to copy victims into the database.
func copyVictims(ctx context.Context, tx pgx.Tx, date string, victims []victimWithKillmailID) error {
	// var allVictimItems []victimItemWithParentID
	// for _, v := range victims {
	// 	for _, i := range v.Items {
	// 		allVictimItems = append(allVictimItems, victimItemWithParentID{i, "victim", int(v.CharacterID.Int32)})
	// 	}
	// }

	var anyVictims [][]any
	for _, v := range victims {
		var x, y, z *float64
		if v.Position != nil {
			x = &v.Position.X
			y = &v.Position.Y
			z = &v.Position.Z
		}
		anyVictims = append(anyVictims, []any{v.killmailID, v.CharacterID, v.CorporationID, v.AllianceID, v.FactionID, v.ShipTypeID, v.DamageTaken, x, y, z})
	}

	cols := []string{"esi_killmail_id", "esi_character_id", "esi_corporation_id", "esi_alliance_id", "esi_faction_id", "ship_type_id", "damage_taken", "position_x", "position_y", "position_z"}

	if err := copyAny(ctx, tx, "killmail", "victim", date, cols, anyVictims); err != nil {
		return fmt.Errorf("error copying victim: %w", err)
	}

	// if err := copyVictimItems(ctx, tx, date, allVictimItems); err != nil {
	// 	return fmt.Errorf("error copying victim items: %w", err)
	// }

	return nil
}

// func copyVictimItems(ctx context.Context, tx pgx.Tx, date string, items []victimItemWithParentID) error {
// 	var anyItems [][]any
// 	for _, i := range items {
// 		anyItems = append(anyItems, []any{i.parentID, i.parentType, i.Flag, i.ItemTypeID, i.QuantityDestroyed, i.QuantityDropped, i.Singleton})
// 	}

// 	cols := []string{"parent_id", "parent_type", "flag", "item_type_id", "quantity_destroyed", "quantity_dropped", "singleton"}

// 	if err := copyAny(ctx, tx, "killmail", "victim_item", date, cols, anyItems); err != nil {
// 		return fmt.Errorf("error copying victim items: %w", err)
// 	}

// 	// need to figure out how to reconcile nested items later. We need the victim_item_id of its
// 	// parent to insert a nested item, which is not available until the parent item is inserted.
// 	// for _, i := range items {
// 	// 	if i.Items != nil {
// 	// 		if err := copyVictimItems(ctx, tx, date, parentID, "item", *i.Items); err != nil {
// 	// 			return fmt.Errorf("error copying victim item details: %w", err)
// 	// 		}
// 	// 	}
// 	// }

// 	return nil
// }

// createTempTable creates a temporary table in the database, using `fromSchema.fromTable` as the
// model for the temp table using LIKE.
func createTempTable(ctx context.Context, tx pgx.Tx, fromSchema, fromTable, date string) error {
	tableName := fromSchema + "_" + fromTable + "_" + date
	query := `
	CREATE TEMP TABLE IF NOT EXISTS ` + tableName + `
	(LIKE ` + fromSchema + "." + fromTable + ` INCLUDING DEFAULTS)
	ON COMMIT DROP;
	`

	_, err := tx.Exec(ctx, query)

	return err
}

// inserAttackers inserts attackers into the database.
func insertAttackers(ctx context.Context, p DBPool, killmailID int, attackers []esi.KillmailAttacker) error {
	// TODO: this was created before copyAny. Let's update to use that to get a little better performance and use fewer transactions.
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

// insertFromTempTable inserts all records from `from` into `to`. On conflict, nothing is done.
func insertFromTempTable(ctx context.Context, tx pgx.Tx, to, from string) error {
	query := `
	INSERT INTO ` + to + `
	SELECT * FROM ` + from + `
	ON CONFLICT DO NOTHING;
	`

	_, err := tx.Exec(ctx, query)

	return err
}

// inserFromTempTableWithConflict inserts all records from `from` into `to`. On conflict of
// conflictCol, the setClauses are used to update the columns.
func inserFromTempTableWithConflict(ctx context.Context, tx pgx.Tx, to, from string, conflictCol string, setClauses ...string) error {
	query := `
	INSERT INTO ` + to + `
	SELECT * FROM ` + from + `
	ON CONFLICT (` + conflictCol + `) DO UPDATE SET ` + strings.Join(setClauses, ", ") + `;
	`

	_, err := tx.Exec(ctx, query)

	return err
}

// insertVictim inserts a victim into the database.
func insertVictim(ctx context.Context, p DBPool, killmailID int, v esi.KillmailVictim) error {
	// TODO: this was also created before copyAny. This will always only be a single row, but is there
	// a way to optimize by using copyAny from the caller?
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

	// if err := insertVictimItems(ctx, p, killmailID, "victim", v.Items); err != nil {
	// 	return fmt.Errorf("error inserting victim items: %w", err)
	// }

	return nil
}

// func insertVictimItems(ctx context.Context, p DBPool, parentID int, parentType string, items []esi.KillmailVictimItem) error {
// 	query := `
// 	INSERT INTO killmail.victim_item(parent_id, parent_type, flag, item_type_id, quantity_destroyed, quantity_dropped, singleton)
// 	VALUES ($1, $2, $3, $4, $5, $6, $7)
// 	RETURNING victim_item_id;
// 	`

// 	for _, i := range items {
// 		row := p.QueryRow(ctx, query, parentID, parentType, i.Flag, i.ItemTypeID, i.QuantityDestroyed, i.QuantityDropped, i.Singleton)

// 		var victimItemID int
// 		if err := row.Scan(&victimItemID); err != nil {
// 			return fmt.Errorf("error inserting victim item: %w", err)
// 		}

// 		if i.Items != nil {
// 			if err := insertVictimItems(ctx, p, victimItemID, "item", *i.Items); err != nil {
// 				return fmt.Errorf("error inserting victim item details: %w", err)
// 			}
// 		}
// 	}

// 	return nil
// }

// upsertEntities upserts entities into the player database, where entityTypes is alliance, character
// or corporation. On conflict, nothing is done because we are only inserting a single column here.
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

// upsertParticipants splits participants into slices of alliance, character, and corporation IDs,
// then uses upserEntities to insert them into the database.
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

// upsertKillmail upserts a killmail into the database. On conflict of KillmailID, MoonID and WarID are updated.
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

// upsertZkillInfo upserts a zkill killmail info into the database. On conflict of ESI killmail ID,
// all fields are updated.
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

// CopyESIKillmails uses Postgres CopyFrom to copy ESI killmails into the database, copying into a
// temp table, then inserting into the persistent table.
func (c *Client) CopyESIKillmails(ctx context.Context, date string, kms []esi.Killmail) error {
	var allAttackers []attackerWithKillmailID
	var allVictims []victimWithKillmailID
	for _, k := range kms {
		for _, a := range k.Attackers {
			allAttackers = append(allAttackers, attackerWithKillmailID{int(k.KillmailID), a})
		}

		allVictims = append(allVictims, victimWithKillmailID{int(k.KillmailID), k.Victim})
	}

	if err := pgx.BeginFunc(ctx, c.pool, func(tx pgx.Tx) error {
		if err := copyParticipants(ctx, tx, date, kms); err != nil {
			return fmt.Errorf("error copying participants: %w", err)
		}

		if err := copyKillmails(ctx, tx, date, kms); err != nil {
			return fmt.Errorf("error copying killmails: %w", err)
		}

		if err := copyAttackers(ctx, tx, date, allAttackers); err != nil {
			return fmt.Errorf("error copying attackers: %w", err)
		}

		if err := copyVictims(ctx, tx, date, allVictims); err != nil {
			return fmt.Errorf("error copying victims: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("transaction error: %w", err)
	}

	return nil
}

// UpsertESIKillmail upserts an ESI killmail into the database.
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

// UpsertZkillKillmail upserts a zkill killmail into the database.
func (c *Client) UpsertZkillKillmail(ctx context.Context, z zkill.RedisQPackage) error {
	if err := c.UpsertESIKillmail(ctx, z.ESIKill); err != nil {
		return fmt.Errorf("error upserting ESI killmail: %w", err)
	}

	if err := upsertZkillInfo(ctx, c.pool, int(z.ESIKill.KillmailID), z.ZKKill); err != nil {
		return fmt.Errorf("error inserting zkill info: %w", err)
	}

	return nil
}
