package coolconf_test

import (
	"os"
	"testing"

	"coolconf"
	"coolconf/test/assert"
)

func TestGetLoadFromEnvVarsWithoutGroup(t *testing.T) {
	os.Setenv("SOME_INT", "9")
	os.Setenv("SOME_STR", "hello")
	os.Setenv("SOME_BOOL", "true")
	os.Setenv(coolconf.SOURCE, coolconf.ENV)
	type MyConfig struct {
		Number int    `envconfig:"SOME_INT"`
		Text   string `envconfig:"SOME_STR"`
		Yes    bool   `envconfig:"SOME_BOOL"`
	}
	var myConfig MyConfig
	conf := coolconf.New(&myConfig)
	conf.Load()
	assert.Equal(t, myConfig.Number, 9)
	assert.Equal(t, myConfig.Text, "hello")
	assert.Equal(t, myConfig.Yes, true)
}

func TestLoadFromEnvVarsUsingGroupAsSuffix(t *testing.T) {
	os.Setenv("OTHER_INT_DEV", "9")
	os.Setenv("OTHER_STR_DEV", "hello")
	os.Setenv("OTHER_BOOL_DEV", "true")
	os.Setenv(coolconf.SOURCE, coolconf.ENV)
	type MyConfig struct {
		Number int    `envconfig:"OTHER_INT"`
		Text   string `envconfig:"OTHER_STR"`
		Yes    bool   `envconfig:"OTHER_BOOL"`
	}
	var myConfig MyConfig
	conf := coolconf.New(&myConfig)
	conf.WithGroup("dev").Load()
	assert.Equal(t, myConfig.Number, 9)
	assert.Equal(t, myConfig.Text, "hello")
	assert.Equal(t, myConfig.Yes, true)
}

func TestLoadFromEnvVarsUsingMultiGroup(t *testing.T) {
	os.Setenv("SOME_STR_QA", "hello")
	os.Setenv("SOME_STR_STAGING", "world")
	os.Setenv(coolconf.SOURCE, coolconf.ENV)
	type MyConfig struct {
		Text   string `envconfig:"SOME_STR"`
	}
	var myConfig1 MyConfig
	var myConfig2 MyConfig
	conf := coolconf.New(&myConfig1)
	conf.WithGroup("qa").Load()
	conf = coolconf.New(&myConfig2)
	conf.WithGroup("staging").Load()
	assert.Equal(t, myConfig1.Text, "hello")
	assert.Equal(t, myConfig2.Text, "world")
}

func TestReloadFromEnvVarsUsingGroupAsSuffix(t *testing.T) {
	os.Setenv(coolconf.SOURCE, coolconf.ENV)
	os.Setenv("some_str_dev", "hello")
	type MyConfig struct {
		Text   string `envconfig:"SOME_STR"`
	}
	var myConfig MyConfig
	conf := coolconf.New(&myConfig)

	config1 := conf.WithGroup("dev").Load().(*MyConfig)
	// assert.Equal(t, ok, true)
	assert.Equal(t, config1.Text, "hello")

	// os.Setenv("some_str_dev", "world")
	var myCleanConfig = MyConfig{}
	conf.WithGroup("dev").Reload(&myCleanConfig)
	// config2 := conf.WithGroup("dev").Load().(*MyConfig)
	// assert.Equal(t, config2.Text, "world")
	assert.Equal(t, myCleanConfig.Text, "world")
}
