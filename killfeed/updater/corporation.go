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

func esiCorpToDBCorp(id int32, in esi.GetCorporationsCorporationIdOk) db.Corporation {
	var (
		outAllID   null.JSONNullInt32
		outCEO     null.JSONNullInt32
		outCreator null.JSONNullInt32
		outFac     null.JSONNullInt32
		outHome    null.JSONNullInt32
		outURL     null.JSONNullString
	)

	if in.AllianceId != 0 {
		outAllID.Valid = true
		outAllID.Int32 = in.AllianceId
	}

	if in.CeoId != 0 {
		outCEO.Valid = true
		outCEO.Int32 = in.CeoId
	}

	if in.CreatorId != 0 {
		outCreator.Valid = true
		outCreator.Int32 = in.CreatorId
	}

	if in.FactionId != 0 {
		outFac.Valid = true
		outFac.Int32 = in.FactionId
	}

	if in.HomeStationId != 0 {
		outHome.Valid = true
		outHome.Int32 = in.HomeStationId
	}

	if in.Url != "" {
		outURL.Valid = true
		outURL.String = in.Url
	}

	return db.Corporation{
		AllianceID:    outAllID,
		CEOID:         outCEO,
		CorporationID: id,
		CreatorID:     outCreator,
		Description:   in.Description,
		FactionID:     outFac,
		HomeStationID: outHome,
		MemberCount:   in.MemberCount,
		Name:          in.Name,
		Shares:        in.Shares,
		TaxRate:       in.TaxRate,
		Ticker:        in.Ticker,
		URL:           outURL,
		WarEligible:   in.WarEligible,
	}
}

func (u *Updater) updateCorporation(ctx context.Context, corpID int32) error {
	corp, res, err := (*u.esiCorp).GetCorporationsCorporationId(ctx, corpID, nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			switch t := esiErr.Model().(type) {
			case esi.GetCorporationsCorporationIdNotFound:
				log.Printf("corporation %d not found: %s\n", corpID, t.Error_)
			}

			log.Printf("error getting corporation %d: %v\n", corpID, err)
			return nil
		}

		return fmt.Errorf("error getting corporation %d: %w", corpID, err)
	}
	defer res.Body.Close()

	log.Println("updating", corp.Name)

	if err := checkLimits(res); err != nil {
		return fmt.Errorf("esi limit met: %w", err)
	}

	dbCorp := esiCorpToDBCorp(corpID, corp)
	if err := u.db.UpdateCorporation(ctx, dbCorp); err != nil {
		return fmt.Errorf("error updating corporation %d: %w", corpID, err)
	}

	return nil
}

func (u *Updater) UpdateCorporations(ctx context.Context, count int) error {
	corpIDs, err := u.db.GetUnnamedCorporationIDs(ctx, count)
	if err != nil {
		return fmt.Errorf("error getting corporation IDs: %w", err)
	}

	if len(corpIDs) == 0 {
		return ErrNoUnnamedCorporations{}
	}

	for i, id := range corpIDs {
		log.Printf("fetching corporation %d (%d/%d)\n", id, i+1, len(corpIDs))
		if err := u.updateCorporation(ctx, id); err != nil {
			var limitErr ESILimitError
			if errors.As(err, &limitErr) {
				log.Printf("ESI limit reached: %v\n", limitErr)
				time.Sleep(time.Duration(limitErr.Reset+1) * time.Second)
			}
			return fmt.Errorf("error updating corporation %d: %w", id, err)
		}
	}

	return nil
}
