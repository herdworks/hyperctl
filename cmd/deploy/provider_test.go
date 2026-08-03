package deploy

import (
	"testing"

	hyperproviders "github.com/herdworks/hyperproviders-go"
	hypertexting "github.com/herdworks/hypertexting-go"
)

func TestDeploymentEnvironment(t *testing.T) {
	production := hyperproviders.Environment{"pages"}
	website := &hypertexting.Website{
		Providers: hyperproviders.HostingProviders{
			"pages": {Kind: "git", Endpoint: "https://example.com"},
		},
		Environments: map[string]*hyperproviders.Environment{
			"production": &production,
		},
	}

	env, environmentName, err := deploymentEnvironment(website, "production", "")
	if err != nil {
		t.Fatalf("deploymentEnvironment: %v", err)
	}
	if environmentName != "production" {
		t.Fatalf("environmentName = %q, want production", environmentName)
	}
	if len(*env) != 1 || (*env)[0] != "pages" {
		t.Fatalf("environment = %#v, want [pages]", *env)
	}

	if _, _, err := deploymentEnvironment(website, "missing", ""); err == nil {
		t.Fatal("deploymentEnvironment missing error = nil, want error")
	}
}

func TestDeploymentEnvironmentProvidersList(t *testing.T) {
	website := &hypertexting.Website{
		Providers: hyperproviders.HostingProviders{
			"pages":  {Kind: "git", Endpoint: "https://example.com"},
			"assets": {Kind: "s3", Endpoint: "https://example.com"},
		},
	}

	env, environmentName, err := deploymentEnvironment(website, "", "pages, assets")
	if err != nil {
		t.Fatalf("deploymentEnvironment: %v", err)
	}
	if environmentName != "" {
		t.Fatalf("environmentName = %q, want empty", environmentName)
	}
	if len(*env) != 2 || (*env)[0] != "pages" || (*env)[1] != "assets" {
		t.Fatalf("environment = %#v, want pages/assets", *env)
	}
}

func TestDeploymentEnvironmentDefault(t *testing.T) {
	defaultEnv := hyperproviders.Environment{"default"}
	website := &hypertexting.Website{
		Providers: hyperproviders.HostingProviders{
			"default": {Kind: "git", Endpoint: "https://example.com"},
		},
		Environments: map[string]*hyperproviders.Environment{
			"default": &defaultEnv,
		},
	}

	env, environmentName, err := deploymentEnvironment(website, "", "")
	if err != nil {
		t.Fatalf("deploymentEnvironment: %v", err)
	}
	if environmentName != "default" {
		t.Fatalf("environmentName = %q, want default", environmentName)
	}
	if len(*env) != 1 || (*env)[0] != "default" {
		t.Fatalf("environment = %#v, want default", *env)
	}
}

func TestDeploymentEnvironmentRejectsEnvironmentAndProviders(t *testing.T) {
	website := &hypertexting.Website{}
	if _, _, err := deploymentEnvironment(website, "production", "pages"); err == nil {
		t.Fatal("deploymentEnvironment both selectors error = nil, want error")
	}
}

func TestParseProviderNames(t *testing.T) {
	got := parseProviderNames("pages, assets,,")
	if len(got) != 2 || got[0] != "pages" || got[1] != "assets" {
		t.Fatalf("parseProviderNames = %#v, want pages/assets", got)
	}
}

func TestDeploymentKind(t *testing.T) {
	if got := deploymentKind(false); got != hyperproviders.CompleteDeployment {
		t.Fatalf("deploymentKind(false) = %q, want complete", got)
	}
	if got := deploymentKind(true); got != hyperproviders.IncrementalDeployment {
		t.Fatalf("deploymentKind(true) = %q, want incremental", got)
	}
}

func TestDeploymentSecretsUsesEnvThenFilesSecrets(t *testing.T) {
	secrets := DeploymentSecrets(nil)
	if len(secrets) != 2 {
		t.Fatalf("SecretProviders len = %d, want 2", len(secrets))
	}
	if got := secrets[0].Name(); got != "env" {
		t.Fatalf("SecretProviders[0].Name() = %q, want env", got)
	}
	if got := secrets[1].Name(); got != "files" {
		t.Fatalf("SecretProviders[1].Name() = %q, want files", got)
	}
}
