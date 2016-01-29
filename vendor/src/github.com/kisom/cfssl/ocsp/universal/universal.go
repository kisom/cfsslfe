package universal

import (
	"github.com/kisom/cfssl/ocsp"
	ocspConfig "github.com/kisom/cfssl/ocsp/config"
	"github.com/kisom/cfssl/ocsp/pkcs11"
)

// NewSignerFromConfig generates a new OCSP signer from a config object.
func NewSignerFromConfig(cfg ocspConfig.Config) (ocsp.Signer, error) {
	if cfg.PKCS11.Module != "" {
		return pkcs11.NewPKCS11Signer(cfg)
	}
	return ocsp.NewSignerFromFile(cfg.CACertFile, cfg.ResponderCertFile,
		cfg.KeyFile, cfg.Interval)
}

