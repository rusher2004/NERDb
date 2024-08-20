package updater

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/antihax/goesi/esi"
	"github.com/jackc/pgx/v5"
	"github.com/rusher2004/nerdb/killfeed/db"
)

type DBClient interface {
	Exec(ctx context.Context, sql string, args ...interface{}) (int64, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

type ESIAllianceClient interface {
	GetAlliances(ctx context.Context, localVarOptionals *esi.GetAlliancesOpts) ([]int32, *http.Response, error)
	GetAlliancesAllianceId(ctx context.Context, allianceId int32, localVarOptionals *esi.GetAlliancesAllianceIdOpts) (esi.GetAlliancesAllianceIdOk, *http.Response, error)
}

type ESICharacterClient interface {
	GetCharactersCharacterId(ctx context.Context, characterId int32, localVarOptionals *esi.GetCharactersCharacterIdOpts) (esi.GetCharactersCharacterIdOk, *http.Response, error)
}

type ESICorporationClient interface {
	GetCorporationsCorporationId(ctx context.Context, corporationId int32, localVarOptionals *esi.GetCorporationsCorporationIdOpts) (esi.GetCorporationsCorporationIdOk, *http.Response, error)
}

type ESILimitError struct {
	Remain int
	Reset  int
}

func (e ESILimitError) Error() string {
	return fmt.Sprintf("ESI error limit reached: remain %d, reset %d", e.Remain, e.Reset)
}

type ESIUnknownError struct {
	Err    error
	Header http.Header
}

func (e ESIUnknownError) Error() string {
	headersStr := ""
	for k, v := range e.Header {
		headersStr += fmt.Sprintf("%s: %s\n", k, v)
	}

	return fmt.Sprintf("unknown ESI error: %s; headers: %s", e.Err, headersStr)
}

type ErrNoUnnamedAlliances struct{}

func (e ErrNoUnnamedAlliances) Error() string {
	return "no unnamed alliances"
}

type ErrNoUnnamedCharacters struct{}

func (e ErrNoUnnamedCharacters) Error() string {
	return "no unnamed characters"
}

type ErrNoUnnamedCorporations struct{}

func (e ErrNoUnnamedCorporations) Error() string {
	return "no unnamed corporations"
}

type Updater struct {
	db      db.Client
	esiAlly *ESIAllianceClient
	esiChar *ESICharacterClient
	esiCorp *ESICorporationClient
}

func NewUpdater(db db.Client, ally ESIAllianceClient, char ESICharacterClient, corp ESICorporationClient) *Updater {
	return &Updater{
		db:      db,
		esiAlly: &ally,
		esiChar: &char,
		esiCorp: &corp,
	}
}

func (u *Updater) Update(ctx context.Context, kind string, limit int) error {
	switch kind {
	case "alliance":
		return u.UpdateAlliances(ctx, limit)
	case "character":
		return u.UpdateCharacters(ctx, limit)
	case "corporation":
		return u.UpdateCorporations(ctx, limit)
	default:
		return fmt.Errorf("unknown type: %s", kind)
	}
}

// checkLimits pulls the values from the X-ESI-Error-Limit-Remain and X-ESI-Error-Limit-Reset
// headers in res. If the remain value is less than 10, an ESILimitError is returned.
func checkLimits(res *http.Response) error {
	remain := res.Header.Get("X-ESI-Error-Limit-Remain")
	reset := res.Header.Get("X-ESI-Error-Limit-Reset")

	remainInt, err := strconv.Atoi(remain)
	if err != nil {
		return ESIUnknownError{Err: fmt.Errorf("remain parse error: %w", err), Header: res.Header.Clone()}
	}
	resetInt, err := strconv.Atoi(reset)
	if err != nil {
		return ESIUnknownError{Err: fmt.Errorf("reset parse error: %w", err), Header: res.Header.Clone()}
	}

	if remainInt < 10 {
		return ESILimitError{Remain: remainInt, Reset: resetInt}
	}

	return nil
}
