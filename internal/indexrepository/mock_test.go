// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package indexrepository_test

import (
	"sync"
)

// stringValidatorMock is a mock implementation of indexrepository.stringValidator.
//
//	func TestSomethingThatUsesstringValidator(t *testing.T) {
//
//		// make and configure a mocked indexrepository.stringValidator
//		mockedstringValidator := &stringValidatorMock{
//			ValidateFunc: func(s string) error {
//				panic("mock out the Validate method")
//			},
//		}
//
//		// use mockedstringValidator in code that requires indexrepository.stringValidator
//		// and then make assertions.
//
//	}
type stringValidatorMock struct {
	// ValidateFunc mocks the Validate method.
	ValidateFunc func(s string) error

	// calls tracks calls to the methods.
	calls struct {
		// Validate holds details about calls to the Validate method.
		Validate []struct {
			// S is the s argument value.
			S string
		}
	}
	lockValidate sync.RWMutex
}

// Validate calls ValidateFunc.
func (mock *stringValidatorMock) Validate(s string) error {
	if mock.ValidateFunc == nil {
		panic("stringValidatorMock.ValidateFunc: method is nil but stringValidator.Validate was just called")
	}
	callInfo := struct {
		S string
	}{
		S: s,
	}
	mock.lockValidate.Lock()
	mock.calls.Validate = append(mock.calls.Validate, callInfo)
	mock.lockValidate.Unlock()
	return mock.ValidateFunc(s)
}

// ValidateCalls gets all the calls that were made to Validate.
// Check the length with:
//
//	len(mockedstringValidator.ValidateCalls())
func (mock *stringValidatorMock) ValidateCalls() []struct {
	S string
} {
	var calls []struct {
		S string
	}
	mock.lockValidate.RLock()
	calls = mock.calls.Validate
	mock.lockValidate.RUnlock()
	return calls
}