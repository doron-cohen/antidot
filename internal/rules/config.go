package rules

import (
	"io/ioutil"
	"log"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)

type RulesConfig struct {
	Version int
	Rules   []Rule
}

var rulesConfig RulesConfig

func LoadRulesConfig(filepath string) error {
	log.Printf("Loading rules config file %s", filepath)
	rulesBytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	var rawConfig map[string]interface{}
	err = yaml.Unmarshal(rulesBytes, &rawConfig)
	if err != nil {
		return err
	}

	config := &mapstructure.DecoderConfig{
		DecodeHook: actionDecodeHook,
		Result:     &rulesConfig,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	err = decoder.Decode(rawConfig)
	if err != nil {
		return err
	}

	log.Printf("Loaded %d rules", len(rulesConfig.Rules))
	return nil
}
