package config

type ProfilerConfig struct {
	Enabled bool   `yaml:"enabled"`
	Host    string `yaml:"host"`
}
