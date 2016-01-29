package revoke

import (
	"testing"

	"database/sql"
	"time"

	"github.com/kisom/cfssl/certdb"
	"github.com/kisom/cfssl/certdb/testdb"
	"github.com/kisom/cfssl/cli"
	"golang.org/x/crypto/ocsp"
)

func prepDB() (db *sql.DB, err error) {
	db = testdb.SQLiteDB("../../certdb/testdb/certstore_development.db")
	expirationTime := time.Now().AddDate(1, 0, 0)
	var cert = &certdb.CertificateRecord{
		Serial: "1",
		Expiry: expirationTime,
		PEM:    "unexpired cert",
	}

	err = certdb.InsertCertificate(db, cert)
	if err != nil {
		return nil, err
	}

	return
}

func TestRevokeMain(t *testing.T) {
	db, err := prepDB()
	if err != nil {
		t.Fatal(err)
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1", DBConfigFile: "../testdata/db-config.json"})
	if err != nil {
		t.Fatal(err)
	}

	var crs *certdb.CertificateRecord
	crs, err = certdb.GetCertificate(db, "1")
	if err != nil {
		t.Fatal("Failed to get certificate")
	}

	if crs.Status != "revoked" {
		t.Fatal("Certificate not marked revoked after we revoked it")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1", Reason: "2", DBConfigFile: "../testdata/db-config.json"})
	if err != nil {
		t.Fatal(err)
	}

	crs, err = certdb.GetCertificate(db, "1")
	if err != nil {
		t.Fatal("Failed to get certificate")
	}

	if crs.Reason != 2 {
		t.Fatal("Certificate revocation reason incorrect")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1", Reason: "Superseded", DBConfigFile: "../testdata/db-config.json"})
	if err != nil {
		t.Fatal(err)
	}

	crs, err = certdb.GetCertificate(db, "1")
	if err != nil {
		t.Fatal("Failed to get certificate")
	}

	if crs.Reason != ocsp.Superseded {
		t.Fatal("Certificate revocation reason incorrect")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1", Reason: "invalid_reason", DBConfigFile: "../testdata/db-config.json"})
	if err == nil {
		t.Fatal("Expected error from invalid reason")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1", Reason: "999", DBConfigFile: "../testdata/db-config.json"})
	if err == nil {
		t.Fatal("Expected error from invalid reason")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "2", DBConfigFile: "../testdata/db-config.json"})
	if err == nil {
		t.Fatal("Expected error from unrecognized serial number")
	}

	err = revokeMain([]string{}, cli.Config{DBConfigFile: "../testdata/db-config.json"})
	if err == nil {
		t.Fatal("Expected error from missing serial number")
	}

	err = revokeMain([]string{}, cli.Config{Serial: "1"})
	if err == nil {
		t.Fatal("Expected error from missing db config")
	}
}
