package models

func init() {
	WalletModel_QmlRegisterType2("WalletsManager", 1, 0, "WalletModel")
	QWallet_QmlRegisterType2("WalletsManager", 1, 0, "QWallet")
	AddressesModel_QmlRegisterType2("WalletsManager", 1, 0, "AddressModel")
	QAddress_QmlRegisterType2("WalletsManager", 1, 0, "QAddress")
	WalletManager_QmlRegisterType2("WalletsManager", 1, 0, "WalletManager")
	ConfigManager_QmlRegisterType2("Config", 1, 0, "ConfigManager")
	WalletSource_QmlRegisterType2("Config", 1, 0, "WalletSource")
}
