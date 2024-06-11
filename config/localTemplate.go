package config

type Local struct {
	BaseLoggerDebugLevel uint32
	LogFileDir           string
	LogArchiveDir        string
	LogArchiveName       string
	LogArchiveMaxAge     string
	LogSizeLimit         uint64

	NetAddress      string
	EndpointAddress string

	RestReadTimeoutSeconds  int
	RestWriteTimeoutSeconds int
}
