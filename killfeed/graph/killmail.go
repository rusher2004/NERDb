package graph

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rusher2004/nerdb/killfeed/esi"
	"github.com/rusher2004/nerdb/killfeed/zkill"
)

func insertAttackers(ctx context.Context, tx neo4j.ManagedTransaction, id int, attackers []esi.KillmailAttacker) error {
	var attackerParams []map[string]any
	for _, a := range attackers {
		m := map[string]any{
			"damageDone":     a.DamageDone,
			"finalBlow":      a.FinalBlow,
			"securityStatus": a.SecurityStatus,
		}

		// we can't use null values in the query, it will throw an error
		if a.CharacterID.Valid && a.CharacterID.Int32 != 0 {
			m["characterId"] = a.CharacterID.Int32
		}
		if a.CorporationID.Valid && a.CorporationID.Int32 != 0 {
			m["corporationId"] = a.CorporationID.Int32
		}
		if a.AllianceID.Valid && a.AllianceID.Int32 != 0 {
			m["allianceId"] = a.AllianceID.Int32
		}
		if a.FactionID.Valid && a.FactionID.Int32 != 0 {
			m["factionId"] = a.FactionID.Int32
		}
		if a.ShipTypeID.Valid && a.ShipTypeID.Int32 != 0 {
			m["shipTypeId"] = a.ShipTypeID.Int32
		}
		if a.WeaponTypeID.Valid && a.WeaponTypeID.Int32 != 0 {
			m["weaponTypeId"] = a.WeaponTypeID.Int32
		}

		attackerParams = append(attackerParams, m)
	}

	query := `
		MATCH (k:Killmail {killmailId: $killmailId})
		WITH k
		CALL apoc.create.nodes($labels, $attackerParams) YIELD node
		CREATE (node)-[:ATTACKER_OF]->(k)
		RETURN count(node) as count;
	`

	params := map[string]any{
		"killmailId":     id,
		"labels":         []string{"Attacker"},
		"attackerParams": attackerParams,
	}

	res, err := tx.Run(ctx, query, params)
	if err != nil {
		return fmt.Errorf("error executing attacker query: %w", err)
	}

	record, err := res.Single(ctx)
	if err != nil {
		return fmt.Errorf("error getting attacker count: %w", err)
	}

	count, ok := record.Get("count")
	if !ok {
		return fmt.Errorf("attacker count not found in result")
	}

	countInt, ok := count.(int64)
	if !ok {
		return fmt.Errorf("attacker count is not an int")
	}

	if int(countInt) != len(attackers) {
		return fmt.Errorf("attacker count mismatch: %d != %d", countInt, len(attackers))
	}

	return nil
}

func insertKillmailWithVictim(ctx context.Context, tx neo4j.ManagedTransaction, k esi.Killmail) (int, error) {
	query := `
				MERGE (k:Killmail {killmailId: $killmailId})
				ON CREATE SET
					k += $kmData
				MERGE (v:Victim)-[:VICTIM_OF]->(k)
				ON CREATE SET
					v += $charData
				RETURN k.killmailId as id;
			`

	kmData := map[string]any{
		"time":          k.KillmailTime,
		"solarSystemId": k.SolarSystemID,
		"created":       time.Now().UTC(),
	}

	if k.MoonID.Valid && k.MoonID.Int32 != 0 {
		kmData["moonId"] = k.MoonID.Int32
	}

	if k.WarID.Valid && k.WarID.Int32 != 0 {
		kmData["warId"] = k.WarID.Int32
	}

	charData := map[string]any{
		"damageTaken": k.Victim.DamageTaken,
		"shipTypeId":  k.Victim.ShipTypeID,
		"created":     time.Now().UTC(),
	}

	if k.Victim.CharacterID.Valid && k.Victim.CharacterID.Int32 != 0 {
		charData["characterId"] = k.Victim.CharacterID.Int32
	}

	if k.Victim.CorporationID.Valid && k.Victim.CorporationID.Int32 != 0 {
		charData["corporationId"] = k.Victim.CorporationID.Int32
	}

	if k.Victim.AllianceID.Valid && k.Victim.AllianceID.Int32 != 0 {
		charData["allianceId"] = k.Victim.AllianceID.Int32
	}

	if k.Victim.FactionID.Valid && k.Victim.FactionID.Int32 != 0 {
		charData["factionId"] = k.Victim.FactionID.Int32
	}

	if k.Victim.Position != nil {
		charData["position"] = neo4j.Point3D{
			X:            k.Victim.Position.X,
			Y:            k.Victim.Position.Y,
			Z:            k.Victim.Position.Z,
			SpatialRefId: 9157, // 3D Cartesian
		}
	}

	params := map[string]any{
		"killmailId": k.KillmailID,
		"kmData":     kmData,
		"charData":   charData,
	}

	res, err := tx.Run(ctx, query, params)
	if err != nil {
		return 0, fmt.Errorf("error executing killmail query: %w", err)
	}

	record, err := res.Single(ctx)
	if err != nil {
		return 0, fmt.Errorf("error getting killmail id: %w", err)
	}

	id, ok := record.Get("id")
	if !ok {
		return 0, fmt.Errorf("killmail id not found in result")
	}

	setID, ok := id.(int64)
	if !ok {
		return 0, fmt.Errorf("killmail id is not an int")
	}

	return int(setID), nil
}

