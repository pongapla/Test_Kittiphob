package repository

import (
	"context"
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

var ErrQueueFull = errors.New("queue is full")

type Ticket struct {
	ID          string `json:"id"`
	Cycle       int    `json:"cycle"`
	QueueNumber string `json:"queueNumber"`
	IssuedAt    string `json:"issuedAt"`
}

type Repository struct {
	db *sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Serialize database operations on this connection pool.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	db.SetConnMaxLifetime(0)

	if _, err := db.Exec("PRAGMA busy_timeout = 5000"); err != nil {
		db.Close()
		return nil, err
	}

	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(string(schema)); err != nil {
		db.Close()
		return nil, err
	}

	return &Repository{db: db}, nil
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func formatNumber(index int) string {
	return fmt.Sprintf("%c%d", 'A'+rune(index/10), index%10)
}

func (r *Repository) Issue(ctx context.Context) (Ticket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Ticket{}, err
	}
	defer tx.Rollback()

	var index, cycle int

	err = tx.QueryRowContext(ctx,
		"SELECT current_index, cycle FROM queue_state WHERE id = 1",
	).Scan(&index, &cycle)
	if err != nil {
		return Ticket{}, err
	}

	if index >= 259 {
		return Ticket{}, ErrQueueFull
	}

	index++

	randomID := make([]byte, 16)
	if _, err := rand.Read(randomID); err != nil {
		return Ticket{}, err
	}

	ticket := Ticket{
		ID:          hex.EncodeToString(randomID),
		Cycle:       cycle,
		QueueNumber: formatNumber(index),
		IssuedAt:    time.Now().UTC().Format(time.RFC3339Nano),
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE queue_state
		 SET current_index = ?, updated_at = ?
		 WHERE id = 1`,
		index, ticket.IssuedAt,
	)
	if err != nil {
		return Ticket{}, err
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO queue_tickets
		 (id, cycle, queue_number, issued_at)
		 VALUES (?, ?, ?, ?)`,
		ticket.ID, ticket.Cycle,
		ticket.QueueNumber, ticket.IssuedAt,
	)
	if err != nil {
		return Ticket{}, err
	}

	if err := tx.Commit(); err != nil {
		return Ticket{}, err
	}

	return ticket, nil
}

func (r *Repository) Reset(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		`UPDATE queue_state
		 SET current_index = -1,
		     cycle = cycle + 1,
		     updated_at = ?
		 WHERE id = 1`,
		time.Now().UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Current(ctx context.Context) (string, error) {
	var index int

	err := r.db.QueryRowContext(ctx,
		"SELECT current_index FROM queue_state WHERE id = 1",
	).Scan(&index)
	if err != nil {
		return "", err
	}

	if index == -1 {
		return "00", nil
	}

	return formatNumber(index), nil
}

func (r *Repository) GetTicket(
	ctx context.Context, id string,
) (Ticket, error) {
	var ticket Ticket

	err := r.db.QueryRowContext(ctx,
		`SELECT id, cycle, queue_number, issued_at
		 FROM queue_tickets WHERE id = ?`, id,
	).Scan(
		&ticket.ID,
		&ticket.Cycle,
		&ticket.QueueNumber,
		&ticket.IssuedAt,
	)

	return ticket, err
}
