package coolconf_test

import (
	"os"
	"testing"

	"coolconf"
	"coolconf/test/assert"
)

func TestGetLoadFromEnvWithoutGroup(t *testing.T) {
	os.Setenv("SOME_INT", "9")
	os.Setenv("SOME_STR", "hi")
	os.Setenv("SOME_BOOL", "true")
	type MyConfig struct {
		Number int    `env:"SOME_INT"`
		Text   string `env:"SOME_STR"`
		Yes    bool   `env:"SOME_BOOL"`
	}
	coolconf.New()
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.Number, 9)
	assert.Equal(t, myConfig.Text, "hi")
	assert.Equal(t, myConfig.Yes, true)
}

func TestLoadFromEnvUsingGroupAsSuffix(t *testing.T) {
	os.Setenv("OTHER_INT_DEV", "9")
	os.Setenv("OTHER_STR_DEV", "hello")
	os.Setenv("OTHER_BOOL_DEV", "true")
	type MyConfig struct {
		Number int    `env:"OTHER_INT"`
		Text   string `env:"OTHER_STR"`
		Yes    bool   `env:"OTHER_BOOL"`
	}
	coolconf.New()
	var myConfig MyConfig
	coolconf.Load(&myConfig, "DEV")
	assert.Equal(t, myConfig.Number, 9)
	assert.Equal(t, myConfig.Text, "hello")
	assert.Equal(t, myConfig.Yes, true)
}

func TestLoadFromEnvUsingMultiGroup(t *testing.T) {
	os.Setenv("SOME_STR", "helloworld")
	os.Setenv("SOME_STR_QA", "hello")
	os.Setenv("SOME_STR_STAGING", "world")
	type MyConfig struct {
		Text string `env:"SOME_STR"`
	}
	coolconf.New()
	var qa, staging, global MyConfig
	coolconf.Load(&qa, "QA")
	coolconf.Load(&staging, "STAGING")
	coolconf.Load(&global)
	assert.Equal(t, qa.Text, "hello")
	assert.Equal(t, staging.Text, "world")
	assert.Equal(t, global.Text, "helloworld")
}

func TestClearOldValuesFromEnvVar(t *testing.T) {
	os.Setenv("another_str_dev", "hello")
	type MyConfig struct {
		Text string `env:"another_str"`
	}
	coolconf.New()
	var old, updated MyConfig

	coolconf.Load(&old, "dev")
	assert.Equal(t, old.Text, "hello")

	// Update the envvar does not change de struct value
	os.Setenv("another_str_dev", "world")
	coolconf.Load(&updated, "dev")
	assert.Equal(t, old.Text, "hello")

	os.Setenv("another_str_dev", "helloworld")
	coolconf.Clear()
	coolconf.Load(&old, "dev")
	assert.Equal(t, old.Text, "helloworld")
}

func TestLoadFromEnvUsingDefaultValue(t *testing.T) {
	type MyConfig struct {
		Text   string `env:"MY_TEXT" default:"hello"`
		Number int    `env:"MY_NUMBER" default:"9"`
	}
	coolconf.New()
	var myConfig MyConfig
	coolconf.Load(&myConfig)
	assert.Equal(t, myConfig.Text, "hello")
	assert.Equal(t, myConfig.Number, 9)
}
