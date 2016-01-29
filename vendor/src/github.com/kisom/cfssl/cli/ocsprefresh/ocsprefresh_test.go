package ocsprefresh

import (
	"testing"

	"time"

	"github.com/kisom/cfssl/certdb"
	"github.com/kisom/cfssl/certdb/testdb"
	"github.com/kisom/cfssl/cli"
	"golang.org/x/crypto/ocsp"
	"io/ioutil"
	"github.com/kisom/cfssl/helpers"
)

func TestOCSPRefreshMain(t *testing.T) {
	db := testdb.SQLiteDB("../../certdb/testdb/certstore_development.db")

	certPEM, err := ioutil.ReadFile("../../ocsp/testdata/cert.pem")
	if err != nil {
		t.Fatal(err)
	}

	expirationTime := time.Now().AddDate(1, 0, 0)
	var cert = &certdb.CertificateRecord{
		Serial: "1333308112180215502", // from cert.pem
		Expiry: expirationTime,
		PEM:    string(certPEM),
		Status: "good",
	}

	err = certdb.InsertCertificate(db, cert)
	if err != nil {
		t.Fatal(err)
	}

	err = ocsprefreshMain([]string{}, cli.Config{
		CAFile: "../../ocsp/testdata/ca.pem",
		ResponderFile:"../../ocsp/testdata/server.crt",
		ResponderKeyFile:"../../ocsp/testdata/server.key",
		DBConfigFile: "../testdata/db-config.json",
		Interval:helpers.OneDay,
	})

	if err != nil {
		t.Fatal(err)
	}

	var records []*certdb.OCSPRecord
	records, err = certdb.GetUnexpiredOCSPs(db)
	if err != nil {
		t.Fatal("Failed to get OCSP responses")
	}

	if len(records) != 1 {
		t.Fatal("Expected one OCSP response")
	}

	var resp *ocsp.Response
	resp, err = ocsp.ParseResponse([]byte(records[0].Body), nil)
	if err != nil {
		t.Fatal("Failed to parse OCSP response")
	}
	if resp.Status != ocsp.Good {
		t.Fatal("Expected cert status 'good'")
	}

	err = certdb.RevokeCertificate(db, cert.Serial, ocsp.KeyCompromise)
	if err != nil {
		t.Fatal("Failed to revoke certificate")
	}

	err = ocsprefreshMain([]string{}, cli.Config{
		CAFile: "../../ocsp/testdata/ca.pem",
		ResponderFile:"../../ocsp/testdata/server.crt",
		ResponderKeyFile:"../../ocsp/testdata/server.key",
		DBConfigFile: "../testdata/db-config.json",
		Interval:helpers.OneDay,
	})

	if err != nil {
		t.Fatal(err)
	}

	records, err = certdb.GetUnexpiredOCSPs(db)
	if err != nil {
		t.Fatal("Failed to get OCSP responses")
	}

	if len(records) != 1 {
		t.Fatal("Expected one OCSP response")
	}

	resp, err = ocsp.ParseResponse([]byte(records[0].Body), nil)
	if err != nil {
		t.Fatal("Failed to parse OCSP response")
	}
	if resp.Status != ocsp.Revoked {
		t.Fatal("Expected cert status 'revoked'")
	}
}
