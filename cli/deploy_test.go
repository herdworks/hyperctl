package cli

import (
	"io"
	"testing"

	"github.com/herdworks/clic"
	"github.com/herdworks/htup"
	"github.com/herdworks/htxt"
)

func TestDeployEnvironment(t *testing.T) {
	production := htup.Environment{"pages"}
	website := &htxt.Website{
		Providers:    htup.HostingProviders{"pages": {Kind: "git", Endpoint: "https://example.com"}},
		Environments: map[string]*htup.Environment{"production": &production},
	}
	execution := &clic.Execution{Stdout: io.Discard}

	environment, name, err := deployEnvironment(website, "production", "", execution)
	if err != nil {
		t.Fatalf("deployEnvironment() error = %v", err)
	}
	if name != "production" {
		t.Fatalf("environment name = %q, want production", name)
	}
	if len(*environment) != 1 || (*environment)[0] != "pages" {
		t.Fatalf("environment = %#v, want [pages]", *environment)
	}
	if _, _, err := deployEnvironment(website, "missing", "", execution); err == nil {
		t.Fatal("missing environment error = nil")
	}
}

func TestDeployEnvironmentProvidersAndDefault(t *testing.T) {
	defaultEnvironment := htup.Environment{"pages", "assets"}
	website := &htxt.Website{
		Providers: htup.HostingProviders{
			"pages":  {Kind: "git", Endpoint: "https://example.com"},
			"assets": {Kind: "s3", Endpoint: "https://example.com"},
		},
		Environments: map[string]*htup.Environment{"default": &defaultEnvironment},
	}
	execution := &clic.Execution{Stdout: io.Discard}

	environment, name, err := deployEnvironment(website, "", "pages, assets", execution)
	if err != nil {
		t.Fatalf("provider selection error = %v", err)
	}
	if name != "" || len(*environment) != 2 || (*environment)[0] != "pages" || (*environment)[1] != "assets" {
		t.Fatalf("provider selection = %#v, %q", environment, name)
	}

	environment, name, err = deployEnvironment(website, "", "", execution)
	if err != nil {
		t.Fatalf("default environment error = %v", err)
	}
	if name != "default" || len(*environment) != 2 {
		t.Fatalf("default environment = %#v, %q", environment, name)
	}
	if _, _, err := deployEnvironment(website, "default", "pages", execution); err == nil {
		t.Fatal("combined environment and providers error = nil")
	}
}

func TestDeployKindAndSecrets(t *testing.T) {
	if got := deployKind(false); got != htup.CompleteDeployment {
		t.Fatalf("deployKind(false) = %q, want complete", got)
	}
	if got := deployKind(true); got != htup.IncrementalDeployment {
		t.Fatalf("deployKind(true) = %q, want incremental", got)
	}
	secrets := deploymentSecrets(nil)
	if len(secrets) != 2 || secrets[0].Name() != "env" || secrets[1].Name() != "files" {
		t.Fatalf("deploymentSecrets(nil) = %#v, want env then files", secrets)
	}
}
