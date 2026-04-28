package deploy

import (
	hyperproviders "github.com/herdworks/hyperproviders-go"
	envsecrets "github.com/herdworks/hyperproviders-go/providers/env"
	_ "github.com/herdworks/hyperproviders-go/providers/git"
	_ "github.com/herdworks/hyperproviders-go/providers/s3"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func deploymentEnvironment(website *hypertexting.Website, name string) (*hyperproviders.Environment, error) {
	return website.DeploymentEnvironment(name)
}

func deploymentSession() *hyperproviders.Session {
	return &hyperproviders.Session{
		SecretProviders: []hyperproviders.SecretProvider{envsecrets.New()},
	}
}
