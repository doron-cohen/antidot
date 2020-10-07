package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

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
	result, err := LoadKeyValuesFromFile(filePath)
	if err != nil {
		return nil, err
	}

	aliasMap := AliasMap(result)
	return aliasMap, nil
}

func LoadKeyValuesFromFile(filePath string) (map[string]string, error) {
	lineRe := regexp.MustCompile(`^alias (?P<key>\w+)="(?P<value>.*)"`)
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

			// TODO error on default
			switch name {
			case "key":
				// TODO: assert key is not an empty string
				key = match[i]
			case "value":
				value = match[i]
			}
		}
		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func WriteAliasesToFile(aliasMap AliasMap, filePath string) error {
	str := aliasMap.String()
	data := []byte(str)
	if err := ioutil.WriteFile(filePath, data, 0o644); err != nil {
		return err
	}

	return nil
}
