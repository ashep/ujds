// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package recordhandler_test

import (
	"context"
	"github.com/ashep/ujds/internal/model"
	"sync"
	"time"
)

// indexRepoMock is a mock implementation of recordhandler.indexRepo.
//
//	func TestSomethingThatUsesindexRepo(t *testing.T) {
//
//		// make and configure a mocked recordhandler.indexRepo
//		mockedindexRepo := &indexRepoMock{
//			GetFunc: func(ctx context.Context, name string) (model.Index, error) {
//				panic("mock out the Get method")
//			},
//		}
//
//		// use mockedindexRepo in code that requires recordhandler.indexRepo
//		// and then make assertions.
//
//	}
type indexRepoMock struct {
	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, name string) (model.Index, error)

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Name is the name argument value.
			Name string
		}
	}
	lockGet sync.RWMutex
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

// recordRepoMock is a mock implementation of recordhandler.recordRepo.
//
//	func TestSomethingThatUsesrecordRepo(t *testing.T) {
//
//		// make and configure a mocked recordhandler.recordRepo
//		mockedrecordRepo := &recordRepoMock{
//			GetFunc: func(ctx context.Context, index string, id string) (model.Record, error) {
//				panic("mock out the Get method")
//			},
//			GetAllFunc: func(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
//				panic("mock out the GetAll method")
//			},
//			HistoryFunc: func(ctx context.Context, index string, id string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
//				panic("mock out the History method")
//			},
//			PushFunc: func(ctx context.Context, indexID uint64, schema []byte, records []model.RecordUpdate) error {
//				panic("mock out the Push method")
//			},
//		}
//
//		// use mockedrecordRepo in code that requires recordhandler.recordRepo
//		// and then make assertions.
//
//	}
type recordRepoMock struct {
	// GetFunc mocks the Get method.
	GetFunc func(ctx context.Context, index string, id string) (model.Record, error)

	// GetAllFunc mocks the GetAll method.
	GetAllFunc func(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error)

	// HistoryFunc mocks the History method.
	HistoryFunc func(ctx context.Context, index string, id string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error)

	// PushFunc mocks the Push method.
	PushFunc func(ctx context.Context, indexID uint64, schema []byte, records []model.RecordUpdate) error

	// calls tracks calls to the methods.
	calls struct {
		// Get holds details about calls to the Get method.
		Get []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Index is the index argument value.
			Index string
			// ID is the id argument value.
			ID string
		}
		// GetAll holds details about calls to the GetAll method.
		GetAll []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Index is the index argument value.
			Index string
			// Since is the since argument value.
			Since time.Time
			// Cursor is the cursor argument value.
			Cursor uint64
			// Limit is the limit argument value.
			Limit uint32
		}
		// History holds details about calls to the History method.
		History []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Index is the index argument value.
			Index string
			// ID is the id argument value.
			ID string
			// Since is the since argument value.
			Since time.Time
			// Cursor is the cursor argument value.
			Cursor uint64
			// Limit is the limit argument value.
			Limit uint32
		}
		// Push holds details about calls to the Push method.
		Push []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// IndexID is the indexID argument value.
			IndexID uint64
			// Schema is the schema argument value.
			Schema []byte
			// Records is the records argument value.
			Records []model.RecordUpdate
		}
	}
	lockGet     sync.RWMutex
	lockGetAll  sync.RWMutex
	lockHistory sync.RWMutex
	lockPush    sync.RWMutex
}

