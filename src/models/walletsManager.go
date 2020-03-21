package models

import (
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin"
	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/config"
	util2 "github.com/fibercrypto/fibercryptowallet/src/models/util"
	"github.com/fibercrypto/fibercryptowallet/src/models/wallets"
	"sync"
	"time"

	// "github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/params"
	"github.com/fibercrypto/fibercryptowallet/src/models/address"
	"github.com/fibercrypto/fibercryptowallet/src/models/transactions"
	"github.com/fibercrypto/fibercryptowallet/src/util"
	"sort"

	"github.com/therecipe/qt/qml"

	sky "github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	local "github.com/fibercrypto/fibercryptowallet/src/main"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"
	qtCore "github.com/therecipe/qt/core"
)

var logWalletManager = logging.MustGetLogger("modelsWalletManager")

var once sync.Once
var walletManager *WalletManager

var updtWltChan = make(chan wallets.UpdateWallet, 100)

type WalletManager struct {
	qtCore.QObject
	WalletEnv                core.WalletEnv
	SeedGenerator            core.SeedGenerator
	wallets                  []*wallets.QWallet
	addresseseByWallets      map[string][]*address.QAddress
	altManager               core.AltcoinManager
	signer                   core.BlockchainSignService
	transactionAPI           core.BlockchainTransactionAPI
	addressesAndWalletsMutex sync.Mutex
	outputsByAddressMutex    sync.Mutex
	walletsIterator          core.WalletIterator
	_                        func()                                                                                                                                                `constructor:"init"`
	_                        func(model *wallets.WalletModel)                                                                                                                      `slot:"loadWallets"`
	_                        func()                                                                                                                                                `slot:"updateAll"`
	_                        func()                                                                                                                                                `slot:"updateWalletEnvs"`
	_                        func(seed string, label string, walletType string, password string, scanN int) *wallets.QWallet                                                       `slot:"createEncryptedWallet"`
	_                        func(seed string, label string, walletType string, scanN int) *wallets.QWallet                                                                        `slot:"createUnencryptedWallet"`
	_                        func(entropy int) string                                                                                                                              `slot:"getNewSeed"`
	_                        func(seed string) int                                                                                                                                 `slot:"verifySeed"`
	_                        func(id string, n int, password string)                                                                                                               `slot:"newWalletAddress"`
	_                        func(id string, password string) int                                                                                                                  `slot:"encryptWallet"`
	_                        func(id string, password string) int                                                                                                                  `slot:"decryptWallet"`
	_                        func() []*wallets.QWallet                                                                                                                             `slot:"getWallets"`
	_                        func() []string                                                                                                                                       `slot:"getCurrencyList"`
	_                        func(wltIds, addresses []string, source string, pwd interface{}, index []int, qTxn *transactions.TransactionDetails) *transactions.TransactionDetails `slot:"signTxn"`
	_                        func(wltId string, destinationAddress string, amount, currency string) *transactions.TransactionDetails                                               `slot:"sendTo"`
	_                        func(id, label string)                                                                                                                                `slot:"editWalletLbl"`
	_                        func(txn *transactions.TransactionDetails) bool                                                                                                       `slot:"broadcastTxn"`
	_                        func(wltIds, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *transactions.TransactionDetails   `slot:"sendFromAddresses"`
	_                        func(wltIds, outs, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *transactions.TransactionDetails   `slot:"sendFromOutputs"`
	_                        func() string                                                                                                                                         `slot:"getDefaultWalletType"`
	_                        func(wltIds, addresses []string, source string, bridgeForPassword *QBridge, index []int, qTxn *transactions.TransactionDetails)                       `slot:"signAndBroadcastTxnAsync"`
	_                        func() []string                                                                                                                                       `slot:"getAvailableWalletTypes"`
	_                        func(wltName string, model *address.ModelAddress)                                                                                                     `slot:"loadAddressModelByWallet"`
	_                        func(model *address.ModelAddress)                                                                                                                     `slot:"loadAddressForAllWallets"`
	_                        func(model *wallets.WalletModel)                                                                                                                      `signal:"initWltModelAsync"`
	_                        func(wallet *wallets.QWallet)                                                                                                                         `slot:"addWltAsync"`
}

func (walletM *WalletManager) init() {
	logWalletManager.Info("Initializing WalletManager")
	logWalletManager.Debug("Starting WalletManager")
	once.Do(func() {
		qml.QQmlEngine_SetObjectOwnership(walletM, qml.QQmlEngine__CppOwnership)
		walletM.ConnectLoadWallets(walletM.loadWallets)
		walletM.ConnectEditWalletLbl(walletM.editWalletLbl)
		walletM.ConnectCreateEncryptedWallet(walletM.createEncryptedWallet)
		walletM.ConnectCreateUnencryptedWallet(walletM.createUnencryptedWallet)
		walletM.ConnectGetNewSeed(walletM.getNewSeed)
		walletM.ConnectVerifySeed(walletM.verifySeed)
		walletM.ConnectNewWalletAddress(walletM.newWalletAddress)
		walletM.ConnectEncryptWallet(walletM.encryptWallet)
		walletM.ConnectDecryptWallet(walletM.decryptWallet)
		walletM.ConnectGetWallets(walletM.getWallets)
		walletM.ConnectGetCurrencyList(walletM.getCurrencyList)
		walletM.ConnectSendTo(walletM.sendTo)
		walletM.ConnectSignTxn(walletM.signTxn)
		walletM.ConnectSendFromAddresses(walletM.sendFromAddresses)
		walletM.ConnectSendFromOutputs(walletM.sendFromOutputs)
		walletM.ConnectBroadcastTxn(walletM.broadcastTxn)
		walletM.ConnectUpdateWalletEnvs(walletM.updateWalletEnvs)
		walletM.ConnectUpdateAll(walletM.updateAll)
		walletM.ConnectSignAndBroadcastTxnAsync(walletM.signAndBroadcastTxnAsync)
		walletM.ConnectGetDefaultWalletType(walletM.getDefaultWalletType)
		walletM.ConnectGetAvailableWalletTypes(walletM.getAvailableWalletTypes)
		walletM.ConnectLoadAddressModelByWallet(walletM.loadAddressModelByWallet)
		walletM.ConnectLoadAddressForAllWallets(walletM.loadAddressForAllWallets)
		walletM.ConnectInitWltModelAsync(walletM.initWltModelAsync)
		walletM.ConnectAddWltAsync(walletM.addWalletAsync)

		walletM.SeedGenerator = new(sky.SeedService)
		walletManager = walletM
	})
	walletM.altManager = local.LoadAltcoinManager()
	walletM.updateTransactionAPI()
	walletM.updateSigner()
	walletM.updateWalletEnvs()
	for walletM.WalletEnv == nil {
		walletM.updateWalletEnvs()
	}
}

func (walletM *WalletManager) updateAll() {
	logWalletManager.Debug("Updating Wallet manager")
	walletM.altManager = local.LoadAltcoinManager()
	walletM.updateTransactionAPI()
	walletM.updateSigner()
	walletM.updateWalletEnvs()
	skycoin.UpdateAltcoin()
}

func GetWalletEnv() core.WalletEnv {
	logWalletManager.Info("Getting WalletEnv")
	wm := GetWalletManager()
	if wm == nil {
		return nil
	}

	return wm.WalletEnv
}

func GetWalletManager() *WalletManager {
	return walletManager
}

func (walletM *WalletManager) getDefaultWalletType() string {
	return walletM.WalletEnv.GetWalletSet().DefaultWalletType()
}

func (walletM *WalletManager) getAvailableWalletTypes() []string {
	return walletM.WalletEnv.GetWalletSet().SupportedWalletTypes()
}

func (walletM *WalletManager) updateSigner() {
	logWalletManager.Info("Updating Signers")
	signers := make([]core.BlockchainSignService, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		sing, err := plug.LoadSignService()
		if err != nil {
			logWalletManager.WithError(err).Errorf("Error loading signer from %s plugin", plug.GetName())
		}
		signers = append(signers, sing)
	}

	walletM.signer = signers[0]
}

func (walletM *WalletManager) updateTransactionAPI() {
	logWalletManager.Info("Updating TransactionAPI")
	txnAPIS := make([]core.BlockchainTransactionAPI, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		txnAPI, err := plug.LoadTransactionAPI("MainNet")
		if err != nil {
			logWalletManager.WithError(err).Errorf("Error loading transaction API from %s plugin", plug.GetName())
		}
		txnAPIS = append(txnAPIS, txnAPI)
	}

	walletM.transactionAPI = txnAPIS[0]
}

func (walletM *WalletManager) updateWalletEnvs() {
	logWalletManager.Info("Updating WalletEnvs")
	var walletsEnvs = make([]core.WalletEnv, 0)

	for _, plug := range walletM.altManager.ListRegisteredPlugins() {
		walletsEnvs = append(walletsEnvs, plug.LoadWalletEnvs()...)
	}
	if len(walletsEnvs) == 0 {
		logWalletManager.Error("Error loading wallet envs")
		return
	}
	walletM.WalletEnv = walletsEnvs[0]
}

func (walletM *WalletManager) broadcastTxn(txn *transactions.TransactionDetails) bool {
	logWalletManager.Info("Broadcasting transaction")
	altManager := local.LoadAltcoinManager()

	// instantiate plugin implementing for main coin of transaction represented by ticked (the first asset)
	plug, _ := altManager.LookupAltcoinPlugin(txn.Txn.SupportedAssets()[0])
	pex, err := plug.LoadPEX("MainNet")

	if err != nil {
		logWalletManager.WithError(err).Warn("Error loading PEX")
		return false
	}

	isSigned, err := txn.Txn.IsFullySigned()
	if err != nil {
		logWalletManager.WithError(err).Warn("Error checking if transaction if fully signed")
		return false
	}
	if !isSigned {
		logWalletManager.Warn("Transaction is not fully signed")
		return false
	}
	err = pex.BroadcastTxn(txn.Txn)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error broadcasting transaction")
		return false
	}
	logWalletManager.Info("Transaction Injected")
	return true
}

func (walletM *WalletManager) sendFromOutputs(wltIds []string, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *transactions.TransactionDetails {
	logWalletManager.Info("Creating transaction")
	wltCache := make(map[string]core.Wallet, 0)
	wlts := make([]core.Wallet, 0)
	for _, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to create transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wlts = append(wlts, wlt)
	}

	outputsFrom := make([]core.TransactionOutput, 0)
	for _, outAddr := range from {
		out := util.NewGenericOutput(nil, outAddr)
		outputsFrom = append(outputsFrom, &out)
	}
	outputsTo := make([]core.TransactionOutput, 0)
	for i := 0; i < len(addrTo); i++ {
		ch := ""
		if !automaticCoinHours {
			ch = coinHoursTo[i]
		}
		addr := util.NewGenericAddress(addrTo[i])
		out := util.NewGenericOutput(&addr, "")
		// FIXME: Remove explicit reference to Skycoin
		err := out.PushCoins(sky.Sky, skyTo[i])
		if err != nil {
			logWalletManager.WithError(err).Warn("Error parsing value for %s", sky.Sky)
			return nil
		}
		// FIXME: Remove explicit reference to Skycoin
		err = out.PushCoins(sky.CoinHour, ch)
		if err != nil {
			logWalletManager.WithError(err).Warn("Error parsing value for %s", sky.Sky)
			return nil
		}
		outputsTo = append(outputsTo, &out)
	}
	changeAddr := util.NewGenericAddress(change)
	opt := util.NewKeyValueMap()
	opt.SetValue("BurnFactor", burnFactor)
	if automaticCoinHours {
		opt.SetValue("CoinHoursSelectionType", "auto")
	} else {
		opt.SetValue("CoinHoursSelectionType", "manual")
	}
	var txn core.Transaction
	var err error
	if len(wltCache) > 1 {
		walletsOutputs := make([]core.WalletOutput, 0)
		for i, wlt := range wlts {
			walletsOutputs = append(walletsOutputs, &util.SimpleWalletOutput{
				Wallet: wlt,
				UxOut:  outputsFrom[i],
			})
		}
		txn, err = walletM.transactionAPI.Spend(walletsOutputs, outputsTo, &changeAddr, opt)
	} else {
		txn, err = wlts[0].Spend(outputsFrom, outputsTo, &changeAddr, opt)
	}

	if err != nil {
		logWalletManager.WithError(err).Info("Error creating transaction")
		return nil
	}

	qTransaction, err := transactions.NewTransactionDetailFromCoreTransaction(txn, transactions.TransactionTypeGeneric)
	if err != nil {
		logWalletManager.WithError(err).Info("Error converting transaction")
		return nil
	}
	return qTransaction
}

func (walletM *WalletManager) sendFromAddresses(wltIds []string, from, addrTo, skyTo, coinHoursTo []string, change string, automaticCoinHours bool, burnFactor string) *transactions.TransactionDetails {
	wltCache := make(map[string]core.Wallet, 0)
	wlts := make([]core.Wallet, 0)
	for _, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to create transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wlts = append(wlts, wlt)
	}

	addrsFrom := make([]core.Address, 0)
	for _, addr := range from {

		addrsFrom = append(addrsFrom, &util.GenericAddress{addr})
	}
	outputsTo := make([]core.TransactionOutput, 0)
	for i := 0; i < len(addrTo); i++ {
		ch := ""
		if !automaticCoinHours {
			ch = coinHoursTo[i]
		}
		addr := util.NewGenericAddress(addrTo[i])
		out := util.NewGenericOutput(&addr, "")
		// FIXME: Remove explicit reference to Skycoin
		err := out.PushCoins(sky.Sky, skyTo[i])
		if err != nil {
			logWalletManager.WithError(err).Warnf("Error parsing value for %s", sky.Sky)
			return nil
		}
		// FIXME: Remove explicit reference to Skycoin
		err = out.PushCoins(sky.CoinHour, ch)
		if err != nil {
			logWalletManager.WithError(err).Warnf("Error parsing value for %s", sky.Sky)
			return nil
		}
		outputsTo = append(outputsTo, &out)
	}
	changeAddr := &util.GenericAddress{change}

	opt := util.NewKeyValueMap()
	opt.SetValue("BurnFactor", burnFactor)
	if automaticCoinHours {
		opt.SetValue("CoinHoursSelectionType", "auto")
	} else {
		opt.SetValue("CoinHoursSelectionType", "manual")
	}
	var txn core.Transaction
	var err error
	if len(wltCache) > 1 {
		walletsAddresses := make([]core.WalletAddress, 0)
		for i, wlt := range wlts {
			walletsAddresses = append(walletsAddresses, &util.SimpleWalletAddress{
				Wallet: wlt,
				UxOut:  addrsFrom[i],
			})
		}
		txn, err = walletM.transactionAPI.SendFromAddress(walletsAddresses, outputsTo, changeAddr, opt)
	} else {
		txn, err = wlts[0].SendFromAddress(addrsFrom, outputsTo, changeAddr, opt)
	}

	if err != nil {
		logWalletManager.WithError(err).Info("Error creating transaction")
		return nil
	}

	qtxn, err := transactions.NewTransactionDetailFromCoreTransaction(txn, transactions.TransactionTypeGeneric)
	if err != nil {
		logWalletManager.WithError(err).Info("Error converting transaction")
		return nil
	}
	logWalletManager.Info("Transaction created")
	return qtxn

}

