package TeaGo

// ServerConfig is a minimal server.yaml config struct used by installer.
type ServerConfig struct {
	Env  string `yaml:"env,omitempty"`
	Http struct {
		On     bool     `yaml:"on,omitempty"`
		Listen []string `yaml:"listen,omitempty"`
	} `yaml:"http,omitempty"`
}
