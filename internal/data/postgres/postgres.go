package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// ConnectPool establishes a connection pool to the PostgreSQL database using the provided connection string.
func ConnectPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	// Parse the connection string into a configuration struct.
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse connection string: %w", err)
	}

	// Optional: Configure pool settings.
	// config.MaxConns = 10                      // Maximum number of connections.
	// config.MinConns = 2                       // Minimum number of connections.
	// config.MaxConnIdleTime = 30 * time.Minute // Idle time before closing a connection.

	// Create the connection pool.
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Verify connectivity with a ping.
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	return pool, nil
}
