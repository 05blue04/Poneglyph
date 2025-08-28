package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

func deleteRecord(db *sql.DB, tableName string, id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
