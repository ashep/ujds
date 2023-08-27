// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package indexhandler_test

import (
	"context"
	"github.com/ashep/ujds/internal/model"
	"sync"
)

// indexRepoMock is a mock implementation of indexhandler.indexRepo.
//
//	func TestSomethingThatUsesindexRepo(t *testing.T) {
//
//		// make and configure a mocked indexhandler.indexRepo
//		mockedindexRepo := &indexRepoMock{
//			ClearFunc: func(ctx context.Context, name string) error {
//				panic("mock out the Clear method")
//			},
//			GetFunc: func(ctx context.Context, name string) (model.Index, error) {
//				panic("mock out the Get method")
//			},
//			UpsertFunc: func(ctx context.Context, name string, schema string) error {
//				panic("mock out the Upsert method")
//			},
//		}
//
//		// use mockedindexRepo in code that requires indexhandler.indexRepo
//		// and then make assertions.
//
//	}
type indexRepoMock struct {
	// ClearFunc mocks the Clear method.
	ClearFunc func(ctx context.Context, name string) error

	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, name string) (model.Index, error)

	// UpsertFunc mocks the Upsert method.
	UpsertFunc func(ctx context.Context, name string, schema string) error

	// calls tracks calls to the methods.
	calls struct {
		// Clear holds details about calls to the Clear method.
		Clear []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
		// Upsert holds details about calls to the Upsert method.
		Upsert []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
			// Schema is the schema argument value.
			Schema string
		}
	}
	lockClear  sync.RWMutex
	lockGet    sync.RWMutex
	lockUpsert sync.RWMutex
}

// Clear calls ClearFunc.
func (mock *indexRepoMock) Clear(ctx context.Context, name string) error {
	if mock.ClearFunc == nil {
		panic("indexRepoMock.ClearFunc: method is nil but indexRepo.Clear was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockClear.Lock()
	mock.calls.Clear = append(mock.calls.Clear, callInfo)
	mock.lockClear.Unlock()
	return mock.ClearFunc(ctx, name)
}

// ClearCalls gets all the calls that were made to Clear.
// Check the length with:
//
//	len(mockedindexRepo.ClearCalls())
func (mock *indexRepoMock) ClearCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockClear.RLock()
	calls = mock.calls.Clear
	mock.lockClear.RUnlock()
	return calls
}

// Get calls GetFunc.
func (mock *indexRepoMock) Get(ctx context.Context, name string) (model.Index, error) {
	if mock.GetFunc == nil {
		panic("indexRepoMock.GetFunc: method is nil but indexRepo.Get was just called")
	}
	callInfo := struct {
		Ctx  context.Context
		Name string
	}{
		Ctx:  ctx,
		Name: name,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, name)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedindexRepo.GetCalls())
func (mock *indexRepoMock) GetCalls() []struct {
	Ctx  context.Context
	Name string
} {
	var calls []struct {
		Ctx  context.Context
		Name string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// Upsert calls UpsertFunc.
func (mock *indexRepoMock) Upsert(ctx context.Context, name string, schema string) error {
	if mock.UpsertFunc == nil {
		panic("indexRepoMock.UpsertFunc: method is nil but indexRepo.Upsert was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Name   string
		Schema string
	}{
		Ctx:    ctx,
		Name:   name,
		Schema: schema,
	}
	mock.lockUpsert.Lock()
	mock.calls.Upsert = append(mock.calls.Upsert, callInfo)
	mock.lockUpsert.Unlock()
	return mock.UpsertFunc(ctx, name, schema)
}

// UpsertCalls gets all the calls that were made to Upsert.
// Check the length with:
//
//	len(mockedindexRepo.UpsertCalls())
func (mock *indexRepoMock) UpsertCalls() []struct {
	Ctx    context.Context
	Name   string
	Schema string
} {
	var calls []struct {
		Ctx    context.Context
		Name   string
		Schema string
	}
	mock.lockUpsert.RLock()
	calls = mock.calls.Upsert
	mock.lockUpsert.RUnlock()
	return calls
}
