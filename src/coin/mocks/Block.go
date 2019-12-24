// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	core "github.com/fibercrypto/fibercryptowallet/src/core"
	mock "github.com/stretchr/testify/mock"
)

// Block is an autogenerated mock type for the Block type
type Block struct {
	mock.Mock
}

// GetFee provides a mock function with given fields: ticker
func (_m *Block) GetFee(ticker string) (uint64, error) {
	ret := _m.Called(ticker)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(string) uint64); ok {
		r0 = rf(ticker)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(ticker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHash provides a mock function with given fields:
func (_m *Block) GetHash() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetHeight provides a mock function with given fields:
func (_m *Block) GetHeight() (uint64, error) {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetPrevHash provides a mock function with given fields:
func (_m *Block) GetPrevHash() ([]byte, error) {
	ret := _m.Called()

	var r0 []byte
	if rf, ok := ret.Get(0).(func() []byte); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTime provides a mock function with given fields:
func (_m *Block) GetTime() (core.Timestamp, error) {
	ret := _m.Called()

	var r0 core.Timestamp
	if rf, ok := ret.Get(0).(func() core.Timestamp); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(core.Timestamp)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVersion provides a mock function with given fields:
func (_m *Block) GetVersion() (uint32, error) {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsGenesisBlock provides a mock function with given fields:
func (_m *Block) IsGenesisBlock() (bool, error) {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
