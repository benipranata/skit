package skit

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	RequireNil(t, err)
	_, err = db.Exec(`CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT NOT NULL)`)
	RequireNil(t, err)
	return db
}

func countUsers(t *testing.T, db *sql.DB) int {
	t.Helper()
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	RequireNil(t, err)
	return count
}

func TestExecTx(t *testing.T) {
	tests := []struct {
		name      string
		ctx       func() context.Context
		opts      *sql.TxOptions
		fn        func(tx *sql.Tx) error
		wantErr   bool
		errMsg    string
		wantCount int
	}{
		{
			name: "successful insert",
			fn: func(tx *sql.Tx) error {
				_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "alice")
				return err
			},
			wantCount: 1,
		},
		{
			name: "successful multiple inserts",
			fn: func(tx *sql.Tx) error {
				if _, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "bob"); err != nil {
					return err
				}
				_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "charlie")
				return err
			},
			wantCount: 2,
		},
		{
			name: "empty fn succeeds",
			fn: func(tx *sql.Tx) error {
				return nil
			},
			wantCount: 0,
		},
		{
			name: "with tx options",
			opts: &sql.TxOptions{Isolation: sql.LevelSerializable},
			fn: func(tx *sql.Tx) error {
				_, err := tx.Exec("INSERT INTO users (name) VALUES (?)", "frank")
				return err
			},
			wantCount: 1,
		},
		{
			name: "rollback on fn error",
			fn: func(tx *sql.Tx) error {
				_, _ = tx.Exec("INSERT INTO users (name) VALUES (?)", "dave")
				return errors.New("business logic error")
			},
			wantErr:   true,
			errMsg:    "fn logic",
			wantCount: 0,
		},
		{
			name: "rollback on sql error",
			fn: func(tx *sql.Tx) error {
				_, _ = tx.Exec("INSERT INTO users (name) VALUES (?)", "eve")
				_, err := tx.Exec("INSERT INTO nonexistent_table (name) VALUES (?)", "fail")
				return err
			},
			wantErr:   true,
			errMsg:    "fn logic",
			wantCount: 0,
		},
		{
			name: "commit fails after manual rollback",
			fn: func(tx *sql.Tx) error {
				_ = tx.Rollback() // manually rollback, causing commit to fail
				return nil
			},
			wantErr:   true,
			errMsg:    "commit tx",
			wantCount: 0,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			fn:        func(tx *sql.Tx) error { return nil },
			wantErr:   true,
			errMsg:    "begin tx",
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := setupTestDB(t)
			defer db.Close()

			ctx := context.Background()
			if tt.ctx != nil {
				ctx = tt.ctx()
			}

			err := ExecTx(ctx, db, tt.opts, tt.fn)

			if tt.wantErr {
				RequireNotNil(t, err)
				require.Contains(t, err.Error(), tt.errMsg)
			} else {
				RequireNil(t, err)
			}
			require.Equal(t, tt.wantCount, countUsers(t, db))
		})
	}
}

func TestExecTx_Panic(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	require.Panics(t, func() {
		_ = ExecTx(context.Background(), db, nil, func(tx *sql.Tx) error {
			_, _ = tx.Exec("INSERT INTO users (name) VALUES (?)", "panic_user")
			panic("unexpected panic")
		})
	})

	require.Equal(t, 0, countUsers(t, db), "should rollback on panic")
}
