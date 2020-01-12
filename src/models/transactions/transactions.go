package transactions

import (
	"fmt"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	"github.com/fibercrypto/fibercryptowallet/src/models/address"
	modelUtil "github.com/fibercrypto/fibercryptowallet/src/models/util"
	"github.com/fibercrypto/fibercryptowallet/src/util"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	qtCore "github.com/therecipe/qt/core"
	"time"
)

var logTransactionDetails = logging.MustGetLogger("TransactionDetails")

func init() {
	TransactionDetails_QmlRegisterType2("HistoryModels", 1, 0, "QTransactionDetail")
}

const (
	Date = int(qtCore.Qt__UserRole) + 1<<iota
	Status
	Type
	Amount
	HoursTraspassed
	HoursBurned
	TransactionID
	BlockHeight
	Addresses
	Inputs
	Outputs
	CoinOptions
)

const (
	TransactionStatusConfirmed = iota
	TransactionStatusPending
	TransactionStatusPreview
)

const (
	TransactionTypeSend = iota
	TransactionTypeReceive
	TransactionTypeInternal
	TransactionTypeGeneric
)

type TransactionDetails struct {
	qtCore.QObject
	_ *qtCore.QDateTime    `property:"date"`
	_ int                  `property:"status"`
	_ int                  `property:"type"`
	_ uint64               `property:"blockHeight"`
	_ string               `property:"amount"`
	_ string               `property:"hoursTraspassed"`
	_ string               `property:"hoursBurned"`
	_ string               `property:"transactionID"`
	_ *address.AddressList `property:"addresses"`
	_ *address.AddressList `property:"inputs"`
	_ *address.AddressList `property:"outputs"`
	_ modelUtil.Map        `property:"coinOptions"`
}

func NewTransactionDetailFromCoreTransaction(transaction core.Transaction, txType int) (*TransactionDetails, error) {

	txnDetails := NewTransactionDetails(nil)
	t := time.Unix(int64(transaction.GetTimestamp()), 0)

	txnDetails.SetDate(qtCore.NewQDateTime3(qtCore.NewQDate3(t.Year(), int(t.Month()), t.Day()),
		qtCore.NewQTime3(t.Hour(), t.Minute(), 0, 0), qtCore.Qt__LocalTime))
	switch transaction.GetStatus() {
	case core.TXN_STATUS_CONFIRMED:
		txnDetails.SetStatus(TransactionStatusConfirmed)
	case core.TXN_STATUS_PENDING:
		txnDetails.SetStatus(TransactionStatusPending)
	default:
		txnDetails.SetStatus(TransactionStatusPreview)
	}

	txnDetails.SetTransactionID(transaction.GetId())

	txnDetails.SetBlockHeight(transaction.GetBlockHeight())

	txnDetails.SetType(txType)

	addresses := address.NewAddressList(nil)
	inputList := address.NewAddressList(nil)
	outputsList := address.NewAddressList(nil)

	for _, input := range transaction.GetInputs() {

		qIn := address.NewAddressDetails(nil)
		qIn.SetAddress(input.GetSpentOutput().GetAddress().String())
		inputCoinOptions := modelUtil.NewMap(nil)

		for _, asset := range input.SupportedAssets() {
			inputCoin, err := input.GetCoins(asset)
			if err != nil {
				logTransactionDetails.WithError(err).Warnf("Couldn't get coin: %s", asset)
				continue
			}

			accuracy, err := util.AltcoinQuotient(asset)

			if err != nil {
				logTransactionDetails.WithError(err).Warnf("Couldn't get quotient of %s", asset)
				continue
			}

			inputCoinOptions.SetValue(asset, util.FormatCoins(inputCoin, accuracy))
		}

		qIn.SetCoinOptions(inputCoinOptions)
		inputList.AddAddress(qIn)
		addresses.AddAddress(qIn)
	}

	txnDetails.SetInputs(inputList)

	var totals = make(map[string]uint64)

	for _, out := range transaction.GetOutputs() {

		qOu := address.NewAddressDetails(nil)
		qOu.SetAddress(out.GetAddress().String())
		outputCoinOptions := modelUtil.NewMap(nil)

		for _, asset := range out.SupportedAssets() {

			outputCoin, err := out.GetCoins(asset)
			if err != nil {
				logTransactionDetails.WithError(err).Warnf("Couldn't get coin %s", asset)
				continue
			}

			accuracy, err := util.AltcoinQuotient(asset)
			if err != nil {
				logTransactionDetails.WithError(err).Warnf("Couldn't get quotient of %s", asset)
				continue
			}
			totals[asset] += outputCoin
			outputCoinOptions.SetValue(asset, util.FormatCoins(outputCoin, accuracy))
		}

		qOu.SetCoinOptions(outputCoinOptions)

		outputsList.AddAddress(qOu)
		addresses.AddAddress(qOu)
	}

	txnCoinOptions := modelUtil.NewMap(nil)

	for asset := range totals {
		fee, err := transaction.ComputeFee(asset)
		if err != nil {
			logTransactionDetails.WithError(err).Warnf("Couldn't get transaction fee of %s coin", asset)
		}

		totals[asset] += fee

		accuracy, err := util.AltcoinQuotient(asset)
		if err != nil {
			logTransactionDetails.WithError(err).Warnf("Couldn't get accuracy of coin %s", asset)
		}

		txnCoinOptions.SetValue(fmt.Sprintf("total %s", asset), util.FormatCoins(totals[asset], accuracy))
	}

	txnDetails.SetCoinOptions(txnCoinOptions)
	txnDetails.SetOutputs(outputsList)
	txnDetails.SetAddresses(addresses)

	return txnDetails, nil
}
