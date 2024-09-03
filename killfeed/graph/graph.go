package graph

import (
	"context"
	"fmt"
	"strings"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphDriver interface {
	Close(ctx context.Context) error
	NewSession(ctx context.Context, config neo4j.SessionConfig) neo4j.SessionWithContext
}

type Client struct {
	driver neo4j.DriverWithContext
}

func NewClient(uri, user, pass string) (Client, error) {
	driver, err := neo4j.NewDriverWithContext(uri, neo4j.BasicAuth(user, pass, ""))
	if err != nil {
		return Client{}, fmt.Errorf("error creating neo4j driver: %w", err)
	}

	return Client{
		driver: driver,
	}, nil
}

func (c Client) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

func (c Client) FetchTotals(ctx context.Context) (map[string]int, error) {
	sess := c.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer sess.Close(ctx)

	totals, err := neo4j.ExecuteRead(ctx, sess, func(tx neo4j.ManagedTransaction) (map[string]int, error) {
		query := `MATCH (k:Killmail)
		WITH {day: date(k.time)} as days
		RETURN days.day as day, COUNT(days.day) as total;
		`

		res, err := tx.Run(ctx, query, nil)
		if err != nil {
			return nil, fmt.Errorf("error running query: %w", err)
		}

		col, err := res.Collect(ctx)
		if err != nil {
			return nil, fmt.Errorf("error collecting results: %w", err)
		}

		out := make(map[string]int)
		for _, row := range col {
			day, ok := row.Get("day")
			if !ok {
				return nil, fmt.Errorf("day not found in result")
			}

			total, ok := row.Get("total")
			if !ok {
				return nil, fmt.Errorf("total not found in result")
			}

			dayStr := day.(neo4j.Date).String()
			formatted := strings.ReplaceAll(dayStr, "-", "")

			out[formatted] = int(total.(int64))
		}

		return out, nil
	})
	if err != nil {
		return nil, fmt.Errorf("error executing read: %w", err)
	}

	return totals, nil
}