func (walletM *WalletManager) sendTo(wltId, destinationAddress, amount, currency string) *transactions.TransactionDetails {
	logWalletManager.Info("Creating Transaction")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
	addr := util.NewGenericAddress(destinationAddress)
	opt := util.GetOptionForCurrencyTxn(currency)

	if wlt == nil {
		logWalletManager.Warn("Couldn't load wallet to create transaction")
		return nil
	}
	txOut := util.NewGenericOutput(&addr, "")

	err := txOut.PushCoins(wlt.GetCryptoAccount().ListAssets()[0], amount)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error parsing value for %s", wlt.GetCryptoAccount().ListAssets()[0])
		return nil
	}
	txn, err := wlt.Transfer(&txOut, opt)
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't create transaction")
		return nil
	}
	qTxn, err := transactions.NewTransactionDetailFromCoreTransaction(txn, transactions.TransactionTypeSend)
	if err != nil {
		logWalletManager.WithError(err).Warn("Couldn't convert transaction")
		return nil
	}

	qTxn.SetAmount(amount)
	logWalletManager.Info("Transaction created")
	return qTxn
}

func (walletM *WalletManager) signTxn(wltIds, address []string, source string, tmpPwd interface{}, index []int, qTxn *transactions.TransactionDetails) *transactions.TransactionDetails {
	pwd, isPwdReader := tmpPwd.(core.PasswordReader)
	if !isPwdReader {
		return nil
	}
	logWalletManager.Info("Signing transaction")

	if len(wltIds) != len(address) {
		logWalletManager.Error("Wallets and addresses provided are incorrect")
		return nil
	}
	wltCache := make(map[string]core.Wallet)
	wltByAddr := make(map[string]core.Wallet)
	wlts := make([]core.Wallet, 0)

	for i, wltId := range wltIds {
		var wlt core.Wallet
		wlt, exist := wltCache[wltId]
		if !exist {
			wlt = walletM.WalletEnv.GetWalletSet().GetWallet(wltId)
			if wlt == nil {
				logWalletManager.Warn("Couldn't load wallet to Sign transaction")
				return nil
			}
			wltCache[wltId] = wlt
		}
		wltByAddr[address[i]] = wlt
		wlts = append(wlts, wlt)
	}

	var txn core.Transaction
	var err error

	if len(wltCache) > 1 {
		signDescriptors := make([]core.InputSignDescriptor, 0)
		for _, in := range qTxn.Txn.GetInputs() {
			sd := core.InputSignDescriptor{
				InputIndex: in.GetId(),
				SignerID:   core.UID(source),
				// Wallet:     wltByAddr[outAddr.String()],
			}
			signDescriptors = append(signDescriptors, sd)
		}
		txn, err = walletM.signer.Sign(qTxn.Txn, signDescriptors, pwd)
	} else {
		signer, err2 := util.LookupSignServiceForWallet(wlts[0], core.UID(source))
		if err2 != nil {
			logWalletManager.WithError(err).Warnf("No signer %s for wallet %v", source, wlts[0])
			return nil
		}
		signerUid, err := signer.GetSignerUID()
		if err != nil {
			logWalletManager.WithError(err).Errorln("unable to ger signer uuid")
			return nil
		}
		if wlts[0].GetId() == string(signerUid) {
			// NOTE the signer is the wallet it self
			signer = nil
		}
		txn, err = wlts[0].Sign(qTxn.Txn, signer, pwd, nil)
	}

	if err != nil {
		logWalletManager.WithError(err).Warn("Error signing txn")
		return nil
	}

	qTxn, err = transactions.NewTransactionDetailFromCoreTransaction(txn, transactions.TransactionTypeGeneric)
	if err != nil {
		logWalletManager.WithError(err).Warn("Error converting transaction")
		return nil
	}

	return qTxn

}

