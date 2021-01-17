package shell

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/doron-cohen/antidot/internal/utils"
)

type KeyValueExist struct {
	Key string
}

func (k *KeyValueExist) Error() string {
	return fmt.Sprintf("Key '%s' already exist with the requested value", k.Key)
}

type keyValueMapFormat struct {
	parseRegexp *regexp.Regexp
	format      string
}

func NewKeyValueMapFormat(pattern, format string) *keyValueMapFormat {
	return &keyValueMapFormat{
		regexp.MustCompile(pattern),
		format,
	}
}

func (k keyValueMapFormat) ParseKeyValuePairs(text string) (map[string]string, error) {
	matches := k.parseRegexp.FindAllStringSubmatch(text, -1)
	results := make(map[string]string, len(matches))
	for _, match := range matches {
		if len(match) != 3 {
			return nil, fmt.Errorf("Invalid match for key value RegEx: %v", match)
		}

		key := match[1]
		value := match[2]
		if _, exists := results[key]; exists {
			return nil, fmt.Errorf("Key %s appears twice", key)
		}

		results[key] = value
	}

	return results, nil
}

func (k keyValueMapFormat) WriteKeyValuePairs(wr io.Writer, keyValues map[string]string) error {
	for key, value := range keyValues {
		_, err := fmt.Fprintf(wr, k.format, key, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func (k keyValueMapFormat) AddKeyValueToString(text, key, value string) (string, error) {
	keyValues, err := k.ParseKeyValuePairs(text)
	if err != nil {
		return "", err
	}

	if currValue, exists := keyValues[key]; exists {
		if currValue != value {
			return "", fmt.Errorf("Key %s already exists with different value", key)
		}
		// Key value already exists
		return "", &KeyValueExist{Key: key}
	}

	keyValues[key] = value
	writer := &strings.Builder{}
	err = k.WriteKeyValuePairs(writer, keyValues)
	if err != nil {
		return "", err
	}

	return writer.String(), nil
}

func AppendKeyValueToFile(path, key, value string, kvMapFormat *keyValueMapFormat) error {
	text := ""
	if utils.FileExists(path) {
		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		text = string(bytes[:])
	}

	result, err := kvMapFormat.AddKeyValueToString(text, key, value)
	if err != nil {
		if _, ok := err.(*KeyValueExist); ok {
			return nil
		}
		return err
	}

	data := []byte(result)
	err = ioutil.WriteFile(path, data, os.FileMode(0o755))
	if err != nil {
		return err
	}

	return nil
}
