// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	core "github.com/fibercrypto/fibercryptowallet/src/core"
	mock "github.com/stretchr/testify/mock"
)

// BlockchainStatus is an autogenerated mock type for the BlockchainStatus type
type BlockchainStatus struct {
	mock.Mock
}

// GetCoinValue provides a mock function with given fields: coinvalue, ticker
func (_m *BlockchainStatus) GetCoinValue(coinvalue core.CoinValueMetric, ticker string) (uint64, error) {
	ret := _m.Called(coinvalue, ticker)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(core.CoinValueMetric, string) uint64); ok {
		r0 = rf(coinvalue, ticker)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(core.CoinValueMetric, string) error); ok {
		r1 = rf(coinvalue, ticker)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLastBlock provides a mock function with given fields:
func (_m *BlockchainStatus) GetLastBlock() (core.Block, error) {
	ret := _m.Called()

	var r0 core.Block
	if rf, ok := ret.Get(0).(func() core.Block); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.Block)
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

// GetNumberOfBlocks provides a mock function with given fields:
func (_m *BlockchainStatus) GetNumberOfBlocks() (uint64, error) {
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
