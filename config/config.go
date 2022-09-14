package config

import (
	"os"

	crypto "github.com/portalnesia/go-crypto"
)

var (
	NODE_ENV     string
	IsProduction bool = true
	Crypto       crypto.CryptoKey
	Prefix       string
)

func SetupConfig() {
	NODE_ENV = os.Getenv("NODE_ENV")
	Crypto = crypto.PortalnesiaCrypto(os.Getenv("SECRET"))
	IsProduction = (NODE_ENV == "production")
	Prefix = os.Getenv("DB_PREFIX")
}
