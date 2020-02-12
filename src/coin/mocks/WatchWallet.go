// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import core "github.com/fibercrypto/fibercryptowallet/src/core"
import mock "github.com/stretchr/testify/mock"

// WatchWallet is an autogenerated mock type for the WatchWallet type
type WatchWallet struct {
	mock.Mock
}

// GenAddresses provides a mock function with given fields: accountIndex, addrType, startIndex, count, pwd
func (_m *WatchWallet) GenAddresses(accountIndex uint32, addrType core.AddressType, startIndex uint32, count uint32, pwd core.PasswordReader) (core.AddressIterator, error) {
	ret := _m.Called(accountIndex, addrType, startIndex, count, pwd)

	var r0 core.AddressIterator
	if rf, ok := ret.Get(0).(func(uint32, core.AddressType, uint32, uint32, core.PasswordReader) core.AddressIterator); ok {
		r0 = rf(accountIndex, addrType, startIndex, count, pwd)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.AddressIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint32, core.AddressType, uint32, uint32, core.PasswordReader) error); ok {
		r1 = rf(accountIndex, addrType, startIndex, count, pwd)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllLoadedAddresses provides a mock function with given fields:
func (_m *WatchWallet) GetAllLoadedAddresses() (core.AddressIterator, error) {
	ret := _m.Called()

	var r0 core.AddressIterator
	if rf, ok := ret.Get(0).(func() core.AddressIterator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.AddressIterator)
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

// GetLoadedAddressesForAccount provides a mock function with given fields: accountIndex, addrType
func (_m *WatchWallet) GetLoadedAddressesForAccount(accountIndex uint32, addrType core.AddressType) (core.AddressIterator, error) {
	ret := _m.Called(accountIndex, addrType)

	var r0 core.AddressIterator
	if rf, ok := ret.Get(0).(func(uint32, core.AddressType) core.AddressIterator); ok {
		r0 = rf(accountIndex, addrType)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.AddressIterator)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint32, core.AddressType) error); ok {
		r1 = rf(accountIndex, addrType)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
