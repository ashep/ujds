package authservice

import (
	"context"
	"encoding/json"
)

func (s *Service) GetEntity(ctx context.Context, id string) (Entity, error) {
	row := s.db.QueryRowContext(ctx, `SELECT permissions, note FROM auth_entity WHERE id=$1`, id)

	var perms []byte
	var note string
	if err := row.Scan(&perms, &note); err != nil {
		return Entity{}, err
	}

	e := Entity{Id: id, Note: note}
	if err := json.Unmarshal(perms, &e.Permissions); err != nil {
		return Entity{}, err
	}

	return e, nil
}
