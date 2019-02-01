package config

import (
	"advanced-database-course-project-server/constants"
	"advanced-database-course-project-server/log"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Config *Configuration

type Configuration struct {
	RestApiConfig       RestApiConfig       `yaml:"restApi"`
	ElasticsearchConfig ElasticsearchConfig `yaml:"elasticsearch"`
	ProfilerConfig      ProfilerConfig      `yaml:"profiler"`
	LogLevel            string              `yaml:"logLevel"`
	PrettyLog           bool                `yaml:"prettyLog"`
	GoMaxProcs          int                 `yaml:"goMaxProcs"`
}

func (config *Configuration) getConfigFromYAMLFile() error {
	yamlFile, err := ioutil.ReadFile(constants.ConfigYamlFileName)
	if err != nil {
		return errors.Wrap(err, "failed to read yaml file")
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return errors.Wrap(err, "failed to Unmarshal yaml file")
	}

	return nil
}

func LoadConfiguration() {

	loadedConfig := Configuration{}
	err := loadedConfig.getConfigFromYAMLFile()

	if err != nil {
		log.StdoutLogger.
			WithError(err).
			WithField("file_name", constants.ConfigYamlFileName).
			Fatal("couldn't load configuration from yaml file")
	}

	Config = &loadedConfig

	log.StdoutLogger.
		WithField("file_name", constants.ConfigYamlFileName).
		Info("configuration loaded from yaml file")
}
