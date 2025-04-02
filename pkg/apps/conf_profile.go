package apps

type ProfileType string

const (
	Pprof  ProfileType = "pprof"
	FGprof ProfileType = "fgprof"
	None   ProfileType = "none"
)

type Profiling struct {
	Enabled      bool        `json:"enabled" yaml:"enabled"`
	AuthDownload bool        `json:"auth_download" yaml:"auth-download"`
	Type         ProfileType `json:"type" yaml:"type"`
	Prefix       string      `json:"prefix" yaml:"prefix"`
}