func (walletM *WalletManager) signAndBroadcastTxnAsync(wltIds, addresses []string, source string, bridgeForPassword *QBridge, index []int, qTxn *transactions.TransactionDetails) {
	channel := make(chan *transactions.TransactionDetails)

	go func() {
		var pwd core.PasswordReader = func(message string, ctx core.KeyValueStore) (string, error) {
			bridgeForPassword.BeginUse()
			defer bridgeForPassword.EndUse()
			bridgeForPassword.lock()
			suffix := ""
			v := ctx.GetValue(core.StrWalletLabel)
			if v == nil {
				v = ctx.GetValue(core.StrWalletName)
			}
			if v != nil {
				if str, isStr := v.(string); isStr {
					suffix = " for " + str
				}
			}
			bridgeForPassword.GetPassword(message + suffix)
			bridgeForPassword.lock()
			pass := bridgeForPassword.getResult()
			bridgeForPassword.unlock()
			return pass, nil
		}

		channel <- walletM.signTxn(wltIds, addresses, source, pwd, index, qTxn)
	}()

	go func() {
		txn := <-channel
		if txn != nil {
			walletM.broadcastTxn(txn)
		}
	}()
}

func (walletM *WalletManager) createEncryptedWallet(seed, label, wltType, password string, scanN int) *wallets.QWallet {
	logWalletManager.Info("Creating encrypted wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, wltType, true, pwd, scanN)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create encrypted wallet")
		return nil
	}

	logWalletManager.Info("Created encrypted wallet")
	qWallet := wallets.FromWalletToQWallet(wlt, true)
	walletM.wallets = append(walletM.wallets, qWallet)

	return qWallet
}

