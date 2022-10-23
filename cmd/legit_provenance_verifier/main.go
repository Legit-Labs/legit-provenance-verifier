package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/legit-labs/legit-provenance-verifier/pkg/legit_provenance_verifier"
	"github.com/legit-labs/legit-registry-tools/pkg/legit_registry_tools"
)

var (
	keyPath          string
	attestationPath  string
	attestationStdin bool
	imageName        string
	digest           string
	checks           legit_provenance_verifier.ProvenanceChecks
)

func main() {
	flag.StringVar(&keyPath, "key", "", "The path of the public key")
	flag.StringVar(&attestationPath, "attestation-path", "", "The path of the attestation document")
	flag.BoolVar(&attestationStdin, "attestation-stdin", false, "Read the attestation from stdin (overwrites -attestation-path if provided)")
	flag.StringVar(&imageName, "image-name", "", "The name of the image (pass instead of providing attestation to pull from remote). note: needs to be logged in for private registries.")
	flag.StringVar(&digest, "digest", "", "The expected subject digest (without the sha256 prefix)")

	checks.Flags()

	flag.Parse()

	if keyPath == "" {
		log.Panicf("please provide a public key path")
	} else if imageName == "" && !attestationStdin && attestationPath == "" {
		log.Panicf("please provide an attestation path or image name (or set -attestation-stdin to read attestation from stdin)")
	} else if digest == "" {
		log.Panicf("please provide the expected digest")
	}

	var attestation []byte
	var err error
	ctx := context.Background()
	if imageName != "" {
		imageRef := legit_registry_tools.ImageRef{
			Name:   imageName,
			Digest: legit_registry_tools.DigestFromShaValue(digest),
		}
		err = legit_provenance_verifier.VerifyRemote(ctx, &imageRef, keyPath, checks)
	} else {
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

		err = legit_provenance_verifier.Verify(ctx, attestation, keyPath, digest, checks)
	}

	if err != nil {
		log.Panicf("verification failed: %v", err)
	}

	log.Printf("provenance verified successfully.")
}
