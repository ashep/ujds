package indexhandler_test

import (
	"context"

	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/stretchr/testify/mock"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Upsert(ctx context.Context, name, title, schema string) error {
	args := m.Called(ctx, name, title, schema)
	return args.Error(0)
}

func (m *repoMock) Get(ctx context.Context, name string) (indexrepo.Index, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(indexrepo.Index), args.Error(1)
}

func (m *repoMock) List(ctx context.Context) ([]indexrepo.Index, error) {
	args := m.Called(ctx)
	return args.Get(0).([]indexrepo.Index), args.Error(1)
}

func (m *repoMock) Clear(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}
