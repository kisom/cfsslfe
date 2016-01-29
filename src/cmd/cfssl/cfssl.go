/*
cfssl is the command line tool to issue/sign/bundle client certificate. It's
also a tool to start a HTTP server to handle web requests for signing, bundling
and verification.

Usage:
	cfssl command [-flags] arguments

	The commands are

	bundle	 create a certificate bundle
	sign	 signs a certificate signing request (CSR)
	serve	 starts a HTTP server handling sign and bundle requests
	version	 prints the current cfssl version
	genkey   generates a key and an associated CSR
	gencert  generates a key and a signed certificate
	selfsign generates a self-signed certificate

Use "cfssl [command] -help" to find out more about a command.
*/
package main

import (
	"flag"
	"os"

	"github.com/kisom/cfssl/cli"
	"github.com/kisom/cfssl/cli/bundle"
	"github.com/kisom/cfssl/cli/certinfo"
	"github.com/kisom/cfssl/cli/gencert"
	"github.com/kisom/cfssl/cli/gencrl"
	"github.com/kisom/cfssl/cli/genkey"
	"github.com/kisom/cfssl/cli/info"
	"github.com/kisom/cfssl/cli/ocspdump"
	"github.com/kisom/cfssl/cli/ocsprefresh"
	"github.com/kisom/cfssl/cli/ocspserve"
	"github.com/kisom/cfssl/cli/ocspsign"
	"github.com/kisom/cfssl/cli/printdefault"
	"github.com/kisom/cfssl/cli/revoke"
	"github.com/kisom/cfssl/cli/scan"
	"github.com/kisom/cfssl/cli/selfsign"
	"github.com/kisom/cfssl/cli/serve"
	"github.com/kisom/cfssl/cli/sign"
	"github.com/kisom/cfssl/cli/version"
	"github.com/kisom/cfssl/log"
)

// main defines the cfssl usage and registers all defined commands and flags.
func main() {
	// Add command names to cfssl usage
	flag.Usage = nil // this is set to nil for testabilty
	flag.IntVar(&log.Level, "loglevel", log.LevelInfo, "Log level")
	// Register commands.
	cmds := map[string]*cli.Command{
		"bundle":         bundle.Command,
		"certinfo":       certinfo.Command,
		"sign":           sign.Command,
		"serve":          serve.Command,
		"version":        version.Command,
		"genkey":         genkey.Command,
		"gencert":        gencert.Command,
		"gencrl":         gencrl.Command,
		"ocspdump":       ocspdump.Command,
		"ocsprefresh":    ocsprefresh.Command,
		"ocspsign":       ocspsign.Command,
		"ocspserve":      ocspserve.Command,
		"selfsign":       selfsign.Command,
		"scan":           scan.Command,
		"info":           info.Command,
		"print-defaults": printdefaults.Command,
		"revoke":         revoke.Command,
	}

	// If the CLI returns an error, exit with an appropriate status
	// code.
	err := cli.Start(cmds)
	if err != nil {
		os.Exit(1)
	}
}
