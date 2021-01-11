package coolconf

import (
	"os"
	"strings"
)

const (
	SOURCE = "COOLCONF_SOURCE"
	ENV    = "env"
	FILE   = "file"
)

var configs = make(map[string]interface{})

type CoolConf struct {
	source      string
	group       string
	destination interface{}
	config      Config
}

type Config struct {
	UseGroupAsPrefix  bool
	UseLowerCaseGroup bool
	Separator         string
}

func New(destination interface{}, params ...Config) *CoolConf {
	source := os.Getenv(SOURCE)
	if source == "" {
		source = ENV
	}
	config := Config{
		Separator: "_",
	}
	if len(params) > 0 {
		config = params[0]
	}
	return &CoolConf{
		source:      source,
		destination: destination,
		config:      config,
	}
}

func (c *CoolConf) NoGroup() *CoolConf {
	c.group = ""
	return c
}

func (c *CoolConf) WithGroup(group string) *CoolConf {
	if c.config.UseLowerCaseGroup {
		c.group = strings.ToLower(group)
	} else {
		c.group = strings.ToUpper(group)
	}
	return c
}

func (c *CoolConf) Load() interface{} {
	if val, ok := configs[c.group]; ok {
		return val
	}
	c.loadConfigFromEnvVar()
	configs[c.group] = c.destination
	return c.destination
}

func (c *CoolConf) Reload(destination interface{}) {
	c.destination = destination
	c.loadConfigFromEnvVar()
	configs[c.group] = destination
}

func (c *CoolConf) loadConfigFromEnvVar() {
	prefix := ""
	suffix := ""
	if c.config.UseGroupAsPrefix {
		prefix = c.group + c.config.Separator
	} else {
		suffix = c.config.Separator + c.group
	}
	Process(prefix, suffix, c.destination)
}
