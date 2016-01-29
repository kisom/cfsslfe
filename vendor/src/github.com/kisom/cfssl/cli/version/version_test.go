package version

import (
	"testing"

	"github.com/kisom/cfssl/cli"
)

func TestVersionString(t *testing.T) {
	version := versionString()
	if version != "1.1.0" {
		t.Fatal("verstion string is not returned correctly")
	}
}

func TestVersionMain(t *testing.T) {
	args := []string{"cfssl", "version"}
	err := versionMain(args, cli.Config{})
	if err != nil {
		t.Fatal("version main failed")
	}
}
