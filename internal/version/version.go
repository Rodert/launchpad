package version

import (
	"fmt"
	"runtime"
)

var (
	Version      string
	Stage        string
	GitCommit    string
	GitMessage   string
	GitCommitter string
	Built        string
)

func VersionInfo() string {
	format := `
Version:      %v
Stage:        %v
GitCommit:    %v
GitMessage:   %v
GitCommitter: %v
GoVersion:    %v
Built:        %v
OSArch:       %v/%v
	`
	return fmt.Sprintf(
		format,
		Version,
		Stage,
		GitCommit,
		GitMessage,
		GitCommitter,
		runtime.Version(),
		Built,
		runtime.GOOS,
		runtime.GOARCH,
	)
}

func GetVersion() string {
	return Version
}