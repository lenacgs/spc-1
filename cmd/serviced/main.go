package main

import (
	"flag"
	"fmt"
	"github.com/gofrs/flock"
	"gitlab.com/sibsfps/spc/spc-1/config"
	"gitlab.com/sibsfps/spc/spc-1/logging"
	"os"
	"path/filepath"
	"strconv"

	service "gitlab.com/sibsfps/spc/spc-1/daemon/serviced"
)

const envLocationTTL = "LOCATION_TTL"
const envUnavailableTTL = "UNAVAILABLE_TTL"

var dataDirectory = flag.String("d", "", "Daemon data path")
var versionCheck = flag.Bool("v", false, "Display and write current build version and exit")
var logToStdout = flag.Bool("o", false, "Write to stdout instead of node.log by overridding config.LogSizeLimit to 0")
var listenIP = flag.String("l", "", "Override config.EndpointAddress (REST listening address) with ip:port")

const envData = "SPC1_DATA"

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

	lockPath := filepath.Join(absolutePath, "serviced.lock")
	fileLock := flock.New(lockPath)
	locked, err := fileLock.TryLock()
	if err != nil {
		fmt.Fprintf(os.Stderr, "unexpected failure in establishing serviced.lock: %s \n", err.Error())
		return 1
	}
	if !locked {
		fmt.Fprintln(os.Stderr, "failed to lock serviced.lock; is an instance of serviced already running in this data directory?")
		return 1
	}
	defer fileLock.Unlock()

	log := logging.Base()

	if logToStdout != nil && *logToStdout {
		cfg.LogSizeLimit = 0
	}

	s := service.Server{
		RootPath: absolutePath,
	}

	if *listenIP != "" {
		cfg.EndpointAddress = *listenIP
	}

	cfg.CacheLocationTTL, cfg.CacheUnavailableTTL, err = readTTLs(log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read TTLs: %v\n", err)
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

func readTTLs(log logging.Logger) (int, int, error) {
	locTTL := os.Getenv(envLocationTTL)
	unavailTTL := os.Getenv(envUnavailableTTL)
	if locTTL == "" || unavailTTL == "" {
		return 0, 0, fmt.Errorf("environment variables UNAVAILABLE_TTL and LOCATION_TTL were not set")
	}

	log.Infof("locationTTL = ", locTTL)
	locationTTL, err := strconv.Atoi(locTTL)
	if err != nil {
		return 0, 0, err
	}

	log.Infof("unavailableTTL = ", unavailTTL)
	unavailableTTL, err := strconv.Atoi(unavailTTL)
	if err != nil {
		return 0, 0, err
	}
	return locationTTL, unavailableTTL, nil
}
