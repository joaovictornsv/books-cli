package buildinfo

import "runtime"

var (
	Version = "0.5.0"
	Commit  = "unknown"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	GoVersion string `json:"go_version"`
}

func Get() Info {
	return Info{
		Version:   Version,
		Commit:    Commit,
		GoVersion: runtime.Version(),
	}
}
