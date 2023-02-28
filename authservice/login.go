package authservice

import (
	"context"
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Login(ctx context.Context, id, secret string) (string, error) {
	var secretHash []byte
	row := s.db.QueryRowContext(ctx, `SELECT secret FROM auth_entity WHERE id=$1`, id)
	err := row.Scan(&secretHash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrUnauthorized
	} else if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(secretHash, []byte(secret))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return "", ErrUnauthorized
	} else if err != nil {
		return "", err
	}

	tok, err := s.getToken(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return s.createToken(ctx, id)
	} else if err != nil {
		return "", err
	}

	return tok, nil
}

func (s *Service) getToken(ctx context.Context, entityId string) (string, error) {
	var tok string

	row := s.db.QueryRowContext(ctx, `SELECT token FROM auth_token WHERE auth_entity_id=$1`, entityId)

	if err := row.Scan(&tok); errors.Is(err, sql.ErrNoRows) {
		return "", ErrNotFound
	} else if err != nil {
		return "", err
	}

	return tok, nil
}

func (s *Service) createToken(ctx context.Context, entityId string) (string, error) {
	_, err := s.db.ExecContext(ctx, `INSERT INTO auth_token (auth_entity_id) VALUES ($1)`, entityId)
	if err != nil {
		return "", err
	}

	return s.getToken(ctx, entityId)
}
