package recordhandler_test

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/ashep/ujds/internal/model"
)

type indexRepoMock struct {
	mock.Mock
}

func (m *indexRepoMock) Get(ctx context.Context, name string) (model.Index, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(model.Index), args.Error(1)
}

type recordRepoMock struct {
	mock.Mock
}

func (m *recordRepoMock) Push(ctx context.Context, records []model.RecordUpdate) error {
	args := m.Called(ctx, records)
	return args.Error(0)
}

func (m *recordRepoMock) Get(
	ctx context.Context,
	index string,
	id string,
) (model.Record, error) {
	args := m.Called(ctx, index, id)
	return args.Get(0).(model.Record), args.Error(1)
}

func (m *recordRepoMock) Find(
	ctx context.Context,
	index string,
	search string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]model.Record, uint64, error) {
	args := m.Called(ctx, index, search, since, cursor, limit)
	return args.Get(0).([]model.Record), args.Get(1).(uint64), args.Error(2)
}

func (m *recordRepoMock) History(
	ctx context.Context,
	index string,
	id string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]model.Record, uint64, error) {
	args := m.Called(ctx, index, id, since, cursor, limit)
	return args.Get(0).([]model.Record), args.Get(1).(uint64), args.Error(2)
}