// Get calls GetFunc.
func (mock *recordRepoMock) Get(ctx context.Context, index string, id string) (model.Record, error) {
	if mock.GetFunc == nil {
		panic("recordRepoMock.GetFunc: method is nil but recordRepo.Get was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Index string
		ID    string
	}{
		Ctx:   ctx,
		Index: index,
		ID:    id,
	}
	mock.lockGet.Lock()
	mock.calls.Get = append(mock.calls.Get, callInfo)
	mock.lockGet.Unlock()
	return mock.GetFunc(ctx, index, id)
}

// GetCalls gets all the calls that were made to Get.
// Check the length with:
//
//	len(mockedrecordRepo.GetCalls())
func (mock *recordRepoMock) GetCalls() []struct {
	Ctx   context.Context
	Index string
	ID    string
} {
	var calls []struct {
		Ctx   context.Context
		Index string
		ID    string
	}
	mock.lockGet.RLock()
	calls = mock.calls.Get
	mock.lockGet.RUnlock()
	return calls
}

// GetAll calls GetAllFunc.
func (mock *recordRepoMock) GetAll(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
	if mock.GetAllFunc == nil {
		panic("recordRepoMock.GetAllFunc: method is nil but recordRepo.GetAll was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Index  string
		Since  time.Time
		Cursor uint64
		Limit  uint32
	}{
		Ctx:    ctx,
		Index:  index,
		Since:  since,
		Cursor: cursor,
		Limit:  limit,
	}
	mock.lockGetAll.Lock()
	mock.calls.GetAll = append(mock.calls.GetAll, callInfo)
	mock.lockGetAll.Unlock()
	return mock.GetAllFunc(ctx, index, since, cursor, limit)
}

// GetAllCalls gets all the calls that were made to GetAll.
// Check the length with:
//
//	len(mockedrecordRepo.GetAllCalls())
func (mock *recordRepoMock) GetAllCalls() []struct {
	Ctx    context.Context
	Index  string
	Since  time.Time
	Cursor uint64
	Limit  uint32
} {
	var calls []struct {
		Ctx    context.Context
		Index  string
		Since  time.Time
		Cursor uint64
		Limit  uint32
	}
	mock.lockGetAll.RLock()
	calls = mock.calls.GetAll
	mock.lockGetAll.RUnlock()
	return calls
}

// History calls HistoryFunc.
func (mock *recordRepoMock) History(ctx context.Context, index string, id string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
	if mock.HistoryFunc == nil {
		panic("recordRepoMock.HistoryFunc: method is nil but recordRepo.History was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		Index  string
		ID     string
		Since  time.Time
		Cursor uint64
		Limit  uint32
	}{
		Ctx:    ctx,
		Index:  index,
		ID:     id,
		Since:  since,
		Cursor: cursor,
		Limit:  limit,
	}
	mock.lockHistory.Lock()
	mock.calls.History = append(mock.calls.History, callInfo)
	mock.lockHistory.Unlock()
	return mock.HistoryFunc(ctx, index, id, since, cursor, limit)
}

// HistoryCalls gets all the calls that were made to History.
// Check the length with:
//
//	len(mockedrecordRepo.HistoryCalls())
func (mock *recordRepoMock) HistoryCalls() []struct {
	Ctx    context.Context
	Index  string
	ID     string
	Since  time.Time
	Cursor uint64
	Limit  uint32
} {
	var calls []struct {
		Ctx    context.Context
		Index  string
		ID     string
		Since  time.Time
		Cursor uint64
		Limit  uint32
	}
	mock.lockHistory.RLock()
	calls = mock.calls.History
	mock.lockHistory.RUnlock()
	return calls
}

// Push calls PushFunc.
func (mock *recordRepoMock) Push(ctx context.Context, indexID uint64, schema []byte, records []model.RecordUpdate) error {
	if mock.PushFunc == nil {
		panic("recordRepoMock.PushFunc: method is nil but recordRepo.Push was just called")
	}
	callInfo := struct {
		Ctx     context.Context
		IndexID uint64
		Schema  []byte
		Records []model.RecordUpdate
	}{
		Ctx:     ctx,
		IndexID: indexID,
		Schema:  schema,
		Records: records,
	}
	mock.lockPush.Lock()
	mock.calls.Push = append(mock.calls.Push, callInfo)
	mock.lockPush.Unlock()
	return mock.PushFunc(ctx, indexID, schema, records)
}

// PushCalls gets all the calls that were made to Push.
// Check the length with:
//
//	len(mockedrecordRepo.PushCalls())
func (mock *recordRepoMock) PushCalls() []struct {
	Ctx     context.Context
	IndexID uint64
	Schema  []byte
	Records []model.RecordUpdate
} {
	var calls []struct {
		Ctx     context.Context
		IndexID uint64
		Schema  []byte
		Records []model.RecordUpdate
	}
	mock.lockPush.RLock()
	calls = mock.calls.Push
	mock.lockPush.RUnlock()
	return calls
}
