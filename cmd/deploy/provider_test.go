package deploy

import (
	"testing"

	hyperproviders "github.com/herdworks/hyperproviders-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func TestDeploymentEnvironment(t *testing.T) {
	website := &hypertexting.Website{
		Environments: map[string]*hyperproviders.Environment{
			"production": {},
		},
	}

	env, err := deploymentEnvironment(website, "production")
	if err != nil {
		t.Fatalf("deploymentEnvironment: %v", err)
	}
	if env.Name != "production" {
		t.Fatalf("Name = %q, want production", env.Name)
	}

	if _, err := deploymentEnvironment(website, "missing"); err == nil {
		t.Fatal("deploymentEnvironment missing error = nil, want error")
	}
}

func TestDeploymentSessionUsesEnvSecrets(t *testing.T) {
	session := deploymentSession()
	if session == nil {
		t.Fatal("deploymentSession returned nil")
	}
	if len(session.SecretProviders) != 1 {
		t.Fatalf("SecretProviders len = %d, want 1", len(session.SecretProviders))
	}
	if got := session.SecretProviders[0].Name(); got != "env" {
		t.Fatalf("SecretProviders[0].Name() = %q, want env", got)
	}
}
