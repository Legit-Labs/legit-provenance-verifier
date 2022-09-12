package main

import (
	"flag"
	"log"

	legit_provenance "github.com/legit-labs/legit-provenance-verifier/legit-provenance"
)

var (
	keyPath         string
	attestationPath string
	digest          string
	checks          legit_provenance.ProvenanceChecks
)

func main() {
	flag.StringVar(&keyPath, "key", "", "The path of the public key")
	flag.StringVar(&attestationPath, "attestation", "", "The path of the attestation document")
	flag.StringVar(&digest, "digest", "", "The expected subject digest")

	checks.Flags()

	flag.Parse()

	if keyPath == "" {
		log.Panicf("please provide a public key path")
	} else if attestationPath == "" {
		log.Panicf("please provide an attestation path")
	}

	if err := legit_provenance.Verify(attestationPath, keyPath, digest, checks); err != nil {
		log.Panicf("verification failed: %v", err)
	}

	log.Printf("provenance verified successfully.")
}
