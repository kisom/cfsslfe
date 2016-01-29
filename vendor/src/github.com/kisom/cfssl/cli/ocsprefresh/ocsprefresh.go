// Package ocsprefresh implements the ocsprefresh command.
package ocsprefresh

import (
	"database/sql"
	"time"

	"github.com/kisom/cfssl/certdb"
	"github.com/kisom/cfssl/cli"
	"github.com/kisom/cfssl/helpers"
	"github.com/kisom/cfssl/log"
	"github.com/kisom/cfssl/ocsp"
)

// Usage text of 'cfssl ocsprefresh'
var ocsprefreshUsageText = `cfssl ocsprefresh -- refreshes the ocsp_responses table
with new OCSP responses for all known unexpired certificates

Usage of ocsprefresh:
        cfssl ocsprefresh -db-config db-config -ca cert -responder cert -responder-key key [-interval 96h]

Flags:
`

// Flags of 'cfssl ocsprefresh'
var ocsprefreshFlags = []string{"ca", "responder", "responder-key", "db-config", "interval"}

// ocsprefreshMain is the main CLI of OCSP refresh functionality.
func ocsprefreshMain(args []string, c cli.Config) (err error) {
	if c.DBConfigFile == "" {
		log.Error("need DB config file (provide with -db-config)")
		return
	}

	if c.ResponderFile == "" {
		log.Error("need responder certificate (provide with -responder)")
		return
	}

	if c.ResponderKeyFile == "" {
		log.Error("need responder key (provide with -responder-key)")
		return
	}

	if c.CAFile == "" {
		log.Error("need CA certificate (provide with -ca)")
		return
	}

	s, err := SignerFromConfig(c)
	if err != nil {
		log.Critical("Unable to create OCSP signer: ", err)
		return err
	}

	var db *sql.DB
	db, err = certdb.DBFromConfig(c.DBConfigFile)
	if err != nil {
		return err
	}

	var certs []*certdb.CertificateRecord
	certs, err = certdb.GetUnexpiredCertificates(db)
	if err != nil {
		return err
	}

	// Set an expiry timestamp for all certificates refreshed in this batch
	ocspExpiry := time.Now().Add(c.Interval)
	for _, certRecord := range certs {
		cert, err := helpers.ParseCertificatePEM([]byte(certRecord.PEM))
		if err != nil {
			log.Critical("Unable to parse certificate: ", err)
			return err
		}

		req := ocsp.SignRequest{
			Certificate: cert,
			Status:      certRecord.Status,
		}

		if certRecord.Status == "revoked" {
			req.Reason = int(certRecord.Reason)
			req.RevokedAt = certRecord.RevokedAt
		}

		resp, err := s.Sign(req)
		if err != nil {
			log.Critical("Unable to sign OCSP response: ", err)
			return err
		}

		err = certdb.UpsertOCSP(db, cert.SerialNumber.String(), string(resp), ocspExpiry)
		if err != nil {
			log.Critical("Unable to save OCSP response: ", err)
			return err
		}
	}

	return nil
}

// SignerFromConfig creates a signer from a cli.Config as a helper for cli and serve
func SignerFromConfig(c cli.Config) (ocsp.Signer, error) {
	//if this is called from serve then we need to use the specific responder key file
	//fallback to key for backwards-compatibility
	k := c.ResponderKeyFile
	if k == "" {
		k = c.KeyFile
	}
	return ocsp.NewSignerFromFile(c.CAFile, c.ResponderFile, k, time.Duration(c.Interval))
}

// Command assembles the definition of Command 'ocsprefresh'
var Command = &cli.Command{UsageText: ocsprefreshUsageText, Flags: ocsprefreshFlags, Main: ocsprefreshMain}
