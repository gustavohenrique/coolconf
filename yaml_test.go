package coolconf_test

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gustavohenrique/coolconf"
	"github.com/gustavohenrique/coolconf/test/assert"
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

func TestCreateYamlIfItDoesNotExists(t *testing.T) {
	type MyConfig struct {
		DatabaseURL string `yaml:"database_url" default:"postgres://root@localhost:5432/db"`
		Database    struct {
			Port int  `yaml:"port" default:"5432"`
			SSL  bool `yaml:"ssl"`
		} `yaml:"database"`
	}
	rand.Seed(time.Now().UnixNano())
	filename := fmt.Sprintf("%d.yaml", 100000+rand.Intn(999999-100000))
	coolconf.New(coolconf.Settings{
		Source:                  filepath.Join(os.TempDir(), filename),
		ShouldCreateDefaultYaml: true,
	})
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.DatabaseURL, "postgres://root@localhost:5432/db")
	assert.Equal(t, myConfig.Database.Port, 5432)
	assert.Equal(t, myConfig.Database.SSL, false)
}
