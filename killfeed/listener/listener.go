package listener

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/rusher2004/nerdb/killfeed/db"
	"github.com/rusher2004/nerdb/killfeed/zkill"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Listener struct {
	hc HTTPClient
	db db.Client
}

func NewListener(hc HTTPClient, db db.Client) *Listener {
	return &Listener{hc, db}
}

func (l Listener) Listen(ctx context.Context) {
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
			l.doRequest(ctx, req)
		}
	}
}

func (l Listener) doRequest(ctx context.Context, req *http.Request) {
	res, err := l.hc.Do(req)
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

	var z zkill.RedisQMessage
	if err := json.Unmarshal(b, &z); err != nil {
		log.Printf("error unmarshalling body: %v", err)
		return
	}

	if z.Package.KillID == 0 {
		log.Printf("no killmails available")
		return
	}

	log.Printf("dispatching to db: %d", z.Package.KillID)
	go dispatchKillmail(ctx, l.db, z.Package)
}

func dispatchKillmail(ctx context.Context, cl db.Client, z zkill.RedisQPackage) {
	if err := cl.UpsertZkillKillmail(ctx, z); err != nil {
		log.Printf("error upserting killmail: %v", err)
	}
}
