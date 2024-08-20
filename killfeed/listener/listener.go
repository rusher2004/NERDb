package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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

// Listen uses Zklilboard's RedisQ to listen for new killmails and dispatches them to the database.
func (l Listener) Listen(ctx context.Context, batchSize int) {
	// https://github.com/zKillboard/RedisQ
	req, err := http.NewRequest("GET", "https://redisq.zkillboard.com/listen.php?queueID=nerdb", nil)
	if err != nil {
		log.Printf("error creating listen request: %v", err)
		return
	}

	packages := make([]zkill.RedisQPackage, 0, batchSize)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			p, err := l.doRequest(req)
			if err != nil {
				log.Println(err)
				continue
			}

			packages = append(packages, p)
			log.Printf("got %d/%d: %d", len(packages), batchSize, p.KillID)

			if len(packages) == batchSize {
				log.Printf("copying %d killmails", batchSize)
				now := time.Now()

				if err := l.db.CopyZkillKillmails(ctx, now.Format("20060102030405"), packages); err != nil {
					log.Printf("error copying killmails: %v", err)
				}

				packages = make([]zkill.RedisQPackage, 0, batchSize)
			}
		}
	}
}

// doRequest makes a request to the RedisQ and dispatches the killmail to the database.
func (l Listener) doRequest(req *http.Request) (zkill.RedisQPackage, error) {
	res, err := l.hc.Do(req)
	if err != nil {

		return zkill.RedisQPackage{}, fmt.Errorf("error making request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {

		return zkill.RedisQPackage{}, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {

		return zkill.RedisQPackage{}, fmt.Errorf("error reading body: %w", err)
	}

	var z zkill.RedisQMessage
	if err := json.Unmarshal(b, &z); err != nil {

		return zkill.RedisQPackage{}, fmt.Errorf("error unmarshalling body: %w", err)
	}

	if z.Package.KillID == 0 {

		return zkill.RedisQPackage{}, fmt.Errorf("no killmails available")
	}

	return z.Package, nil
	// log.Printf("dispatching to db: %d", z.Package.KillID)
	// go dispatchKillmail(ctx, l.db, z.Package)
}

// dispatchKillmail upserts the killmail into the database.
func dispatchKillmail(ctx context.Context, cl db.Client, z zkill.RedisQPackage) {
	if err := cl.UpsertZkillKillmail(ctx, z); err != nil {
		log.Printf("error upserting killmail: %v", err)
	}
}
