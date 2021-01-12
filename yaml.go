package coolconf

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func loadConfigFromYamlFile(destination interface{}, group string) {
	prefix := ""
	suffix := ""
	source := settings.Source
	if group != "" {
		dirname, filename, ext := breakFilePath(source)
		if settings.Yaml.UseGroupAsPrefix {
			prefix = group + settings.Yaml.Separator
			source = fmt.Sprintf("%s%s%s%s", dirname, prefix, filename, ext)
		} else {
			suffix = settings.Yaml.Separator + group
			source = fmt.Sprintf("%s%s%s%s", dirname, filename, suffix, ext)
		}
	}
	b, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalf("[coolconf] Error reading %s: %s", source, err)
	}
	err = yaml.Unmarshal(b, destination)
	if err != nil {
		log.Fatalf("[coolconf] Error unmarshal %s: %s", source, err)
	}
}

func breakFilePath(source string) (string, string, string) {
	filename_with_ext := filepath.Base(source)
	ext := filepath.Ext(filename_with_ext)
	filename := filename_with_ext[:len(filename_with_ext)-len(ext)]
	dirname := source[:len(source)-len(filename_with_ext)]
	return dirname, filename, ext
}
