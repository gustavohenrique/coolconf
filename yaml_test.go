package coolconf_test

import (
	"testing"

	"coolconf"
	"coolconf/test/assert"
)

func getSettingsForYaml() coolconf.Settings {
	return coolconf.Settings{
		Source: "test/testdata/database.yaml",
	}
}

func TestLoadFromYamlWithoutGroup(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	coolconf.New(getSettingsForYaml())
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root:root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}

func TestLoadFromYamlUsingGroupAsFilenameSuffix(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	coolconf.New(getSettingsForYaml())
	var myConfig MyConfig
	coolconf.Load(&myConfig, "dev")
	assert.Equal(t, myConfig.DatabaseURL, "postgres://admin:admin@dev:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}

func TestLoadFromYamlUsingGroupAsFilenamePrefix(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	settings := getSettingsForYaml()
	settings.Yaml.UseGroupAsPrefix = true
	coolconf.New(settings)
	var myConfig MyConfig
	coolconf.Load(&myConfig, "dev")
	assert.Equal(t, myConfig.DatabaseURL, "postgres://admin:root@dev:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}
