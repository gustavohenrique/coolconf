package coolconf

import (
	"context"
	"log"
	"os"
	"strings"
)

const (
	ENV               = "env"
	YAML_FILE         = "yaml_file"
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
	Source    string
	SecretKey string
	Encrypted bool
	Key       string
	Env       Option
	Yaml      Option
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

func DecryptYamlFile(params ...string) {
	if settings.SecretKey == "" || !settings.Encrypted || !isYamlFile(settings.Source) {
		log.Fatalln("[error] Settings does not contains the secret key or encrypted=false or source is not a YAML file")
	}
	var group string
	if len(params) > 0 {
		group = params[0]
	}
	err := decryptYamlFile(group)
	if err != nil {
		log.Fatalln("[error] Failed to decrypt YAML file:", err)
	}
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
	source := settings.Source
	mode := source
	if isYamlFile(source) {
		mode = YAML_FILE
	}
	switch mode {
	case YAML_FILE:
		loadConfigFromYamlFile(destination, group)
	case ENV:
		loadConfigFromEnv(destination, group)
	default:
		unmarshalYaml([]byte(settings.Source), destination)
	}
}

func isYamlFile(filename string) bool {
	f := strings.ToLower(filename)
	isYaml := strings.HasSuffix(f, ".yaml") || strings.HasSuffix(f, ".yml")
	return isYaml && isFile(filename)
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
