package coolconf

import (
	"context"
	"errors"
	"os"
	"strings"
)

const (
	ENV               = "env"
	YAML_FILE         = "yaml_file"
	DEFAULT_SEPARATOR = "_"
)

var configs context.Context
var settings *Settings

type Option struct {
	UseGroupAsPrefix bool
	Separator        string
}

type Settings struct {
	Source                  string
	SecretKey               string
	Encrypted               bool
	Key                     string
	Env                     Option
	Yaml                    Option
	ShouldCreateDefaultYaml bool
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

func DecryptYamlFile(params ...string) error {
	if settings.SecretKey == "" || !settings.Encrypted || !isYamlFile(settings.Source) {
		return errors.New("Settings does not contains the secret key or encrypted=false or source is not a YAML file")
	}
	var group string
	if len(params) > 0 {
		group = params[0]
	}
	err := decryptYamlFile(group)
	if err != nil {
		return errors.New("Failed to decrypt YAML file: " + err.Error())
	}
	return nil
}

func DecryptYaml(encoded []byte) error {
	if settings.SecretKey == "" || !settings.Encrypted {
		return errors.New("Cannot decrypt YAML string")
	}
	err := decryptYaml(encoded)
	if err != nil {
		return errors.New("Failed to decrypt YAML: " + err.Error())
	}
	return nil
}

func Load(destination interface{}, params ...string) error {
	var group string
	if len(params) > 0 {
		group = params[0]
	}
	if val := configs.Value(group); val != nil {
		destination = val
		return nil
	}
	err := loadTo(destination, group)
	if err != nil {
		return err
	}
	configs = context.WithValue(configs, group, destination)
	return nil
}

func loadTo(destination interface{}, group string) error {
	source := settings.Source
	mode := source
	if isYamlFile(source) {
		mode = YAML_FILE
	}
	var err error
	switch mode {
	case YAML_FILE:
		if settings.ShouldCreateDefaultYaml {
			err = createYamlUsingDefaultConfigIfItDoesNotExists(settings.Source, destination)
			if err != nil {
				return err
			}
		}
		err = loadConfigFromYamlFile(destination, group)
	case ENV:
		err = loadConfigFromEnv(destination, group)
	default:
		err = unmarshalYaml([]byte(settings.Source), destination)
	}
	return err
}

func isYamlFile(filename string) bool {
	f := strings.ToLower(filename)
	isYaml := strings.HasSuffix(f, ".yaml") || strings.HasSuffix(f, ".yml")
	return isYaml // && isFile(filename)
}

func isFile(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
