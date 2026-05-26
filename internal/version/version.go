package version

// Set at link time via -ldflags.
var (
	Version   = "0.1.0"
	Commit    = "unknown"
	BuildTime = "unknown"
)
