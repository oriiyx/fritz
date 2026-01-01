package version

import (
	"fmt"
	"runtime"
)

// Build information. Populated at build-time via -ldflags.
var (
	Version   = "dev"
	BuildDate = "unknown"
	GitCommit = "unknown"
)

// BuildInfo contains version and build information.
type BuildInfo struct {
	Version   string `json:"version"`
	BuildDate string `json:"buildDate"`
	GitCommit string `json:"gitCommit"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// GetBuildInfo returns build information.
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:   Version,
		BuildDate: BuildDate,
		GitCommit: GitCommit,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// GetVersion returns the version string.
func GetVersion() string {
	if Version == "dev" {
		return fmt.Sprintf("%s-%s", Version, GitCommit)
	}
	return Version
}

// GetFullVersion returns a detailed version string.
func GetFullVersion() string {
	info := GetBuildInfo()
	return fmt.Sprintf("Fritz CLI %s\nBuild Date: %s\nGit Commit: %s\nGo Version: %s\nPlatform: %s",
		info.Version,
		info.BuildDate,
		info.GitCommit,
		info.GoVersion,
		info.Platform,
	)
}
