package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gitlab.com/sibsfps/spc/spc-1/config"
	workers "gitlab.com/sibsfps/spc/spc-1/daemon/workersd"
	"gitlab.com/sibsfps/spc/spc-1/logging"

	"github.com/gofrs/flock"
)

const envData = "SPC1_DATA"

var dataDirectory = flag.String("d", "", "Daemon data path")
var versionCheck = flag.Bool("v", false, "Display and write current build version and exit")
var logToStdout = flag.Bool("o", false, "Write to stdout instead of node.log by overriding config.LogSizeLimit to 0")
var listenIP = flag.String("l", "", "Override config.EndpointAddress (REST listening address) with ip:port")

func main() {
	flag.Parse()

	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	var err error
	cfg := config.GetDefaultLocal()

	if *versionCheck {
		fmt.Println(config.FormatVersionAndLicense())
		return 0
	}

	dataDir := resolveDataDir()
	absolutePath, absPathErr := filepath.Abs(dataDir)

	if len(dataDir) == 0 {
		fmt.Fprintf(os.Stderr, "Data directory not specified. Please use -d or set $%s in your environment.\n", envData)
		return 1
	}

	if absPathErr != nil {
		fmt.Fprintf(os.Stderr, "Can't convert data directory's path to absolute, %v\n", dataDir)
		return 1
	}

	if _, err := os.Stat(absolutePath); err != nil {
		fmt.Fprintf(os.Stderr, "Data directory %s does not appear to be valid\n", dataDir)
		return 1
	}

	lockPath := filepath.Join(absolutePath, "workersd.lock")
	fileLock := flock.New(lockPath)
	locked, err := fileLock.TryLock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected failure in establishing workersd.lock: %s \n", err.Error())
		return 1
	}
	if !locked {
		fmt.Fprintln(os.Stderr, "failed to lock workersd.lock; is an instance of workersd already running in this data directory?")
		return 1
	}
	defer fileLock.Unlock()

	log := logging.Base()

	if logToStdout != nil && *logToStdout {
		cfg.LogSizeLimit = 0
	}

	s := workers.Server{
		RootPath: absolutePath,
	}

	if *listenIP != "" {
		cfg.EndpointAddress = *listenIP
	}

	err = s.Initialize(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		log.Error(err)
		return 1
	}

	s.Start()

	return 0
}

func resolveDataDir() string {
	// If not specified on cmdline with '-d', look for default in environment.
	var dir string
	if dataDirectory == nil || *dataDirectory == "" {
		dir = os.Getenv(envData)
	} else {
		dir = *dataDirectory
	}
	return dir
}
