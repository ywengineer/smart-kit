package apps

type ProfileType string

const (
	Pprof  ProfileType = "pprof"
	FGprof ProfileType = "fgprof"
	None   ProfileType = "none"
)

type Profiling struct {
	Enabled bool        `json:"enabled" yaml:"enabled"`
	Type    ProfileType `json:"type" yaml:"type"`
}
