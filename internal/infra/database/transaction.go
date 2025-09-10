package database

import (
	"context"
	"database/sql"
	"fmt"
)

type SqlExecer interface {
	Exec(query string, args ...any) (sql.Result, error)
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
}

type TransactionManager interface {
	StartTransaction(ctx context.Context) (*sql.Tx, error)
	Commit(tx *sql.Tx) error
	Rollback(tx *sql.Tx) error
	WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error
}

type transactionManager struct {
	db *sql.DB
}

func NewTransactionManager(db *sql.DB) TransactionManager {
	return &transactionManager{db: db}
}

func (t *transactionManager) StartTransaction(ctx context.Context) (*sql.Tx, error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	return tx, nil
}

func (t *transactionManager) Commit(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (t *transactionManager) Rollback(tx *sql.Tx) error {
	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}
	return nil
}

func (t *transactionManager) WithTransaction(ctx context.Context, fn func(tx *sql.Tx) error) error {
	tx, err := t.StartTransaction(ctx)
	if err != nil {
		return err
	}
	defer t.Rollback(tx)

	if err := fn(tx); err != nil {
		return err
	}

	if err := t.Commit(tx); err != nil {
		return err
	}

	return nil
}
