package legit_provenance

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	verifyattestation "github.com/legit-labs/legit-verify-attestation/verify-attestation"
)

func Verify(attestationPath string, keyPath string, digest string, checks ProvenanceChecks) error {
	attestation, err := os.ReadFile(attestationPath)
	if err != nil {
		return fmt.Errorf("failed to open attestation at %v: %v", attestationPath, err)
	}

	payload, err := verifyattestation.VerifiedPayload(context.Background(), keyPath, attestation)
	if err != nil {
		return fmt.Errorf("attestation verification failed: %v", err)
	}

	var statement intoto.ProvenanceStatement
	err = json.Unmarshal(payload, &statement)
	if err != nil {
		return fmt.Errorf("failed to unmarshal predicate: %v", err)
	}

	statementDigest := statement.Subject[0].Digest["sha256"]
	if statementDigest != digest {
		return fmt.Errorf("expected digest %v does not match actual: %v", digest, statementDigest)
	}

	if err = checks.Verify(&statement); err != nil {
		return fmt.Errorf("provenance verification failed: %v", err)
	}

	return nil
}
