package shell_test

import (
	"os"
	"strings"
	"testing"

	"github.com/doron-cohen/antidot/internal/shell"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFishShell(t *testing.T) {
	fish := &shell.Fish{}

	t.Run("RenderInit with empty kvstore", func(t *testing.T) {
		// Create a fresh kvstore for this test
		kv := &shell.KeyValueStore{
			EnvVars: make(map[string]string),
			Aliases: make(map[string]string),
		}

		result := fish.RenderInit(kv)

		// Should contain XDG exports with fish syntax
		assert.Contains(t, result, "set -q XDG_CONFIG_HOME")
		assert.Contains(t, result, "set -q XDG_CACHE_HOME")
		assert.Contains(t, result, "set -q XDG_DATA_HOME")

		// Should not contain any aliases or env vars since kvstore is empty
		assert.NotContains(t, result, "alias")
		assert.NotContains(t, result, "set -gx TEST_VAR")
	})

	t.Run("RenderInit with populated kvstore", func(t *testing.T) {
		// Create a temporary file for the kvstore with proper JSON content
		tmpFile, err := os.CreateTemp("", "antidot_test_*.json")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Write initial JSON content
		initialJSON := `{"env":{},"alias":{}}`
		_, err = tmpFile.WriteString(initialJSON)
		require.NoError(t, err)
		tmpFile.Close()

		// Create a fresh kvstore for this test
		kv, err := shell.LoadKeyValueStore(tmpFile.Name())
		require.NoError(t, err)

		// Add test data
		err = kv.AddEnv("TEST_VAR", "test_value")
		require.NoError(t, err)
		err = kv.AddEnv("ANOTHER_VAR", "${XDG_CONFIG_HOME}/test")
		require.NoError(t, err)
		err = kv.AddAlias("test_alias", "test_command --flag")
		require.NoError(t, err)
		err = kv.AddAlias("another_alias", "another_command ${XDG_DATA_HOME}/config")
		require.NoError(t, err)

		result := fish.RenderInit(kv)

		// Should contain XDG exports with fish syntax
		assert.Contains(t, result, "set -q XDG_CONFIG_HOME")
		assert.Contains(t, result, "set -q XDG_CACHE_HOME")
		assert.Contains(t, result, "set -q XDG_DATA_HOME")

		// Should contain the env vars with fish syntax
		assert.Contains(t, result, "set -gx TEST_VAR \"test_value\"")
		assert.Contains(t, result, "set -gx ANOTHER_VAR \"$XDG_CONFIG_HOME/test\"")

		// Should contain the aliases with fish syntax
		assert.Contains(t, result, "alias test_alias \"test_command --flag\"")
		assert.Contains(t, result, "alias another_alias \"another_command $XDG_DATA_HOME/config\"")

		// Verify order: XDG exports first, then env vars, then aliases
		lines := strings.Split(strings.TrimSpace(result), "\n")
		xdgCount := 0
		envCount := 0
		aliasCount := 0

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "set -q XDG_") {
				xdgCount++
			} else if strings.HasPrefix(line, "set -gx ") {
				envCount++
			} else if strings.HasPrefix(line, "alias ") {
				aliasCount++
			}
		}

		assert.Equal(t, 3, xdgCount, "Should have 3 XDG exports")
		assert.Equal(t, 2, envCount, "Should have 2 env vars")
		assert.Equal(t, 2, aliasCount, "Should have 2 aliases")

		// Verify XDG exports come before env vars, and env vars come before aliases
		xdgIndex := -1
		envIndex := -1
		aliasIndex := -1

		for i, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "set -q XDG_") {
				if xdgIndex == -1 {
					xdgIndex = i
				}
			} else if strings.HasPrefix(line, "set -gx ") {
				if envIndex == -1 {
					envIndex = i
				}
			} else if strings.HasPrefix(line, "alias ") {
				if aliasIndex == -1 {
					aliasIndex = i
				}
			}
		}

		assert.True(t, xdgIndex < envIndex, "XDG exports should come before env vars")
		assert.True(t, envIndex < aliasIndex, "Env vars should come before aliases")
	})

	t.Run("RenderInit with nil kvstore", func(t *testing.T) {
		result := fish.RenderInit(nil)

		// Should contain XDG exports with fish syntax
		assert.Contains(t, result, "set -q XDG_CONFIG_HOME")
		assert.Contains(t, result, "set -q XDG_CACHE_HOME")
		assert.Contains(t, result, "set -q XDG_DATA_HOME")

		// Should not contain any aliases or env vars
		assert.NotContains(t, result, "alias")
		assert.NotContains(t, result, "set -gx TEST_VAR")
	})

	t.Run("unbracketEnvVar functionality", func(t *testing.T) {
		// Test that ${VAR} gets converted to $VAR in fish syntax through RenderInit
		// Create a temporary file for the kvstore
		tmpFile, err := os.CreateTemp("", "antidot_test_*.json")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		// Write initial JSON content
		initialJSON := `{"env":{},"alias":{}}`
		_, err = tmpFile.WriteString(initialJSON)
		require.NoError(t, err)
		tmpFile.Close()

		// Create a fresh kvstore for this test
		kv, err := shell.LoadKeyValueStore(tmpFile.Name())
		require.NoError(t, err)

		// Add test data with bracketed variables
		err = kv.AddEnv("PATH", "${XDG_DATA_HOME}/bin:${PATH}")
		require.NoError(t, err)
		err = kv.AddAlias("grep", "grep --color=auto ${GREP_COLORS}")
		require.NoError(t, err)

		// Add test data with empty braces
		err = kv.AddEnv("SPECIAL_CHARS", "value with $ and {} and \"quotes\"")
		require.NoError(t, err)
		err = kv.AddAlias("special_alias", "command with $ and {} and \"quotes\"")
		require.NoError(t, err)

		result := fish.RenderInit(kv)

		// Should contain unbracketed variables
		assert.Contains(t, result, "set -gx PATH \"$XDG_DATA_HOME/bin:$PATH\"")
		assert.Contains(t, result, "alias grep \"grep --color=auto $GREP_COLORS\"")

		// Should contain single braces for empty braces
		assert.Contains(t, result, "set -gx SPECIAL_CHARS \"value with $ and { and \"quotes\"\"")
		assert.Contains(t, result, "alias special_alias \"command with $ and { and \"quotes\"\"")
	})
}
