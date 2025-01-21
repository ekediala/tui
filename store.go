package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type Note struct {
	ID    int64
	Title string
	Body  string
}

type Store struct {
	conn *sql.DB
}

func NewStore(ctx context.Context) (*Store, error) {
	s := Store{}
	err := s.Init(ctx)
	return &s, err
}

func (s *Store) Init(ctx context.Context) error {
	var err error

	if s.conn, err = sql.Open("sqlite", "./notes.db"); err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}

	if err = s.conn.PingContext(ctx); err != nil {
		return fmt.Errorf("pinging database: %w", err)
	}

	if _, err = s.conn.ExecContext(ctx, CreateTableStmt); err != nil {
		return fmt.Errorf("creating notes table: %w", err)
	}

	return nil
}

func (s *Store) GetNotes(ctx context.Context) ([]Note, error) {
	var notes = []Note{}

	r, err := s.conn.QueryContext(ctx, GetNotesStmt)
	if err != nil {
		return notes, fmt.Errorf("getting notes: %w", err)
	}

	defer r.Close()

	for r.Next() {
		var note Note
		err := r.Scan(&note.ID, &note.Title, &note.Body)
		if err != nil {
			return notes, fmt.Errorf("scanning notes: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil

}

func (s *Store) Upsert(ctx context.Context, note Note) error {
	if note.ID == 0 {
		note.ID = time.Now().UnixNano()
	}

	r, err := s.conn.ExecContext(ctx, UpsertNoteStmt, note.ID, note.Title, note.Body)

	if err != nil {
		return fmt.Errorf("upserting note: %w", err)
	}

	n, err := r.RowsAffected()
	if err != nil {
		return fmt.Errorf("getting affected rows: %w", err)
	}

	if n != 1 {
		return errors.New(fmt.Sprintf("query error: incorrect number of rows affected: %d rows affected", n))
	}

	return nil
}
