package coolconf

import (
	"coolconf/aes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

func loadConfigFromYamlFile(destination interface{}, group string) {
	source := getFilePathWithPrefixOrSuffix(settings, group)
	content := readFile(source)
	unmarshalYaml(content, destination)
}

func decryptYamlFile(group string) error {
	source := getFilePathWithPrefixOrSuffix(settings, group)
	content := readFile(source)
	return decryptYaml(content)
}

func decryptYaml(content []byte) error {
	decoded, err := hex.DecodeString(string(content))
	if err != nil {
		return err
	}
	decrypted, err := aes.Decrypt(settings.SecretKey, decoded)
	if err != nil {
		return err
	}
	settings.Source = decrypted
	return nil
}

func readFile(source string) []byte {
	b, err := ioutil.ReadFile(source)
	if err != nil {
		log.Fatalf("[coolconf] Error reading %s: %s", source, err)
	}
	return b
}

func unmarshalYaml(b []byte, destination interface{}) error {
	return yaml.Unmarshal(b, destination)
}

func getFilePathWithPrefixOrSuffix(settings *Settings, group string) string {
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
	return source
}

func breakFilePath(source string) (string, string, string) {
	filename_with_ext := filepath.Base(source)
	ext := filepath.Ext(filename_with_ext)
	filename := filename_with_ext[:len(filename_with_ext)-len(ext)]
	dirname := source[:len(source)-len(filename_with_ext)]
	return dirname, filename, ext
}
