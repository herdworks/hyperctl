package cli

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/herdworks/clic"
	"github.com/herdworks/hyperctl"
)

func TestVersionCommand(t *testing.T) {
	var output bytes.Buffer
	result := CLI.Run(clic.Execution{Context: context.Background(), Stdout: &output}, clic.Arguments{"version"})
	if result != (clic.Result{}) {
		t.Fatalf("CLI.Run(version) = %#v", result)
	}
	want := "hyperctl version " + hyperctl.CLIVersion + " (" + hyperctl.CLIReleaseDate.Format(time.RFC3339) + ")\n"
	if got := output.String(); got != want {
		t.Fatalf("version output = %q, want %q", got, want)
	}
}

func TestVersionCommandVerbose(t *testing.T) {
	var output bytes.Buffer
	result := CLI.Run(clic.Execution{Context: context.Background(), Stdout: &output}, clic.Arguments{"version", "--verbose"})
	if result != (clic.Result{}) {
		t.Fatalf("CLI.Run(version --verbose) = %#v", result)
	}
	want := "Package     Version        Release Date\n-------     -------        ------------\n"
	if got := output.String(); len(got) < len(want) || got[:len(want)] != want {
		t.Fatalf("verbose version output = %q, want prefix %q", got, want)
	}
}
