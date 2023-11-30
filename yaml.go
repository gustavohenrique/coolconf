package coolconf

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/gustavohenrique/coolconf/aes"
)

func loadConfigFromYamlFile(destination interface{}, group string) error {
	source := getFilePathWithPrefixOrSuffix(settings, group)
	content, err := readFile(source)
	if err != nil {
		return err
	}
	return unmarshalYaml(content, destination)
}

func decryptYamlFile(group string) error {
	source := getFilePathWithPrefixOrSuffix(settings, group)
	content, err := readFile(source)
	if err != nil {
		return err
	}
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

func readFile(source string) ([]byte, error) {
	b := bytes.NewBuffer(nil)
	f, err := os.Open(filepath.Clean(source))
	if err != nil {
		return b.Bytes(), fmt.Errorf("Error reading %s: %s", source, err)
	}
	defer f.Close()
	_, err = io.Copy(b, f)
	if err != nil {
		return b.Bytes(), fmt.Errorf("Error copying bytes %s: %s", source, err)
	}
	return []byte(os.ExpandEnv(b.String())), nil
}

func unmarshalYaml(b []byte, destination interface{}) error {
	return yaml.Unmarshal(b, destination)
}

func createYamlUsingDefaultConfigIfItDoesNotExists(source string, destination interface{}) error {
	err := process("", "", "yaml", destination)
	if err != nil {
		return err
	}
	if !isFileExist(source) {
		b, err := yaml.Marshal(destination)
		if err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Dir(source), os.ModePerm); err != nil {
			return err
		}
		return os.WriteFile(source, b, 0600)
	}
	return nil
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
