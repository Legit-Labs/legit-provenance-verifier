package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/legit-labs/legit-provenance-verifier/pkg/legit_provenance_verifier"
)

var (
	keyPath          string
	attestationPath  string
	attestationStdin bool
	digest           string
	checks           legit_provenance_verifier.ProvenanceChecks
)

func main() {
	flag.StringVar(&keyPath, "key", "", "The path of the public key")
	flag.StringVar(&attestationPath, "attestation-path", "", "The path of the attestation document")
	flag.BoolVar(&attestationStdin, "attestation-stdin", false, "Read the attestation from stdin (overwrites -attestation-path if provided)")
	flag.StringVar(&digest, "digest", "", "The expected subject digest")

	checks.Flags()

	flag.Parse()

	if keyPath == "" {
		log.Panicf("please provide a public key path")
	} else if !attestationStdin && attestationPath == "" {
		log.Panicf("please provide an attestation path (or set -attestation-stdin to read it from stdin)")
	} else if digest == "" {
		log.Panicf("please provide the expected digest")
	}

	var attestation []byte
	var err error
	if attestationStdin {
		if attestation, err = ioutil.ReadAll(os.Stdin); err != nil {
			log.Panicf("failed to read payload from stdin: %v", err)
		}
	} else {
		attestation, err = os.ReadFile(attestationPath)
		if err != nil {
			log.Panicf("failed to open payload at %v: %v", attestationPath, err)
		}
	}

	if err := legit_provenance_verifier.Verify(context.Background(), attestation, keyPath, digest, checks); err != nil {
		log.Panicf("verification failed: %v", err)
	}

	log.Printf("provenance verified successfully.")
}
