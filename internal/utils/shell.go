package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type KeyValueMap interface {
	String() string
}

type EnvMap map[string]string

func (e EnvMap) String() string {
	var line string
	var builder strings.Builder
	for key, value := range e {
		line = fmt.Sprintf("export %s=\"%s\"\n", key, value)
		builder.WriteString(line)
	}

	return builder.String()
}

func EnvMapFromFile(filePath string) (EnvMap, error) {
	result, err := LoadKeyValuesFromFile(filePath, "export")
	if err != nil {
		return nil, err
	}

	envMap := EnvMap(result)
	return envMap, nil
}

type AliasMap map[string]string

func (a AliasMap) String() string {
	var line string
	var builder strings.Builder
	for alias, command := range a {
		line = fmt.Sprintf("alias %s=\"%s\"\n", alias, command)
		builder.WriteString(line)
	}

	return builder.String()
}

func AliasMapFromFile(filePath string) (AliasMap, error) {
	result, err := LoadKeyValuesFromFile(filePath, "alias")
	if err != nil {
		return nil, err
	}

	aliasMap := AliasMap(result)
	return aliasMap, nil
}

func LoadKeyValuesFromFile(filePath, prefix string) (map[string]string, error) {
	pattern := fmt.Sprintf("^%s (?P<key>\\w+)=\"(?P<value>.*)\"", prefix)
	lineRe, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		match := lineRe.FindStringSubmatch(scanner.Text())
		if len(match) == 0 {
			continue
		}

		var key, value string
		for i, name := range lineRe.SubexpNames() {
			if i == 0 {
				continue
			}

			switch name {
			case "key":
				key = match[i]
			case "value":
				value = match[i]
			default:
				errMessage := fmt.Sprintf("Unknown RegEx group name '%s'", name)
				return nil, errors.New(errMessage)
			}
		}
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func WriteKeyValuesToFile(keyValueMap KeyValueMap, filePath string) error {
	str := keyValueMap.String()
	data := []byte(str)
	if err := ioutil.WriteFile(filePath, data, 0o644); err != nil {
		return err
	}

	return nil
}
