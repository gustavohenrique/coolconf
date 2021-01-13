package coolconf_test

import (
	"testing"

	"coolconf"
	"coolconf/test/assert"
)

func getSettingsForEncryptedYaml() coolconf.Settings {
	return coolconf.Settings{
		Source:    "test/testdata/database_encrypted.yaml",
		Encrypted: true,
		SecretKey: "strongpass",
	}
}

func TestLoadFromEncryptedYamlFile(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	coolconf.New(getSettingsForEncryptedYaml())
	coolconf.DecryptYamlFile()
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root:root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}

func TestLoadFromEncryptedYamlFileUsingGroup(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	coolconf.New(getSettingsForEncryptedYaml())
	group := "dev"
	coolconf.DecryptYamlFile(group)
	var myConfig MyConfig
	coolconf.Load(&myConfig, group)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root:root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}
