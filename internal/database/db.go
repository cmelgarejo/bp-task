package database

import (
	"bp-task/internal/models"
	"context"
	"log"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB holds the PostgreSQL database connection pool.
type DB struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

// NewDB initializes a new PostgreSQL database connection pool.
func NewDB(ctx context.Context, dbURL string) (*DB, error) {
	// Create a PostgreSQL connection pool
	poolConfig, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatalln("Unable to parse DATABASE_URL:", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalln("Unable to create connection pool:", err)
	}

	return &DB{pool: pool, ctx: ctx}, nil
}

// Close closes the database connection pool.
func (db *DB) Close() {
	db.pool.Close()
}

// InsertMetadata inserts a metadata from a ipfs address (cid) into the "ipfs_metadata" table.
func (db *DB) InsertMetadata(cid string, token []byte) error {
	_, err := db.pool.Exec(db.ctx, "INSERT INTO ipfs_metadata (cid, token) VALUES ($1, $2)", cid, token)
	return err
}

// GetToken returns the token associated with a ipfs address (cid).
func (db *DB) GetToken(cid string) (models.IPFSToken, error) {
	var token models.IPFSToken
	err := db.pool.QueryRow(db.ctx, "SELECT token FROM ipfs_metadata WHERE cid=$1", cid).Scan(&token)
	// slog.InfoContext(db.ctx, "Token Info", "cid", cid, "token", token)
	return token, err
}

// GetTokens returns all the tokens in the "tokens" table.
func (db *DB) GetTokens() ([]*models.IPFSMetadata, error) {
	var allMetadata []*models.IPFSMetadata
	rows, err := db.pool.Query(db.ctx, "SELECT cid, token FROM ipfs_metadata")
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var cid string
		var token models.IPFSToken
		err := rows.Scan(&cid, &token)
		if err != nil {
			return nil, err
		}
		allMetadata = append(allMetadata, &models.IPFSMetadata{Cid: cid, Token: token})
	}

	return allMetadata, nil
}
