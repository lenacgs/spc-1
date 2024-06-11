package config

import (
	"fmt"
	"strconv"
)

type Version struct {
	BuildNumber int
	CommitHash  string
}

func (v Version) String() string {
	return fmt.Sprintf("%d.%s", v.BuildNumber, v.CommitHash)
}

func convertToInt(val string) int {
	if val == "" {
		return 0
	}
	value, _ := strconv.ParseInt(val, 10, 0)
	return int(value)
}

var currentVersion = Version{
	BuildNumber: convertToInt(BuildNumber), // set using -ldflags
	CommitHash:  CommitHash,
}

func GetCurrentVersion() Version {
	return currentVersion
}

func GetLicenseInfo() string {
	return "workersd is licensed with _____\nsource code available at https://gitlab.com/sibsfps/spc/spc-1"
}