func (walletM *WalletManager) createUnencryptedWallet(seed, label, wltType string, scanN int) *wallets.QWallet {
	logWalletManager.Info("Creating unencrypted wallet")
	pwd := util.EmptyPassword
	wlt, err := walletM.WalletEnv.GetWalletSet().CreateWallet(label, seed, wltType, false, pwd, scanN)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create unencrypted wallet")
		return nil
	}
	logWalletManager.Info("Created unencrypted wallet")

	qWallet := wallets.FromWalletToQWallet(wlt, true)
	walletM.wallets = append(walletM.wallets, qWallet)

	return qWallet
}

func (walletM *WalletManager) getNewSeed(entropy int) string {
	logWalletManager.Info("Getting new seed")
	seed, err := walletM.SeedGenerator.GenerateMnemonic(entropy)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't get new seed")
		return ""
	}
	logWalletManager.Info("Generated new seed")
	return seed
}

func (walletM *WalletManager) verifySeed(seed string) int {
	logWalletManager.Info("Verifying seed")
	ok, err := walletM.SeedGenerator.VerifyMnemonic(seed)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't verify seed")
		return 0
	}
	logWalletManager.Info("Verified seed")
	if ok {
		return 1
	}
	return 0

}

func (walletM *WalletManager) encryptWallet(id, password string) int {
	logWalletManager.Info("Encrypting wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	walletM.WalletEnv.GetStorage().Encrypt(id, pwd)
	ret, err := walletM.WalletEnv.GetStorage().IsEncrypted(id)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't create encrypted wallets")
	}
	logWalletManager.Info("Wallet encrypted")
	loadWlt := func(id string, encrypted bool) {
		updtWltChan <- wallets.UpdateWallet{
			Wlt:   wallets.FromWalletToQWallet(walletM.WalletEnv.GetWalletSet().GetWallet(id), encrypted),
			Roles: []int{wallets.EncryptionEnabled},
		}
	}
	if ret {
		loadWlt(id, true)
		return 1
	}
	loadWlt(id, false)
	return 0
}

