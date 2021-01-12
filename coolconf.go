package coolconf

import (
	"context"
)

const (
	ENV               = "env"
	FILE              = "file"
	DEFAULT_SEPARATOR = "_"
)

var configs context.Context
var settings *Settings

type Env struct {
	UseGroupAsPrefix bool
	Separator        string
}

type Settings struct {
	Source string
	Env    Env
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

func LoadConfigFromEnvVar(destination interface{}, group string) {
	prefix := ""
	suffix := ""
	if group != "" {
		if settings.Env.UseGroupAsPrefix {
			prefix = group + settings.Env.Separator
		} else {
			suffix = settings.Env.Separator + group
		}
	}
	Process(prefix, suffix, destination)
}

func loadTo(destination interface{}, group string) {
	switch settings.Source {
	default:
		LoadConfigFromEnvVar(destination, group)
	}
}
