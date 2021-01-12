package coolconf

import (
	"context"
	"os"
	"strings"
)

const (
	ENV               = "env"
	YAML              = "yaml"
	DEFAULT_SEPARATOR = "_"
)

var configs context.Context
var settings *Settings

type Option struct {
	UseGroupAsPrefix bool
	Separator        string
}

type Settings struct {
	Source string
	Env    Option
	Yaml   Option
}

func New(params ...Settings) {
	Clear()
	settings = &Settings{}
	if len(params) > 0 {
		settings = &params[0]
	}
	if settings.Env.Separator == "" {
		settings.Env.Separator = DEFAULT_SEPARATOR
	}
	if settings.Yaml.Separator == "" {
		settings.Yaml.Separator = DEFAULT_SEPARATOR
	}
	if settings.Source == "" {
		settings.Source = ENV
	}
}

func Clear() {
	configs = context.Background()
}

func Load(destination interface{}, params ...string) interface{} {
	var group string
	if len(params) > 0 {
		group = params[0]
	}
	if val := configs.Value(group); val != nil {
		destination = val
		return val
	}
	loadTo(destination, group)
	configs = context.WithValue(configs, group, destination)
	return destination
}

func loadTo(destination interface{}, group string) {
	var mode string
	source := settings.Source
	if isYaml(source) && isFile(source) {
		mode = YAML
	}
	switch mode {
	case YAML:
		loadConfigFromYamlFile(destination, group)
	default:
		loadConfigFromEnv(destination, group)
	}
}

func isYaml(filename string) bool {
	f := strings.ToLower(filename)
	return strings.HasSuffix(f, ".yaml") || strings.HasSuffix(f, ".yml")
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