func (walletM *WalletManager) decryptWallet(id, password string) int {
	logWalletManager.Info("Decrypt wallet")
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	walletM.WalletEnv.GetStorage().Decrypt(id, pwd)
	ret, err := walletM.WalletEnv.GetStorage().IsEncrypted(id)
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't decrypt wallet")
	}
	logWalletManager.Info("Wallet decrypted")

	loadWlt := func(id string, encrypted bool) {
		updtWltChan <- wallets.UpdateWallet{
			Wlt:   wallets.FromWalletToQWallet(walletM.WalletEnv.GetWalletSet().GetWallet(id), encrypted),
			Roles: []int{wallets.EncryptionEnabled},
		}
	}
	if ret {
		loadWlt(id, true)
		return 1
	}
	loadWlt(id, false)
	return 0
}

func (walletM *WalletManager) newWalletAddress(id string, n int, password string) {
	logWalletManager.Info("Creating new wallet addresses")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(id)
	pwd := util.ConstantPassword(password)
	// NOTE: No easy way to get plain passwords in memory
	password = ""
	wltEntriesLen := 0
	it, err := wlt.GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Error("Couldn't load addresses")
		return
	}
	for it.Next() {
		wltEntriesLen++
	}
	wlt.GenAddresses(core.AccountAddress, uint32(wltEntriesLen), uint32(n), pwd)
}

