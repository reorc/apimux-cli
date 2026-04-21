package buildinfo

import "strings"

var (
	Version            = "dev"
	Commit             = "unknown"
	BuildDate          = ""
	ReleaseManifestURL = ""
)

type Info struct {
	Version            string `json:"version"`
	Commit             string `json:"commit"`
	BuildDate          string `json:"build_date,omitempty"`
	ReleaseManifestURL string `json:"release_manifest_url,omitempty"`
}

func Current() Info {
	return Info{
		Version:            fallback(Version, "dev"),
		Commit:             fallback(Commit, "unknown"),
		BuildDate:          strings.TrimSpace(BuildDate),
		ReleaseManifestURL: strings.TrimSpace(ReleaseManifestURL),
	}
}

func fallback(value string, fallbackValue string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return fallbackValue
	}
	return value
}
