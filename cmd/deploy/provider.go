package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/herdworks/htup"
	envsecrets "github.com/herdworks/htup/providers/env"
	filesecrets "github.com/herdworks/htup/providers/files"
	_ "github.com/herdworks/htup/providers/git"
	_ "github.com/herdworks/htup/providers/s3"
	"github.com/herdworks/htxt"
)

const defaultEnvironment = "default"

func deploymentEnvironment(website *htxt.Website, environmentName, providerNames string) (*htup.Environment, string, error) {
	if environmentName != "" && providerNames != "" {
		return nil, "", fmt.Errorf("--providers and --environment are mutually exclusive")
	}
	if providerNames != "" {
		names := parseProviderNames(providerNames)
		env, err := website.DeploymentProviders(names)
		if err != nil {
			return nil, "", err
		}
		return &env, "", nil
	}
	if environmentName == "" {
		environmentName = defaultEnvironment
		fmt.Printf("No deployment target specified; using default environment %q\n", environmentName)
	}
	env, err := website.DeploymentEnvironment(environmentName)
	if err != nil {
		return nil, "", err
	}
	return env, environmentName, nil
}

func deploymentKind(incremental bool) htup.DeploymentKind {
	if incremental {
		return htup.IncrementalDeployment
	}
	return htup.CompleteDeployment
}

func parseProviderNames(value string) []string {
	parts := strings.Split(value, ",")
	names := make([]string, 0, len(parts))
	for _, part := range parts {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		names = append(names, name)
	}
	return names
}

func DeploymentSecrets(website *htxt.Website) htup.SecretsProviders {
	fileSecrets := filesecrets.New()
	if website != nil && website.URI != nil && website.URI.Path != "" {
		fileSecrets = filesecrets.NewWithBaseDirectory(filepath.Dir(website.URI.Path))
	}

	return htup.SecretsProviders{envsecrets.New(), fileSecrets}
}
