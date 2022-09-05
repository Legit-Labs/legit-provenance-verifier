package legitprovenance

import (
	"flag"
	"fmt"
	"strings"

	intoto "github.com/in-toto/in-toto-golang/in_toto"
)

type ProvenanceChecks struct {
	RepoUrl   string
	Branch    string
	Tag       string
	BuilderId string
}

const (
	defaultBuilderID = `https://github.com/legit-labs/legit-provenance-generator/.github/workflows/legit_provenance_generator.yml@refs/tags/v0.1.0`
	branchRefPrefix  = "refs/heads/"
	tagRefPrefix     = "refs/tags/"
)

func (pc *ProvenanceChecks) Flags() {
	flag.StringVar(&pc.RepoUrl, "repo-url", "", "The source repository url (default: no check)")
	flag.StringVar(&pc.Branch, "branch", "", "The source branch (default: no check)")
	flag.StringVar(&pc.Tag, "tag", "", "The tag of the commit (default: no check)")
	flag.StringVar(&pc.BuilderId, "builder-id", defaultBuilderID, "The builder ID of the provenance generator (default: Legit's provenance generator)")
}

func (pc *ProvenanceChecks) Verify(provenance intoto.ProvenanceStatement) error {
	if err := pc.verifyRepo(provenance); err != nil {
		return err
	}

	if err := pc.verifyBuilderID(provenance); err != nil {
		return err
	}

	if err := pc.verifyBranch(provenance); err != nil {
		return err
	}

	if err := pc.verifyTag(provenance); err != nil {
		return err
	}

	return nil
}

func (pc *ProvenanceChecks) verifyTag(provenance intoto.ProvenanceStatement) error {
	if pc.Tag == "" {
		return nil
	}

	tag, err := pullTag(provenance)
	if err != nil {
		return err
	}

	if tag != pc.Tag {
		return fmt.Errorf("expected tag %v does not match actual: %v", pc.Tag, tag)
	}

	return nil
}

func (pc *ProvenanceChecks) verifyBranch(provenance intoto.ProvenanceStatement) error {
	if pc.Branch == "" {
		return nil
	}

	branch, err := pullBranch(provenance)
	if err != nil {
		return err
	}

	if branch != pc.Branch {
		return fmt.Errorf("expected branch %v does not match actual: %v", pc.Branch, branch)
	}

	return nil
}

func (pc *ProvenanceChecks) verifyBuilderID(provenance intoto.ProvenanceStatement) error {
	if pc.BuilderId == "" {
		return nil
	}

	builderID, err := pullBuilderID(provenance)
	if err != nil {
		return err
	}

	if builderID != pc.BuilderId {
		return fmt.Errorf("expected builder ID %v does not match actual: %v", pc.BuilderId, builderID)
	}

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

func pullBuilderID(provenance intoto.ProvenanceStatement) (string, error) {
	return provenance.Predicate.Builder.ID, nil
}

func pullBranch(provenance intoto.ProvenanceStatement) (string, error) {
	event, err := pullGithubEventPayload(provenance)
	if err != nil {
		return "", err
	}

	branch, ok := event["base_ref"]
	if !ok {
		return "", fmt.Errorf("failed to pull base ref (branch)")
	}

	asString, ok := branch.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type of branch: %T", branch)
	}

	clean := strings.TrimPrefix(asString, branchRefPrefix)

	return clean, nil
}

func pullTag(provenance intoto.ProvenanceStatement) (string, error) {
	event, err := pullGithubEventPayload(provenance)
	if err != nil {
		return "", err
	}

	tag, ok := event["ref"]
	if !ok {
		return "", fmt.Errorf("failed to pull ref (tag)")
	}

	asString, ok := tag.(string)
	if !ok {
		return "", fmt.Errorf("unexpected type of tag: %T", tag)
	}

	clean := strings.TrimPrefix(asString, tagRefPrefix)

	return clean, nil
}
