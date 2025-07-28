package nacos

type Nacos struct {
	Ip          string `json:"ip" yaml:"ip"`
	Port        uint64 `json:"port" yaml:"port"`
	ContextPath string `json:"context_path" yaml:"context-path"`
	TimeoutMs   uint64 `json:"timeout_ms" yaml:"timeout-ms"`
	Namespace   string `json:"namespace" yaml:"namespace"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	Cluster     string `json:"cluster" yaml:"cluster"`
	Group       string `json:"group" yaml:"group"`
}