func (walletM *WalletManager) getWallets() []*wallets.QWallet {
	logWalletManager.Info("Getting wallets")
	walletM.wallets = make([]*wallets.QWallet, 0)
	if walletM.WalletEnv == nil {
		walletM.updateWalletEnvs()
	}
	it := walletM.WalletEnv.GetWalletSet().ListWallets()

	if it == nil {
		logWalletManager.WithError(nil).Error("Couldn't load wallets")
		return walletM.wallets
	}

	for it.Next() {
		encrypted, err := walletM.WalletEnv.GetStorage().IsEncrypted(it.Value().GetId())
		if err != nil {
			logWalletManager.WithError(err).Error("Couldn't get wallets")
			return walletM.wallets
		}
		if encrypted {
			qw := wallets.FromWalletToQWallet(it.Value(), true)
			walletM.wallets = append(walletM.wallets, qw)
		} else {
			qw := wallets.FromWalletToQWallet(it.Value(), false)
			walletM.wallets = append(walletM.wallets, qw)
		}
	}

	logWalletManager.Info("Wallets obtained")
	return walletM.wallets
}

func (walletM *WalletManager) initWltModelAsync(model *wallets.WalletModel) {
	logWalletManager.Info("Init wallet model async")
	getPos := func(list []*wallets.QWallet, obj *wallets.QWallet) int {
		for e := range list {
			if list[e].FileName() == obj.FileName() {
				return e
			}
		}
		return -1
	}

	go func() {
		mutex := sync.Mutex{}
		for {
			select {
			case <-time.After(time.Duration(config.GetDataUpdateTime()) * time.Second):
				for _, qWlt := range model.Wallets() {
					if qWlt.EncryptionEnabled() == 0 {
						wallets.LoadCoinFtrFromWalletAsync(qWlt.GetCorWlt(), false, updtWltChan)
					} else {
						wallets.LoadCoinFtrFromWalletAsync(qWlt.GetCorWlt(), true, updtWltChan)
					}
				}
				break
			case updtWlt := <-updtWltChan:
				if len(updtWlt.Roles) == 0 {
					util2.Helper.RunInMain(func() {
						mutex.Lock()
						defer mutex.Unlock()
						model.AddWallet(updtWlt.Wlt)
					})
				} else {
					util2.Helper.RunInMain(func() {
						mutex.Lock()
						defer mutex.Unlock()
						model.EditWallet(getPos(model.Wallets(), updtWlt.Wlt), updtWlt.Wlt, updtWlt.Roles)
					})
				}
			}
		}
	}()
}

