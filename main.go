package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	legitprovenance "github.com/legit-labs/legit-provenance-verifier/legit-provenance"
	verifyattestation "github.com/legit-labs/legit-verify-attestation/verify-attestation"
)

var (
	keyPath         string
	attestationPath string
	checks          legitprovenance.ProvenanceChecks
)

const (
	debugPayloadEnvKey = "DEBUG_PAYLAOD"
)

func main() {
	flag.StringVar(&keyPath, "key", "", "The path of the public key")
	flag.StringVar(&attestationPath, "attestation", "", "The path of the attestation document")
	checks.Flags()

	flag.Parse()

	if keyPath == "" {
		log.Panicf("please provide a public key path")
	} else if attestationPath == "" {
		log.Panicf("please provide an attestation path")
	}

	attestation, err := os.ReadFile(attestationPath)
	if err != nil {
		log.Panicf("failed to open attestation at %v: %v", attestationPath, err)
	}

	payload, err := verifyattestation.VerifiedPayload(context.Background(), keyPath, attestation)
	if err != nil {
		log.Panicf("attestation verification failed: %v", err)
	}

	if os.Getenv(debugPayloadEnvKey) == "1" {
		fmt.Printf("Payload:\n%v\n", string(payload))
	}

	var statement intoto.ProvenanceStatement
	err = json.Unmarshal(payload, &statement)
	if err != nil {
		log.Panicf("failed to unmarshal predicate: %v", err)
	}

	if err = checks.Verify(&statement); err != nil {
		log.Panicf("provenance verification failed: %v", err)
	}

	log.Printf("provenance verified successfully.")
}
