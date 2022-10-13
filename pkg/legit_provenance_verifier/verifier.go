package legit_provenance_verifier

import (
	"context"
	"fmt"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
	"github.com/legit-labs/legit-attestation/pkg/legit_verify_attestation"
)

var verifyPayload = legit_verify_attestation.VerifiedTypedPayload[intoto.ProvenanceStatement]

func Verify(ctx context.Context, attestation []byte, keyPath string, digest string, checks ProvenanceChecks) error {
	statement, err := verifyPayload(ctx, keyPath, attestation, digest)
	if err != nil {
		return fmt.Errorf("provenance payload verification failed: %v", err)
	}

	if err = checks.Verify(statement); err != nil {
		return fmt.Errorf("provenance checks failed: %v", err)
	}

	return nil
}
