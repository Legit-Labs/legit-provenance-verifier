package legitprovenance

import (
	"flag"
	"fmt"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

type ProvenanceChecks struct {
	RepoUrl   string
	Branch    string
	BuilderId string
}

const (
	defaultBuilderID = `https://github.com/legit-labs/legit-provenance-generator/.github/workflows/legit_provenance_generator.yml@refs/tags/v0.1.0`
)

func (pc *ProvenanceChecks) Flags() {
	flag.StringVar(&pc.RepoUrl, "repo-url", "", "The source repository url (default: no check)")
	flag.StringVar(&pc.Branch, "branch", "", "The source branch (default: no check)")
	flag.StringVar(&pc.BuilderId, "builder-id", defaultBuilderID, "The builder ID of the provenance generator (default: Legit's provenance generator)")
}

func (pc *ProvenanceChecks) Verify(provenance intoto.ProvenanceStatement) error {
	err := pc.verifyRepo(provenance)
	if err != nil {
		return err
	}

	// TODO verify branch
	// TODO verify builder id

	return nil
}

func (pc *ProvenanceChecks) verifyRepo(provenance intoto.ProvenanceStatement) error {
	if pc.RepoUrl == "" {
		return nil
	}

	repo, err := pullProvenanceRepoUrl(provenance)
	if err != nil {
		return err
	}

	if repo != pc.RepoUrl {
		return fmt.Errorf("expected repo %v does not match actual: %v", pc.RepoUrl, repo)
	}

	return nil
}

func pullInvocationEnv(provenance intoto.ProvenanceStatement) (map[string]interface{}, error) {
	env, ok := provenance.Predicate.Invocation.Environment.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to pull environment info")
	}

	return env, nil
}

func pullGithubEventPayload(provenance intoto.ProvenanceStatement) (map[string]interface{}, error) {
	env, err := pullInvocationEnv(provenance)
	if err != nil {
		return nil, err
	}

	event, ok := env["github_event_payload"]
	if !ok {
		return nil, fmt.Errorf("failed to pull github event payload")
	}

	asMap, ok := event.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type of github event payload: %T", event)
	}

	return asMap, nil
}

func pullProvenanceRepoInfo(provenance intoto.ProvenanceStatement) (map[string]interface{}, error) {
	event, err := pullGithubEventPayload(provenance)
	if err != nil {
		return nil, err
	}

	repo, ok := event["repository"]
	if !ok {
		return nil, fmt.Errorf("failed to pull repository info")
	}

	asMap, ok := repo.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type of repository info: %T", repo)
	}

	return asMap, nil
}

func pullProvenanceRepoUrl(provenance intoto.ProvenanceStatement) (string, error) {
	repo, err := pullProvenanceRepoInfo(provenance)
	if err != nil {
		return "", err
	}

	url, ok := repo["url"]
	if !ok {
		return "", fmt.Errorf("failed to pull repository url")
	}

	asString, ok := url.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type of repository url: %T", url)
	}

	return asString, nil
}
