package dataservice

import (
	"time"
)

type Item struct {
	Id      string
	Type    string
	Version uint64
	Time    time.Time
	Data    []byte
}

func (s *Service) ValidateItemData(tp string, data []byte) error {
	return nil
}
