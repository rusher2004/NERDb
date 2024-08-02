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

type ESIClient interface {
	GetCharactersCharacterId(ctx context.Context, characterId int32, localVarOptionals *esi.GetCharactersCharacterIdOpts) (esi.GetCharactersCharacterIdOk, *http.Response, error)
}

type ESILimitError struct {
	Remain int
	Reset  int
}

func (e ESILimitError) Error() string {
	return fmt.Sprintf("ESI error limit reached: remain %d, reset %d", e.Remain, e.Reset)
}

type ErrNoUnnamedCharacters struct{}

func (e ErrNoUnnamedCharacters) Error() string {
	return "no unnamed characters"
}

type Updater struct {
	db  db.Client
	esi *ESIClient
}

func NewUpdater(db db.Client, ec ESIClient) *Updater {
	return &Updater{
		db:  db,
		esi: &ec,
	}
}

// checkLimits pulls the values from the X-ESI-Error-Limit-Remain and X-ESI-Error-Limit-Reset
// headers in res. If the remain value is less than 10, an ESILimitError is returned.
func checkLimits(res *http.Response) error {
	remain := res.Header.Get("X-ESI-Error-Limit-Remain")
	reset := res.Header.Get("X-ESI-Error-Limit-Reset")

	remainInt, err := strconv.Atoi(remain)
	if err != nil {
		return fmt.Errorf("error converting remain to int: %w", err)
	}
	resetInt, err := strconv.Atoi(reset)
	if err != nil {
		return fmt.Errorf("error converting reset to int: %w", err)
	}

	if remainInt < 10 {
		return ESILimitError{Remain: remainInt, Reset: resetInt}
	}

	return nil
}
