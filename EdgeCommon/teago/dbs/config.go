package dbs

// Config mirrors TeaGo/dbs Config for compatibility.
type Config struct {
	DBs     map[string]*DBConfig `yaml:"dbs"`
	Default *DefaultConfig       `yaml:"default,omitempty"`
	Fields  map[string][]string  `yaml:"fields,omitempty"`
}

type DefaultConfig struct {
	DB string `yaml:"db"`
}

// DBConfig holds database connection parameters.
type DBConfig struct {
	Driver          string                `yaml:"driver"`
	Dsn             string                `yaml:"dsn"`
	Prefix          string                `yaml:"prefix"`
	Connections     any                   `yaml:"connections,omitempty"`
	MaxIdle         int                   `yaml:"maxIdle"`
	MaxOpen         int                   `yaml:"maxOpen"`
	MaxLifeSeconds  int                   `yaml:"maxLifeSeconds"`
	ConnTimeoutSecs int                   `yaml:"connTimeoutSeconds"`
	Models          DBModelsConfig        `yaml:"models,omitempty"`
	Fields          map[string]*FieldName `yaml:"fields,omitempty"`
}

type DBModelsConfig struct {
	Package string `yaml:"package"`
}

var globalConfig = &Config{}

// GlobalConfig returns the singleton config.
func GlobalConfig() *Config {
	return globalConfig
}
