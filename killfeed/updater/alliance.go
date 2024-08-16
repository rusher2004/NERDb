package updater

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/antihax/goesi/esi"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/null"
)

// esiAllianceToDBAlliance will convert a goesi alliance response to a database Alliance struct.
func esiAllianceToDBAlliance(id int32, in esi.GetAlliancesAllianceIdOk) db.Alliance {
	var (
		outExecCorpID null.JSONNullInt32
		outFacID      null.JSONNullInt32
	)

	if in.ExecutorCorporationId != 0 {
		outExecCorpID.Valid = true
		outExecCorpID.Int32 = in.ExecutorCorporationId
	}

	if in.FactionId != 0 {
		outFacID.Valid = true
		outFacID.Int32 = in.FactionId
	}

	return db.Alliance{
		AllianceID:    id,
		CreatorCorpID: in.CreatorCorporationId,
		CreatorID:     in.CreatorId,
		DateFounded:   in.DateFounded,
		ExecCorpID:    outExecCorpID,
		FactionID:     outFacID,
		Name:          in.Name,
		Ticker:        in.Ticker,
	}
}

// fetchESIAllianceIDs will fetch all the alliance IDs from ESI.
func (u *Updater) fetchESIAllianceIDs(ctx context.Context) ([]int32, error) {
	ids, res, err := (*u.esiAlly).GetAlliances(ctx, nil)
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

	return ids, nil
}

// fetchAllianceInfo will fetch the alliance information from ESI for the given ID.
func (u *Updater) fetchAllianceInfo(ctx context.Context, id int32) (db.Alliance, error) {
	ally, res, err := (*u.esiAlly).GetAlliancesAllianceId(ctx, id, nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			switch t := esiErr.Model().(type) {
			case esi.GetAlliancesAllianceIdNotFound:
				return db.Alliance{}, fmt.Errorf("alliance %d not found: %s", id, t.Error_)
			}
		}

		return db.Alliance{}, fmt.Errorf("error getting alliance %d: %w", id, err)
	}
	defer res.Body.Close()

	if err := checkLimits(res); err != nil {
		return db.Alliance{}, fmt.Errorf("esi limit met: %w", err)
	}

	return esiAllianceToDBAlliance(id, ally), nil
}

// fetchAlliancesInfo will fetch the alliance information from ESI for the given IDs.
func (u *Updater) fetchAlliancesInfo(ctx context.Context, ids []int32) ([]db.Alliance, error) {
	out := make([]db.Alliance, 0, len(ids))

	for i, id := range ids {
		ally, err := u.fetchAllianceInfo(ctx, id)
		if err != nil {
			return nil, fmt.Errorf("error getting alliance %d: %w", id, err)
		}

		log.Printf("got alliance %d/%d: %s\n", i+1, len(ids), ally.Name)

		out = append(out, ally)
	}

	return out, nil
}

// updateAlliance will fetch the alliance information from ESI and update the database with the result.
func (u *Updater) updateAlliance(ctx context.Context, id int32) error {
	ally, res, err := (*u.esiAlly).GetAlliancesAllianceId(ctx, id, nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			switch t := esiErr.Model().(type) {
			case esi.GetAlliancesAllianceIdNotFound:
				log.Printf("alliance %d not found: %s\n", id, t.Error_)
			}

			log.Printf("error getting alliance %d: %v\n", id, err)
			return nil
		}

		return fmt.Errorf("error getting alliance %d: %w", id, err)
	}
	defer res.Body.Close()

	log.Println("updating", ally.Name)

	if err := checkLimits(res); err != nil {
		return fmt.Errorf("esi limit met: %w", err)
	}

	dbAlly := esiAllianceToDBAlliance(id, ally)
	if err := u.db.UpdateAlliance(ctx, dbAlly); err != nil {
		return fmt.Errorf("error updating alliance %d: %w", id, err)
	}

	return nil
}

// updateMissingAlliances updates the database with alliances that are missing.
func (u *Updater) updateMissingAlliances(ctx context.Context) error {
	ids, err := u.fetchESIAllianceIDs(ctx)
	if err != nil {
		var limitErr ESILimitError
		if errors.As(err, &limitErr) {
			log.Printf("ESI limit reached: %v\n", limitErr)
			time.Sleep(time.Duration(limitErr.Reset+1) * time.Second)
		}

		return fmt.Errorf("error getting alliance IDs: %w", err)
	}

	missing, err := u.db.GetMissingAllianceIDs(ctx, ids)
	if err != nil {
		return fmt.Errorf("error getting missing alliance IDs: %w", err)
	}

	log.Printf("found %d missing alliances\n", len(missing))

	missingInfo, err := u.fetchAlliancesInfo(ctx, missing)
	if err != nil {
		return fmt.Errorf("error getting alliances info: %w", err)
	}

	if err := u.db.CopyAlliances(ctx, missingInfo); err != nil {
		return fmt.Errorf("error inserting alliances: %w", err)
	}

	return nil
}

// updateUnnamedAlliances updates the database with alliances that have no name.
func (u *Updater) updateUnnamedAlliances(ctx context.Context, count int) error {
	unnamed, err := u.db.GetUnnamedAllianceIDs(ctx, count)
	if err != nil {
		return fmt.Errorf("error getting unnamed alliance IDs: %w", err)
	}

	if len(unnamed) == 0 {
		return ErrNoUnnamedAlliances{}
	}

	log.Printf("found %d unnamed alliances\n", len(unnamed))

	unnamedInfo, err := u.fetchAlliancesInfo(ctx, unnamed)
	if err != nil {
		return fmt.Errorf("error getting unnamed alliances info: %w", err)
	}

	if err := u.db.CopyAlliances(ctx, unnamedInfo); err != nil {
		return fmt.Errorf("error inserting unnamed alliances: %w", err)
	}

	return nil
}

// UpdateAlliances updates the database with missing alliances. It will both check for missing
// alliances against the set provided by ESI, and alliances in the db with no name.
func (u *Updater) UpdateAlliances(ctx context.Context, count int) error {
	if err := u.updateMissingAlliances(ctx); err != nil {
		return fmt.Errorf("error updating missing alliances: %w", err)
	}

	for {
		if err := u.updateUnnamedAlliances(ctx, count); err != nil {
			return fmt.Errorf("error updating unnamed alliances: %w", err)
		}
	}
}
