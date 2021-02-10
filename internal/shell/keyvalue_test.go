package shell_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/google/go-cmp/cmp"
)

func TestKeyValueStore(t *testing.T) {
	tmpDir := os.TempDir()
	kvPath := filepath.Join(tmpDir, "antidot_kvstore_test.json")
	defer os.Remove(kvPath)

	t.Run("Open a new file", func(t *testing.T) {
		kv, err := shell.LoadKeyValueStore(kvPath)
		if err != nil {
			t.Fatalf("Failed to open new key value store in %s: %v", kvPath, err)
		}

		if kv.Aliases == nil || len(kv.Aliases) > 0 {
			t.Fatalf("Aliases not initialized: %v", err)
		}

		if kv.EnvVars == nil || len(kv.EnvVars) > 0 {
			t.Fatalf("EnvVars not initialized: %v", err)
		}
	})

	t.Run("Add aliases and env vars", func(t *testing.T) {
		kv, err := shell.LoadKeyValueStore(kvPath)
		if err != nil {
			t.Fatalf("Failed to open existing key value store in %s: %v", kvPath, err)
		}

		aliases := map[string]string{
			"test":      "test --color",
			"tool":      "tool --config ${XDG_CONFIG_HOME}/tool.toml",
			"grrrrreat": "grrrrrrrrreat",
		}
		for k, v := range aliases {
			err = kv.AddAlias(k, v)
			if err != nil {
				t.Fatalf("Failed to add alias %s=%s: %v", k, v, err)
			}
		}

		envVars := map[string]string{
			"TEST_HOME": "${XDG_CONFIG_HOME}/test",
			"MEMES":     "${XDG_MEMES_HOME}",
			"SHELL":     "/bin/bash",
		}
		for k, v := range envVars {
			err = kv.AddEnv(k, v)
			if err != nil {
				t.Fatalf("Failed to add env %s=%s: %v", k, v, err)
			}
		}

		storedAliases, err := kv.ListAliases()
		if err != nil {
			t.Fatalf("Failed listing aliases: %v", err)
		}
		if !cmp.Equal(aliases, storedAliases) {
			t.Fatal("Aliases mismatch")
		}

		storedEnvVars, err := kv.ListEnvVars()
		if err != nil {
			t.Fatalf("Failed listing env vars: %v", err)
		}
		if !cmp.Equal(envVars, storedEnvVars) {
			t.Fatal("Env vars mismatch")
		}
	})
}
