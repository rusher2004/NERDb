package esi

import (
	"context"
	"net/http"

	"github.com/antihax/goesi/esi"
)

type ESICharacterClient interface {
	GetCharactersCharacterId(ctx context.Context, characterId int32, localVarOptionals *esi.GetCharactersCharacterIdOpts) (esi.GetCharactersCharacterIdOk, *http.Response, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Updater struct {
	esi *ESICharacterClient
}

func NewUpdater(cl HTTPClient, ec *ESICharacterClient) *Updater {
	return &Updater{
		esi: ec,
	}
}
