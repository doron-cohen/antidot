package rules_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/dotfile"

	"github.com/doron-cohen/antidot/internal/rules"
	"github.com/google/go-cmp/cmp"
)

type configTest struct {
	Name                string
	ConfigYaml          string
	ExpectedRulesConfig rules.RulesConfig
}

func TestLoadConfig(t *testing.T) {
	tests := []configTest{
		{
			Name: "EmptyRules",
			ConfigYaml: `version: 1
rules: []`,
			ExpectedRulesConfig: rules.RulesConfig{Version: 1, Rules: []rules.Rule{}},
		},
		{
			Name: "AllActions",
			ConfigYaml: `version: 1
rules:
- name: all_actions
  description: a rule containing all actions
  dotfile:
    name: .test
    is_dir: true
  actions:
    - type: migrate
      source: source_path
      dest: dest_path
      symlink: true
    - type: delete
      path: delete_path`,
			ExpectedRulesConfig: rules.RulesConfig{
				Version: 1,
				Rules: []rules.Rule{
					{
						Name:        "all_actions",
						Description: "a rule containing all actions",
						Dotfile:     &dotfile.Dotfile{Name: ".test", IsDir: true},
						Ignore:      false,
						Actions: []rules.Action{
							rules.Migrate{
								Source:  "source_path",
								Dest:    "dest_path",
								Symlink: true,
							},
							rules.Delete{
								Path: "delete_path",
							},
						},
					},
				},
			},
		},
		{
			Name: "IgnoredRule",
			ConfigYaml: `version: 1
rules:
- name: ignored_rule
  description: an ignored rule
  dotfile:
    name: .test
    is_dir: false
  ignore: true`,
			ExpectedRulesConfig: rules.RulesConfig{
				Version: 1,
				Rules: []rules.Rule{
					{
						Name:        "ignored_rule",
						Description: "an ignored rule",
						Dotfile:     &dotfile.Dotfile{Name: ".test", IsDir: false},
						Ignore:      true,
						Actions:     nil,
					},
				},
			},
		},
	}

	for _, test := range tests {
		tmpDir, err := ioutil.TempDir("", "test")
		if err != nil {
			t.Errorf("Failed creating temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		t.Run(test.Name, func(t *testing.T) {
			configFileName := fmt.Sprintf("%s.yml", test.Name)
			configPath := filepath.Join(tmpDir, configFileName)
			configFile, err := os.Create(configPath)
			if err != nil {
				t.Errorf("Failed creating test config file: %v", err)
			}

			_, err = configFile.WriteString(test.ConfigYaml)
			if err != nil {
				t.Errorf("Failed writing test config to file: %v", err)
			}

			loaded, err := rules.LoadRulesConfig(configPath)
			if err != nil {
				t.Errorf("Failed loading test config from file: %v", err)
			}

			if !cmp.Equal(loaded, test.ExpectedRulesConfig) {
				t.Errorf("Found rules config are not what we expected:\n%v", cmp.Diff(loaded, test.ExpectedRulesConfig))
			}
		})
	}
}
