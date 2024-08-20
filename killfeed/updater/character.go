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

// esiCharToDBChar converts in to a db.Character struct.
func esiCharToDBChar(id int32, in esi.GetCharactersCharacterIdOk) db.Character {
	var (
		outAllID null.JSONNullInt32
		outDesc  null.JSONNullString
		outFac   null.JSONNullInt32
		outTitle null.JSONNullString
	)

	if in.AllianceId != 0 {
		outAllID.Valid = true
		outAllID.Int32 = in.AllianceId
	}

	if in.Description != "" {
		outDesc.Valid = true
		outDesc.String = in.Description
	}

	if in.FactionId != 0 {
		outFac.Valid = true
		outFac.Int32 = in.FactionId
	}

	if in.Title != "" {
		outTitle.Valid = true
		outTitle.String = in.Title
	}

	return db.Character{
		AllianceID:     outAllID,
		Birthday:       in.Birthday,
		BloodlineID:    in.BloodlineId,
		CharacterID:    id,
		CorporationID:  in.CorporationId,
		Description:    outDesc,
		Gender:         db.Gender(in.Gender),
		FactionID:      outFac,
		Name:           in.Name,
		RaceID:         in.RaceId,
		SecurityStatus: in.SecurityStatus,
		Title:          outTitle,
	}
}

func fetchCharacter(ctx context.Context, cc *ESICharacterClient, charID int32) (db.Character, bool, error) {
	ok, res, err := (*cc).GetCharactersCharacterId(ctx, int32(charID), nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			switch t := esiErr.Model().(type) {
			case esi.GetCharactersCharacterIdNotFound:
				log.Printf("character %d not found: %s\n", charID, t.Error_)
				return db.Character{}, true, nil
			case esi.GatewayTimeout:
				log.Printf("gateway timeout (%d seconds): %s\n", t.Timeout, t.Error_)
				return db.Character{}, false, ESILimitError{Remain: 0, Reset: int(t.Timeout)}
			}

			if err := checkLimits(res); err != nil {
				return db.Character{}, false, fmt.Errorf("esi limit met: %w", err)
			}

			return db.Character{}, false, fmt.Errorf("unkown ESI error: %w", err)
		}

		return db.Character{}, false, ESIUnknownError{Err: err, Header: res.Header.Clone()}
	}
	defer res.Body.Close()

	return esiCharToDBChar(charID, ok), false, nil
}

// updateCharacter fetches character info from the ESI API and then updates the corresponding row in the database.
func (u *Updater) updateCharacter(ctx context.Context, charID int32) error {
	ok, res, err := (*u.esiChar).GetCharactersCharacterId(ctx, int32(charID), nil)
	if err != nil {
		var esiErr esi.GenericSwaggerError
		if errors.As(err, &esiErr) {
			switch t := esiErr.Model().(type) {
			case esi.GetCharactersCharacterIdNotFound:
				log.Printf("character %d not found: %s\n", charID, t.Error_)
				if err := u.db.SetCharacterDeleted(ctx, charID); err != nil {
					log.Printf("error setting character %d as deleted: %v", charID, err)
				}
			}

			log.Printf("error getting character %d: %v\n", charID, err)
			if err := checkLimits(res); err != nil {
				return fmt.Errorf("esi limit met: %w", err)
			}

			return nil
		}

		return fmt.Errorf("error getting character %d: %w", charID, err)
	}
	defer res.Body.Close()

	log.Println("updating", ok.Name)

	if err := checkLimits(res); err != nil {
		return fmt.Errorf("esi limit met: %w", err)
	}

	dbChar := esiCharToDBChar(charID, ok)
	if err := u.db.UpdateCharacter(ctx, dbChar); err != nil {
		return fmt.Errorf("error updating character %d: %w", charID, err)
	}

	return nil
}

// UpdateCharacters updates count number of characters in the database whose row contains null
// values for name and esi_deleted. For each row, a call to the ESI API is made for the data,
// pausing as needed, based on ESI's limit values.
func (u *Updater) UpdateCharacters(ctx context.Context, count int) error {
	charIDs, err := u.db.GetUnnamedCharacterIDs(ctx, count)
	if err != nil {
		return fmt.Errorf("error getting character IDs: %w", err)
	}

	if len(charIDs) == 0 {
		return ErrNoUnnamedCharacters{}
	}

	chars := make([]db.Character, 0)
	deletedChars := make([]int32, 0)
	for i, id := range charIDs {
		log.Printf("fetching character %d (%d/%d)\n", id, i+1, len(charIDs))
		char, deleted, err := fetchCharacter(ctx, u.esiChar, id)
		if err != nil {
			var limitErr ESILimitError
			if errors.As(err, &limitErr) {
				log.Printf("ESI limit reached: %v\n", limitErr)
				time.Sleep(time.Duration(limitErr.Reset+1) * time.Second)
				continue
			}

			return fmt.Errorf("error fetching character %d: %w", id, err)
		}

		if deleted {
			deletedChars = append(deletedChars, id)
			continue
		}

		chars = append(chars, char)
	}

	if err := u.db.CopyCharacters(ctx, chars); err != nil {
		return fmt.Errorf("error copying characters: %w", err)
	}

	if err := u.db.CopyDeletedCharacters(ctx, deletedChars); err != nil {
		return fmt.Errorf("error copying deleted characters: %w", err)
	}

	return nil
}
