package coolconf_test

import (
	"testing"

	"github.com/gustavohenrique/coolconf"
	"github.com/gustavohenrique/coolconf/test/assert"
)

func getSettingsForEncryptedYamlFile() coolconf.Settings {
	return coolconf.Settings{
		Source:    "test/testdata/database_encrypted.yaml",
		Encrypted: true,
		SecretKey: "strongpass",
	}
}

func getSettingsForEncryptedYaml() coolconf.Settings {
	return coolconf.Settings{
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
	coolconf.New(getSettingsForEncryptedYamlFile())
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
	coolconf.New(getSettingsForEncryptedYamlFile())
	group := "dev"
	coolconf.DecryptYamlFile(group)
	var myConfig MyConfig
	coolconf.Load(&myConfig, group)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root:root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}

func TestLoadFromEncryptedYamlString(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url"`
		Database    struct {
			Port int  `yaml:"port"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	coolconf.New(getSettingsForEncryptedYaml())
	encoded := "460e3aae555acab69f9003db7fa5ae6121a195312173c1d20a32a27bd12250b7396bdb0c2c859b00bd8e30fed9c31828a2157e97d8433e8daf92b1d808c9cd5ffa74fa4d7ab564275c40f418a715dcd487965b6be7969c1ea39239e1988d6119348d4c2e6f08247c2e02e62df476b0a079e6a1b852bf487ba210dd9d89b829c78d2358da5ce0d1cdc44ce0de2dd690bb3dc283ddb8b7"
	coolconf.DecryptYaml([]byte(encoded))
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root:root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, true)
}
