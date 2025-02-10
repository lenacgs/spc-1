package config

var defaultLocal = Local{
	BaseLoggerDebugLevel:    4,
	LogFileDir:              "",
	LogArchiveDir:           "",
	LogArchiveName:          "node.archive.log",
	LogArchiveMaxAge:        "",
	LogSizeLimit:            1073741824,
	NetAddress:              "",
	EndpointAddress:         "127.0.0.1:0",
	RestReadTimeoutSeconds:  15,
	RestWriteTimeoutSeconds: 120,
	CacheLocationTTL:        1000,
	CacheUnavailableTTL:     100,
	CacheMaxCapacity:        100,
}
