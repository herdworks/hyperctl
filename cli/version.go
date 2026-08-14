package cli

import (
	"fmt"
	"time"

	"github.com/herdworks/clic"
	hqversion "github.com/herdworks/hq/version"
	"github.com/herdworks/hsync"
	htskills "github.com/herdworks/htskills"
	"github.com/herdworks/httt"
	"github.com/herdworks/htup"
	"github.com/herdworks/htxt"
	"github.com/herdworks/hyperctl"
	hyperfeedversion "github.com/herdworks/hyperfeed/version"
	"github.com/herdworks/hypermedia"
)

// VersionCommand prints Hyperctl version information.
func VersionCommand() clic.Command {
	verbose := GlobalVerboseOption()
	verbose.Meta.Description = "Print HyperTexting package versions."
	return clic.Command{
		Meta: clic.CommandMeta{
			Name:        "version",
			Description: "Print version information.",
			Usage:       []string{"hyperctl version"},
			Order:       1000,
		},
		Spec: clic.CommandSpec{
			Key:       "version",
			Help:      clic.DefaultHelp,
			Validator: expectedPositionalArguments(0),
			Options: clic.Options{
				verbose.Spec.Key: &verbose,
			},
			Handler: versionHandler,
		},
	}
}

func versionHandler(execution *clic.Execution) error {
	if execution.Options["verbose"].BoolValue() {
		return clic.WriteList(execution.Stdout, clic.ColumnarFormat, packageVersions(),
			clic.Field{Key: "package", Header: "Package"},
			clic.Field{Key: "version", Header: "Version"},
			clic.Field{Key: "release_date", Header: "Release Date"},
		)
	}
	_, err := fmt.Fprintf(execution.Stdout, "hyperctl version %s (%s)\n", hyperctl.CLIVersion, hyperctl.CLIReleaseDate.Format(time.RFC3339))
	return err
}

// PackageVersion identifies one bundled Herdworks module release.
type PackageVersion struct {
	Package     string
	Version     string
	ReleaseDate time.Time
}

// PackageVersions is a collection of package version metadata.
type PackageVersions []PackageVersion

// ColumnValues returns the package version table columns.
func (versions PackageVersions) ColumnValues(field string) (values []string) {
	for _, version := range versions {
		switch field {
		case "package":
			values = append(values, version.Package)
		case "version":
			values = append(values, version.Version)
		case "release_date":
			values = append(values, version.ReleaseDate.Format(time.RFC3339))
		}
	}
	return values
}

// packageVersions returns the Herdworks modules linked into Hyperctl.
func packageVersions() PackageVersions {
	return PackageVersions{
		{Package: "clic", Version: clic.LibVersion, ReleaseDate: clic.LibReleaseDate},
		{Package: "hq", Version: hqversion.LibVersion, ReleaseDate: hqversion.LibReleaseDate},
		{Package: "hsync", Version: hsync.LibVersion, ReleaseDate: hsync.LibReleaseDate},
		{Package: "htskills", Version: htskills.LibVersion, ReleaseDate: htskills.LibReleaseDate},
		{Package: "httt", Version: httt.LibVersion, ReleaseDate: httt.LibReleaseDate},
		{Package: "htup", Version: htup.LibVersion, ReleaseDate: htup.LibReleaseDate},
		{Package: "htxt", Version: htxt.LibVersion, ReleaseDate: htxt.LibReleaseDate},
		{Package: "hyperfeed", Version: hyperfeedversion.LibVersion, ReleaseDate: hyperfeedversion.LibReleaseDate},
		{Package: "hypermedia", Version: hypermedia.LibVersion, ReleaseDate: hypermedia.LibReleaseDate},
		{Package: "hyperctl", Version: hyperctl.CLIVersion, ReleaseDate: hyperctl.CLIReleaseDate},
	}
}
