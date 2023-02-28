package authservice

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) CreateEntity(
	ctx context.Context,
	adminSecret string,
	entitySecret string,
	perms Permissions,
	note string,
) (Entity, error) {
	if adminSecret != s.cfg.AdminSecret {
		return Entity{}, ErrUnauthorized
	}

	permsJSON, err := json.Marshal(perms)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to marshal permissions: %w", err)
	}

	if len(perms) == 0 {
		return Entity{}, ErrInvalidArg{Msg: "empty permissions set"}
	}

	for k, v := range perms {
		if !(v.Read || v.Write) {
			return Entity{}, ErrInvalidArg{Msg: "empty permissions set for " + k}
		}
	}

	sec, err := bcrypt.GenerateFromPassword([]byte(entitySecret), bcrypt.DefaultCost)
	if err != nil {
		return Entity{}, fmt.Errorf("failed to hash a secret: %w", err)
	}

	tx, err := s.db.Begin()
	if err != nil {
		return Entity{}, err
	}

	id := uuid.NewString()

	q := `INSERT INTO auth_entity (id, secret, permissions, note) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, q, id, sec, permsJSON, note)
	if err != nil {
		return Entity{}, err
	}

	if err := tx.Commit(); err != nil {
		return Entity{}, err
	}

	return s.GetEntity(ctx, id)
}