func (walletM *WalletManager) loadWallets(model *wallets.WalletModel) {
	logWalletManager.Info("Load qWallets async")
	go func() {
		for _, qWlt := range walletM.getWallets() {
			updtWltChan <- wallets.UpdateWallet{
				Wlt:   qWlt,
				Roles: []int{},
			}
		}
	}()
}

func (walletM *WalletManager) editWalletLbl(id, label string) {
	logWalletManager.Info("Editing wallet")
	wlt := walletM.WalletEnv.GetWalletSet().GetWallet(id)
	wlt.SetLabel(label)

	updtWltChan <- wallets.UpdateWallet{
		Wlt:   wallets.FromWalletToQWallet(wlt, true),
		Roles: []int{wallets.Name},
	}
}

func (walletM *WalletManager) getCurrencyList() []string {
	logWalletManager.Info("Obtaining list of currency")

	var currencyList = make([]string, 0)
	exist := func(currency string) bool {
		for e := range currencyList {
			if currencyList[e] == currency {
				return true
			}
		}
		return false
	}

	for e := range walletM.wallets {
		if !exist(walletM.wallets[e].Currency()) {
			currencyList = append(currencyList, walletM.wallets[e].Currency())
		}
	}
	sort.Strings(currencyList)
	return currencyList
}

func (walletM *WalletManager) loadAddressModelByWallet(wltName string, model *address.ModelAddress) {
	addrIter, err := walletM.WalletEnv.GetWalletSet().GetWallet(wltName).GetLoadedAddresses()
	if err != nil {
		logWalletManager.WithError(err).Warnf("Couldn't get address iterator for wallet %s", wltName)
		return
	}

	addrList := make([]*address.QAddress, 0)
	for addrIter.Next() {
		addrList = append(addrList, address.FromCorAddrToQAddr(addrIter.Value(), wltName))
	}
	model.LoadModelAsync(addrList)
}

func (walletM *WalletManager) loadAddressForAllWallets(model *address.ModelAddress) {
	wltIter := walletM.WalletEnv.GetWalletSet().ListWallets()
	if wltIter == nil {
		panic("implement me")
	}

	addrList := make([]*address.QAddress, 0)

	for wltIter.Next() {
		addrIter, err := wltIter.Value().GetLoadedAddresses()
		if err != nil {
			panic("implement me")
		}
		for addrIter.Next() {
			addrList = append(addrList, address.FromCorAddrToQAddr(addrIter.Value(), wltIter.Value().GetLabel()))
		}
	}
	model.LoadModel(addrList)
}

func (walletM *WalletManager) addWalletAsync(wallet *wallets.QWallet) {
	logWalletManager.Info("Add wallet async ", wallet.GetCorWlt().GetLabel())
	updtWltChan <- wallets.UpdateWallet{
		Wlt:   wallet,
		Roles: []int{},
	}
}
