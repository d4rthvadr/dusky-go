package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/d4rthvadr/dusky-go/internal/models"
)

type RoleStore struct {
	db *sql.DB
}

func (r *RoleStore) GetByName(ctx context.Context, name models.RoleStr) (*models.Role, error) {
	query := `SELECT id, name, description, level, created_at, updated_at FROM roles WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	var role models.Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description, &role.Level, &role.CreatedAt, &role.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("role not found: %w", err)
		}
		return nil, fmt.Errorf("error scanning role: %w", err)
	}

	return &role, nil
}
