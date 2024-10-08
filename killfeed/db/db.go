package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPool interface {
	Close()
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Ping(ctx context.Context) error
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type Client struct {
	pool *pgxpool.Pool
}

func NewClient(ctx context.Context, connString string) (*Client, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error parsing connection string: %w", err)
	}

	// https://github.com/jackc/pgx/issues/1847
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeCacheDescribe

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging db: %w", err)
	}

	return &Client{pool}, nil
}

func (c *Client) Close() {
	c.pool.Close()
}

// FetchTotals returns a map of unique days to the number of killmails on that day
func (c *Client) FetchTotals(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT to_char(date_trunc('day', time), 'YYYYMMDD') as day, count(*) as total
		FROM killmail.esi_killmail
		GROUP BY day;
	`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying day_total: %w", err)
	}

	totals := make(map[string]int)
	for rows.Next() {
		var day string
		var total int
		if err := rows.Scan(&day, &total); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		totals[day] = total
	}

	return totals, nil
}
