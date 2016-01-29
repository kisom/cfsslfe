// +build !nopkcs11

package pkcs11

import (
	"io/ioutil"

	"github.com/kisom/cfssl/config"
	"github.com/kisom/cfssl/crypto/pkcs11key"
	"github.com/kisom/cfssl/errors"
	"github.com/kisom/cfssl/helpers"
	"github.com/kisom/cfssl/log"
	"github.com/kisom/cfssl/signer"
	"github.com/kisom/cfssl/signer/local"
)

// Enabled is set to true if PKCS #11 support is present.
const Enabled = true

// New returns a new PKCS #11 signer.
func New(caCertFile string, policy *config.Signing, cfg *Config) (signer.Signer, error) {
	if cfg == nil {
		return nil, errors.New(errors.PrivateKeyError, errors.ReadFailed)
	}

	log.Debugf("Loading PKCS #11 module %s", cfg.Module)
	certData, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, errors.New(errors.PrivateKeyError, errors.ReadFailed)
	}

	cert, err := helpers.ParseCertificatePEM(certData)
	if err != nil {
		return nil, err
	}

	priv, err := pkcs11key.New(cfg.Module, cfg.Token, cfg.PIN, cfg.Label)
	if err != nil {
		return nil, errors.New(errors.PrivateKeyError, errors.ReadFailed)
	}
	sigAlgo := signer.DefaultSigAlgo(priv)

	return local.NewSigner(priv, cert, sigAlgo, policy)
}
