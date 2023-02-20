package dataservice

import (
	"context"
)

func (s *Service) UpsertSchema(ctx context.Context, schema []byte) error {
	return nil
}

func (s *Service) GetSchema(ctx context.Context, tp string) ([]byte, error) {
	return nil, nil
}
