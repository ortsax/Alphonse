package main

// Build-time version metadata.  Injected via:
//
//	go build -ldflags "-X main.Version=x.y.z -X main.Commit=<sha> -X main.BuildDate=<rfc3339>"
var (
	Version   = "0.0.1"
	Commit    = "unknown"
	BuildDate = "unknown"
)
