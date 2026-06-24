package infrastructure

import (
	"context"
	"database/sql"
	"fmt"
)

func NewTxManager(db *sql.DB) *TxManager {
	return &TxManager{
		db: db,
	}
}

type TxManager struct {
	db *sql.DB
}

func (tm *TxManager) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {

	//Begin Datbase Transaction
	tx, err := tm.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	//Defer Rollback to end
	if _, err := tx.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to force immediate transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	//Execute function
	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error: %v, rollback failed: %w", err, rbErr)
		}
		return err
	}

	//Commit Changes
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
