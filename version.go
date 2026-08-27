package hyperctl

import (
	"time"
)

// ============================ Package Variables ============================
var (
	CLIVersion        = "v0.24.0-rc.14"
	CLIReleaseDate, _ = time.Parse(time.RFC3339, "2026-08-26T17:00:00-07:00") // Usage: CLIReleaseDate.Format(time.RFC3339)
)
