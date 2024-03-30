package indexhandler_test

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/ashep/ujds/internal/model"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Upsert(ctx context.Context, name, title, schema string) error {
	args := m.Called(ctx, name, title, schema)
	return args.Error(0)
}

func (m *repoMock) Get(ctx context.Context, name string) (model.Index, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(model.Index), args.Error(1)
}

func (m *repoMock) List(ctx context.Context) ([]model.Index, error) {
	args := m.Called(ctx)
	return args.Get(0).([]model.Index), args.Error(1)
}

func (m *repoMock) Clear(ctx context.Context, name string) error {
	args := m.Called(ctx, name)
	return args.Error(0)
}
