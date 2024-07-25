package listener

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/rusher2004/nerdb/esi"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type zkillResponse struct {
	Package zkillPackage `json:"package"`
}

type zkillPackage struct {
	ESIKill esi.Killmail `json:"killmail"`
	KillID  int          `json:"killID"`
	ZKKill  zkKill       `json:"zkb"`
}

type zkKill struct {
	LocationID     int64    `json:"locationID,omitempty"`
	Hash           string   `json:"hash"`
	FittedValue    float64  `json:"fittedValue"`
	TotalValue     float64  `json:"totalValue"`
	Points         int64    `json:"points"`
	DestroyedValue float64  `json:"destroyedValue"`
	DroppedValue   float64  `json:"droppedValue"`
	Npc            bool     `json:"npc"`
	Solo           bool     `json:"solo"`
	Awox           bool     `json:"awox"`
	Labels         []string `json:"labels"`
}

func Listen(ctx context.Context, hc HTTPClient) {
	req, err := http.NewRequest("GET", "https://redisq.zkillboard.com/listen.php?queueID=nerdb", nil)
	if err != nil {
		log.Printf("error creating listen request: %v", err)
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			// listen for killmails
			doRequest(hc, req)
		}
	}
}

func doRequest(hc HTTPClient, req *http.Request) {
	res, err := hc.Do(req)
	if err != nil {
		log.Printf("error making request: %v", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", res.StatusCode)
		return
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
		return
	}

	var z zkillResponse
	if err := json.Unmarshal(b, &z); err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return
	}

	if z.Package.KillID == 0 {
		log.Printf("no killmails available")
		return
	}

	log.Printf("got kill: %d\n%+v", z.Package.KillID, z.Package.ESIKill)
	log.Printf("victim: %+v", z.Package.ESIKill.Victim.CharacterID)
	// log.Printf("attackers: %+v", z.Package.ESIKill.Attackers)
}
