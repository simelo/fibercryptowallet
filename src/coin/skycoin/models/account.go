package models

import (
	"github.com/fibercrypto/FiberCryptoWallet/src/core"
	"github.com/fibercrypto/FiberCryptoWallet/src/util"
)

func (addr SkycoinAddress) GetBalance(ticker string) (uint64, error) {
	c := util.NewClient()
	bl, err := c.Balance([]string{addr.address})
	if err != nil {
		return 0, err
	}

	if ticker == Sky {
		return bl.Confirmed.Coins, nil
	} else if ticker == CoinHour {
		return bl.Confirmed.Hours, nil
	} else {
		return 0, errorTickerInvalid{ticker}
	}
}
func (addr SkycoinAddress) ListAssets() []string {
	return []string{Sky, CoinHour}
}
func (addr SkycoinAddress) ScanUnspentOutputs() core.TransactionOutputIterator { //------TODO
	return nil
}
func (addr SkycoinAddress) ListTransactions() core.TransactionIterator { //------TODO
	return nil
}

func (wlt RemoteWallet) GetBalance(ticker string) (uint64, error) {
	c := wlt.newClient()
	bl, err := c.WalletBalance(wlt.Id)
	if err != nil {
		return 0, err
	}

	if ticker == Sky {
		return bl.Confirmed.Coins, nil
	} else if ticker == CoinHour {
		return bl.Confirmed.Hours, nil
	} else {
		return 0, errorTickerInvalid{ticker}
	}

}

func (wlt RemoteWallet) ListAssets() []string {
	return []string{Sky, CoinHour}
}

func (wlt RemoteWallet) ScanUnspentOutputs() core.TransactionOutputIterator { //------TODO
	return nil
}

func (wlt RemoteWallet) ListTransactions() core.TransactionIterator { //------TODO
	return nil
}
