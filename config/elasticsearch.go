package config

import "log"

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

type ElasticsearchConfig struct {
	Hosts              []string `yaml:"hosts"`
	LogLevel           string   `yaml:"logLevel"`
	SniffingEnabled    bool     `yaml:"sniffingEnabled"`
	HealthcheckEnabled bool     `yaml:"healthcheckEnabled"`
	Index              string   `yaml:"index"`
	Type               string   `yaml:"type"`
}
