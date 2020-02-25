package hardware

import (
	"errors"
	hardware_wallet "github.com/fibercrypto/fibercryptowallet/src/contrib/hardware-wallet"
	"github.com/fibercrypto/fibercryptowallet/src/contrib/hardware-wallet/skywallet/mocks"
	"github.com/fibercrypto/skywallet-go/src/skywallet"
	"github.com/fibercrypto/skywallet-go/src/skywallet/wire"
	messages "github.com/fibercrypto/skywallet-protob/go"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func createDeviceInteraction2(dev skywallet.Devicer) hardware_wallet.DeviceInteraction {
	return &SkyWalletInteraction{
		dev: dev,
		initializeWasWarn: false,
		uploadFirmwareWasWarn: false,
		secureWasWarn: false,
	}
}

func TestAddressGenShouldWorkOk(t *testing.T) {
	// Giving
	dev := &mocks.Devicer{}
	orgAddrs := []string{"jhfjdhfjd", "dfd787fd8"}
	addrResp := messages.ResponseSkycoinAddress{Addresses: orgAddrs}
	data, err := proto.Marshal(&addrResp)
	require.NoError(t, err)
	msg := wire.Message{
		Kind: uint16(messages.MessageType_MessageType_ResponseSkycoinAddress),
		Data: data,
	}
	dev.On("AddressGen", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(msg, nil)
	di := createDeviceInteraction2(dev)

	// When
	addrs, err := di.AddressGen(uint32(1), uint32(1), false, "deterministic").Then(func(data interface{}) interface{} {
		return data
	}).Await()

	// Then
	require.NoError(t, err)
	addrStrs, ok := addrs.([]string)
	require.True(t, ok)
	require.Len(t, addrStrs, 2)
	require.Equal(t, orgAddrs[0], addrStrs[0])
	require.Equal(t, orgAddrs[1], addrStrs[1])
}

func TestAddressGenShouldHandleDeviceError(t *testing.T) {
	// Giving
	dev := &mocks.Devicer{}
	dev.On("AddressGen", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(wire.Message{}, errors.New(""))
	di := createDeviceInteraction2(dev)

	// When
	_, err := di.AddressGen(uint32(1), uint32(1), false, "deterministic").Then(func(data interface{}) interface{} {
		return data
	}).Await()

	// Then
	require.Error(t, err)
}

func TestAddressGenShouldHandleErrorDecodingResponse(t *testing.T) {
	// Giving
	dev := &mocks.Devicer{}
	orgAddrs := []string{"jhfjdhfjd", "dfd787fd8"}
	addrResp := messages.ResponseSkycoinAddress{Addresses: orgAddrs}
	data, err := proto.Marshal(&addrResp)
	require.NoError(t, err)
	msg := wire.Message{
		Kind: uint16(messages.MessageType_MessageType_Success),
		Data: data,
	}
	dev.On("AddressGen", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(msg, nil)
	di := createDeviceInteraction2(dev)

	// When
	_, err = di.AddressGen(uint32(1), uint32(1), false, "deterministic").Then(func(data interface{}) interface{} {
		return data
	}).Await()

	// Then
	require.Error(t, err)
	require.Equal(t, err, errors.New("calling DecodeResponseSkycoinAddress with wrong message type: MessageType_Success"))
}