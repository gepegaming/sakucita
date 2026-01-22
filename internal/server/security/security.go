package security

import (
	"crypto/rsa"

	"sakucita/pkg/config"

	"github.com/rs/zerolog"
)

type Security struct {
	config    config.App
	log       zerolog.Logger
	activeKID string
	rsaKeys   map[string]*RSAKeys
}

type RSAKeys struct {
	private *rsa.PrivateKey
	public  *rsa.PublicKey
}

// ! inget usahain security itu primitive agar tidak menjadi god object yang semuanya bergantung kesini

func NewSecurity(cfg config.App, log zerolog.Logger) *Security {
	return &Security{
		config: cfg,
		log:    log,
	}
}
