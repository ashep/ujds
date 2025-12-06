package recordhandler_test

import (
	"context"
	"time"

	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/ashep/ujds/internal/recordrepo"
	"github.com/stretchr/testify/mock"
)

type indexRepoMock struct {
	mock.Mock
}

func (m *indexRepoMock) Get(ctx context.Context, name string) (indexrepo.Index, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(indexrepo.Index), args.Error(1)
}

type recordRepoMock struct {
	mock.Mock
}

func (m *recordRepoMock) Push(ctx context.Context, records []recordrepo.RecordUpdate) error {
	args := m.Called(ctx, records)
	return args.Error(0)
}

func (m *recordRepoMock) Get(
	ctx context.Context,
	index string,
	id string,
) (recordrepo.Record, error) {
	args := m.Called(ctx, index, id)
	return args.Get(0).(recordrepo.Record), args.Error(1)
}

func (m *recordRepoMock) Find(ctx context.Context, req recordrepo.FindRequest) ([]recordrepo.Record, uint64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]recordrepo.Record), args.Get(1).(uint64), args.Error(2)
}

func (m *recordRepoMock) History(
	ctx context.Context,
	index string,
	id string,
	since time.Time,
	cursor uint64,
	limit uint32,
) ([]recordrepo.Record, uint64, error) {
	args := m.Called(ctx, index, id, since, cursor, limit)
	return args.Get(0).([]recordrepo.Record), args.Get(1).(uint64), args.Error(2)
}

type stringValidatorMock struct {
	mock.Mock
}

func (m *stringValidatorMock) Validate(data string) error {
	args := m.Called(data)
	return args.Error(0)
}