func (c Client) copyESIKillmails(ctx context.Context, kms []esi.Killmail) error {
	sess := c.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sess.Close(ctx)

	for _, k := range kms {
		if _, err := sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
			setID, err := insertKillmailWithVictim(ctx, tx, k)
			if err != nil {
				return nil, fmt.Errorf("error inserting killmail with victim: %w", err)
			}

			if err := insertAttackers(ctx, tx, setID, k.Attackers); err != nil {
				return nil, fmt.Errorf("error inserting attackers: %w", err)
			}

			return nil, nil
		}); err != nil {
			return fmt.Errorf("error executing write: %w", err)
		}
	}

	return nil
}

func (c Client) insertZkillInfos(ctx context.Context, z []zkill.RedisQPackage) error {
	sess := c.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sess.Close(ctx)

	var zkillParams []map[string]any
	for _, p := range z {
		m := map[string]any{
			"awox":           p.ZKKill.Awox,
			"destroyedValue": p.ZKKill.DestroyedValue,
			"droppedValue":   p.ZKKill.DroppedValue,
			"fittedValue":    p.ZKKill.FittedValue,
			"hash":           p.ZKKill.Hash,
			"npc":            p.ZKKill.NPC,
			"points":         p.ZKKill.Points,
			"solo":           p.ZKKill.Solo,
			"totalValue":     p.ZKKill.TotalValue,
			"killmailId":     p.KillID,
		}

		if p.ZKKill.LocationID.Valid && p.ZKKill.LocationID.Int64 != 0 {
			m["locationId"] = p.ZKKill.LocationID.Int64
		}

		zkillParams = append(zkillParams, m)
	}

	query := `
		CALL apoc.create.nodes($labels, $zkillParams) YIELD node
		WITH collect(node) as nodes
		UNWIND nodes as n
		MATCH (k:Killmail {killmailId: n.killmailId})
		MERGE (n)-[:REFERENCES]->(k)
		REMOVE n.killmailId
		RETURN count(n) as count;
	`

	params := map[string]any{
		"labels":      []string{"ZKILL_INFO"},
		"zkillParams": zkillParams,
	}

	count, err := sess.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		res, err := tx.Run(ctx, query, params)
		if err != nil {
			return nil, fmt.Errorf("error executing zkill query: %w", err)
		}

		record, err := res.Single(ctx)
		if err != nil {
			return nil, fmt.Errorf("error getting zkill count: %w", err)
		}

		count, ok := record.Get("count")
		if !ok {
			return nil, fmt.Errorf("zkill count not found in result")
		}

		countInt, ok := count.(int64)
		if !ok {
			return nil, fmt.Errorf("zkill count is not an int")
		}

		if int(countInt) != len(z) {
			return nil, fmt.Errorf("zkill count mismatch: %d != %d", count.(int64), len(z))
		}

		return countInt, nil
	})
	if err != nil {
		return fmt.Errorf("error executing write: %w", err)
	}

	// log.Printf("inserted %d zkill infos", count.(int64))
	log.Printf("inserted %v zkill infos\n", count)

	return nil
}

func (c Client) CopyESIKillmails(ctx context.Context, date string, kms []esi.Killmail) error {
	return c.copyESIKillmails(ctx, kms)
}

func (c Client) CopyZkillKillmails(ctx context.Context, _ string, z []zkill.RedisQPackage) error {
	esiKkills := make([]esi.Killmail, len(z))
	for i, p := range z {
		esiKkills[i] = p.ESIKill
	}

	if err := c.copyESIKillmails(ctx, esiKkills); err != nil {
		return fmt.Errorf("error copying zkill killmails: %w", err)
	}

	if err := c.insertZkillInfos(ctx, z); err != nil {
		return fmt.Errorf("error inserting zkill infos: %w", err)
	}

	return nil
}
