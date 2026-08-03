package deploy

import (
	"fmt"
	"path/filepath"
	"strings"

	hyperproviders "github.com/herdworks/hyperproviders-go"
	envsecrets "github.com/herdworks/hyperproviders-go/providers/env"
	filesecrets "github.com/herdworks/hyperproviders-go/providers/files"
	_ "github.com/herdworks/hyperproviders-go/providers/git"
	_ "github.com/herdworks/hyperproviders-go/providers/s3"
	hypertexting "github.com/herdworks/hypertexting-go"
)

const defaultEnvironment = "default"

func deploymentEnvironment(website *hypertexting.Website, environmentName, providerNames string) (*hyperproviders.Environment, string, error) {
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

func deploymentKind(incremental bool) hyperproviders.DeploymentKind {
	if incremental {
		return hyperproviders.IncrementalDeployment
	}
	return hyperproviders.CompleteDeployment
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

func DeploymentSecrets(website *hypertexting.Website) hyperproviders.SecretsProviders {
	fileSecrets := filesecrets.New()
	if website != nil && website.URI != nil && website.URI.Path != "" {
		fileSecrets = filesecrets.NewWithBaseDirectory(filepath.Dir(website.URI.Path))
	}

	return hyperproviders.SecretsProviders{envsecrets.New(), fileSecrets}
}
