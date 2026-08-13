package hyperctl

import (
	"time"
)

// ============================ Package Variables ============================
var (
	CLIVersion        = "v0.23.0"
	CLIReleaseDate, _ = time.Parse(time.RFC3339, "2026-07-30T14:00:00-07:00") // Usage: CLIReleaseDate.Format(time.RFC3339)
)
