package everef

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/antihax/goesi"
	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/null"
)

type Alliance struct {
	ID                      OID          `json:"_id"`
	AllianceID              int          `json:"alliance_id"`
	CreatorCorporationID    int          `json:"creator_corporation_id"`
	CreatorID               int          `json:"creator_id"`
	DateFounded             time.Time    `json:"date_founded"`
	ExecutorCorporationID   int          `json:"executor_corporation_id"`
	FactionID               int          `json:"faction_id"`
	LastModified            LastModified `json:"last_modified"`
	Name                    string       `json:"name"`
	Ticker                  string       `json:"ticker"`
	Losses                  int          `json:"losses"`
	CreatorCorporationName  string       `json:"creator_corporation_name"`
	CreatorName             string       `json:"creator_name"`
	ExecutorCorporationName string       `json:"executor_corporation_name"`
	FactionName             string       `json:"faction_name"`
	Points                  int          `json:"points"`
	Kills                   int          `json:"kills"`
}

type Character struct {
	ID              OID          `json:"_id"`
	CharacterID     int          `json:"character_id"`
	AllianceID      int          `json:"alliance_id"`
	Birthday        time.Time    `json:"birthday"`
	BloodlineID     int          `json:"bloodline_id"`
	CorporationID   int          `json:"corporation_id"`
	Description     string       `json:"description"`
	Gender          string       `json:"gender"`
	LastModified    LastModified `json:"last_modified"`
	Name            string       `json:"name"`
	RaceID          int          `json:"race_id"`
	SecurityStatus  float64      `json:"security_status"`
	Losses          int          `json:"losses"`
	AllianceName    string       `json:"alliance_name"`
	CorporationName string       `json:"corporation_name"`
	FactionName     string       `json:"faction_name"`
	Points          float64      `json:"points"`
	Kills           int          `json:"kills"`
}

type Corporation struct {
	ID              OID          `json:"_id"`
	CorporationID   int          `json:"corporation_id"`
	CeoID           int          `json:"ceo_id"`
	CreatorID       int          `json:"creator_id"`
	DateFounded     time.Time    `json:"date_founded"`
	Description     string       `json:"description"`
	FactionID       int          `json:"faction_id"`
	HomeStationID   int          `json:"home_station_id"`
	LastModified    LastModified `json:"last_modified"`
	MemberCount     int          `json:"member_count"`
	Name            string       `json:"name"`
	Shares          int          `json:"shares"`
	TaxRate         int          `json:"tax_rate"`
	Ticker          string       `json:"ticker"`
	URL             string       `json:"url"`
	Losses          int          `json:"losses"`
	AllianceName    string       `json:"alliance_name"`
	CeoName         string       `json:"ceo_name"`
	CreatorName     string       `json:"creator_name"`
	FactionName     string       `json:"faction_name"`
	HomeStationName string       `json:"home_station_name"`
	Points          float64      `json:"points"`
	Kills           int          `json:"kills"`
}

type LastModified struct {
	Date time.Time `json:"$date"`
}

type OID struct {
	OID string `json:"$oid"`
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// toDBCharacter returns a db.Character from a Character, with nullable fields based off of zero
// values. Title is not set, as it is not available from the EveRef dataset.
func (c Character) toDBCharacter() db.Character {
	var (
		outAllID null.JSONNullInt32
		outDesc  null.JSONNullString
		outFac   null.JSONNullInt32
		outTitle null.JSONNullString
	)

	if c.AllianceID != 0 {
		outAllID.Valid = true
		outAllID.Int32 = int32(c.AllianceID)
	}

	if c.Description != "" {
		outDesc.Valid = true
		outDesc.String = c.Description
	}

	if c.FactionName != "" {
		outFac.Valid = true
		outFac.Int32 = goesi.FactionNameToID(c.FactionName)
	}

	// title needs to be done on demand
	// if c.Title != "" {
	// 	outTitle.Valid = true
	// 	outTitle.String = c.Title
	// }

	return db.Character{
		AllianceID:     outAllID,
		Birthday:       c.Birthday,
		BloodlineID:    int32(c.BloodlineID),
		CharacterID:    int32(c.CharacterID),
		CorporationID:  int32(c.CorporationID),
		Description:    outDesc,
		Gender:         db.Gender(c.Gender),
		FactionID:      outFac,
		Name:           c.Name,
		RaceID:         int32(c.RaceID),
		SecurityStatus: float32(c.SecurityStatus),
		Title:          outTitle,
	}
}

// FetchTotals returns a map of unique days to the number of killmails on that day, sourced from
// the EveRef API.
func FetchTotals(cl httpClient) (map[string]int, error) {
	req, err := http.NewRequest("GET", "https://data.everef.net/killmails/totals.json", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	res, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("enexpected status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %w", err)
	}

	var totals map[string]int
	if err := json.Unmarshal(b, &totals); err != nil {
		return nil, fmt.Errorf("error unmarshalling body: %w", err)
	}

	return totals, nil
}
