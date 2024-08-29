package updater

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/antihax/goesi/esi"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/null"
)

func esiFactionToDBFaction(id int32, in esi.GetUniverseFactions200Ok) db.Faction {
	var (
		corpID null.JSONNullInt32
		milID  null.JSONNullInt32
		sysID  null.JSONNullInt32
	)

	if in.CorporationId != 0 {
		corpID.Valid = true
		corpID.Int32 = in.CorporationId
	}

	if in.MilitiaCorporationId != 0 {
		milID.Valid = true
		milID.Int32 = in.MilitiaCorporationId
	}

	if in.SolarSystemId != 0 {
		sysID.Valid = true
		sysID.Int32 = in.SolarSystemId
	}

	return db.Faction{
		CorporationID:      corpID,
		Description:        in.Description,
		FactionID:          id,
		IsUnique:           in.IsUnique,
		MilitiaCorpID:      milID,
		Name:               in.Name,
		SizeFactor:         in.SizeFactor,
		SolarSystemID:      sysID,
		StationCount:       in.StationCount,
		StationSystemCount: in.StationSystemCount,
	}
}

func (u *Updater) fetchFactions(ctx context.Context) ([]esi.GetUniverseFactions200Ok, error) {
	factions, res, err := (*u.esiUni).GetUniverseFactions(ctx, nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			return nil, fmt.Errorf("esi error: %w", esiErr)
		}

		return nil, fmt.Errorf("unknown error: %w", err)
	}

	defer res.Body.Close()

	if err := checkLimits(res); err != nil {
		return nil, fmt.Errorf("esi limit met: %w", err)
	}

	return factions, nil
}

func (u *Updater) UpdateFactions(ctx context.Context) error {
	log.Printf("updating factions\n")
	factions, err := u.fetchFactions(ctx)
	if err != nil {
		return fmt.Errorf("fetch error: %w", err)
	}

	dbFactions := make([]db.Faction, 0, len(factions))
	for _, faction := range factions {
		dbFactions = append(dbFactions, esiFactionToDBFaction(faction.FactionId, faction))
	}

	if err := u.db.CopyFactions(ctx, dbFactions); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}
