package config

import (
	"fmt"
	"path/filepath"
)

func GetDefaultLocal() Local {
	return defaultLocal
}

func (cfg *Local) ResolveLogPaths(rootDir string) (liveLog, archive string) {
	// the default locations of log and archive are root
	liveLog = filepath.Join(rootDir, "node.log")
	archive = filepath.Join(rootDir, cfg.LogArchiveName)

	// if LogFileDir is set, use it instead
	if cfg.LogFileDir != "" {
		liveLog = filepath.Join(cfg.LogFileDir, "node.log")
		archive = filepath.Join(cfg.LogFileDir, cfg.LogArchiveName)
	}

	// if LogArchivePath is set, use it instead
	if cfg.LogArchiveDir != "" {
		archive = filepath.Join(cfg.LogArchiveDir, cfg.LogArchiveName)
	}

	return liveLog, archive
}

func FormatVersionAndLicense() string {
	version := GetCurrentVersion()
	return fmt.Sprintf("%d\n (commit #%s)\n%s",
		version.BuildNumber, version.CommitHash, GetLicenseInfo(),
	)
}
