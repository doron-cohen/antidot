package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/doron-cohen/antidot/tests"
	"github.com/stretchr/testify/require"
)

// loadSnapshot loads a snapshot file from testdata directory
func loadSnapshot(t *testing.T, snapshotName string) string {
	snapshotPath := filepath.Join("testdata", snapshotName)
	data, err := os.ReadFile(snapshotPath)
	require.NoError(t, err, "Failed to load snapshot %s", snapshotName)
	return string(data)
}

// normalizeOutput sorts lines to make output deterministic for comparison
func normalizeOutput(output string) string {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	// Sort lines to make comparison deterministic
	// Note: This assumes that order doesn't matter for the output
	// If order matters, we should not sort
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// normalizePaths replaces absolute paths with placeholders for comparison
func normalizePaths(output string) string {
	// Replace XDG fallback paths with standard values using regex
	// For bash/zsh: ${XDG_VAR:-/path/to/fallback}
	bashZshPattern := regexp.MustCompile(`\$\{XDG_CACHE_HOME:-[^}]+\}`)
	output = bashZshPattern.ReplaceAllString(output, "${XDG_CACHE_HOME:-/USER_HOME/.cache}")

	bashZshConfigPattern := regexp.MustCompile(`\$\{XDG_CONFIG_HOME:-[^}]+\}`)
	output = bashZshConfigPattern.ReplaceAllString(output, "${XDG_CONFIG_HOME:-/USER_HOME/.config}")

	bashZshDataPattern := regexp.MustCompile(`\$\{XDG_DATA_HOME:-[^}]+\}`)
	output = bashZshDataPattern.ReplaceAllString(output, "${XDG_DATA_HOME:-/USER_HOME/.local/share}")

	// For fish: set -x XDG_VAR "/path/to/fallback"
	fishCachePattern := regexp.MustCompile(`set -x XDG_CACHE_HOME "[^"]*"`)
	output = fishCachePattern.ReplaceAllString(output, "set -x XDG_CACHE_HOME \"/USER_HOME/.cache\"")

	fishConfigPattern := regexp.MustCompile(`set -x XDG_CONFIG_HOME "[^"]*"`)
	output = fishConfigPattern.ReplaceAllString(output, "set -x XDG_CONFIG_HOME \"/USER_HOME/.config\"")

	fishDataPattern := regexp.MustCompile(`set -x XDG_DATA_HOME "[^"]*"`)
	output = fishDataPattern.ReplaceAllString(output, "set -x XDG_DATA_HOME \"/USER_HOME/.local/share\"")

	// Replace any remaining temporary directory paths with placeholders
	output = strings.ReplaceAll(output, "/var/folders/", "/TEMP_DIR/")
	output = strings.ReplaceAll(output, "/tmp/", "/TEMP_DIR/")
	output = strings.ReplaceAll(output, "/Users/doron/", "/USER_HOME/")

	return output
}

func TestInitCommand(t *testing.T) {
	env := tests.SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("init without shell override", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init")

		// Pass the current environment to ensure XDG variables are available
		initCmd.Env = os.Environ()

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_bash_default.bash")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("init with bash shell", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "bash")

		// Pass the current environment to ensure XDG variables are available
		initCmd.Env = os.Environ()

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init --shell bash failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_bash_default.bash")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("init with fish shell", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "fish")

		// Pass the current environment to ensure XDG variables are available
		initCmd.Env = os.Environ()

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init --shell fish failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_fish_default.fish")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("init with zsh shell", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "zsh")

		// Pass the current environment to ensure XDG variables are available
		initCmd.Env = os.Environ()

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init --shell zsh failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_zsh_default.zsh")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("init with invalid shell", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "invalid_shell")
		_, err := initCmd.Output()
		require.Error(t, err, "antidot init with invalid shell should fail")

		// Check that the error message mentions supported shells
		if exitErr, ok := err.(*exec.ExitError); ok {
			errorOutput := string(exitErr.Stderr)
			require.Contains(t, errorOutput, "not supported", "Error should mention shell is not supported")
		}
	})
}

func TestInitWithKeyValueStore(t *testing.T) {
	env := tests.SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("init with existing key-value store", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_basic.json")
		require.NoError(t, err, "Failed to read kvstore_basic.json")
		tmpFile, err := os.CreateTemp("", "test_kvstore_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init with the custom key-value store
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "bash")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_bash_with_kvstore.bash")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("init with key-value store for fish shell", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_basic.json")
		require.NoError(t, err, "Failed to read kvstore_basic.json")
		tmpFile, err := os.CreateTemp("", "test_kvstore_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init with fish shell
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "fish")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_fish_with_kvstore.fish")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})
}

func TestInitCommandOutputConsistency(t *testing.T) {
	env := tests.SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("init output is consistent across runs", func(t *testing.T) {
		antidotPath := tests.GetAntidotPath(t)

		// Run antidot init multiple times
		initCmd1 := exec.Command(antidotPath, "init", "--shell", "bash")
		initCmd1.Env = os.Environ()
		output1, err := initCmd1.Output()
		require.NoError(t, err, "First antidot init failed")

		initCmd2 := exec.Command(antidotPath, "init", "--shell", "bash")
		initCmd2.Env = os.Environ()
		output2, err := initCmd2.Output()
		require.NoError(t, err, "Second antidot init failed")

		// The output should contain the same content, but order may vary due to map iteration
		output1Str := string(output1)
		output2Str := string(output2)

		// Verify both outputs contain the expected XDG exports
		require.Contains(t, output1Str, "XDG_CONFIG_HOME", "First output should contain XDG_CONFIG_HOME")
		require.Contains(t, output1Str, "XDG_CACHE_HOME", "First output should contain XDG_CACHE_HOME")
		require.Contains(t, output1Str, "XDG_DATA_HOME", "First output should contain XDG_DATA_HOME")

		require.Contains(t, output2Str, "XDG_CONFIG_HOME", "Second output should contain XDG_CONFIG_HOME")
		require.Contains(t, output2Str, "XDG_CACHE_HOME", "Second output should contain XDG_CACHE_HOME")
		require.Contains(t, output2Str, "XDG_DATA_HOME", "Second output should contain XDG_DATA_HOME")

		// Verify both outputs have the same length (same content, potentially different order)
		require.Equal(t, len(output1), len(output2), "Outputs should have the same length")

		// Verify normalized outputs are identical
		normalized1 := normalizeOutput(output1Str)
		normalized2 := normalizeOutput(output2Str)
		require.Equal(t, normalized1, normalized2, "Normalized outputs should be identical")
	})
}

func TestInitCommandFormatting(t *testing.T) {
	env := tests.SetupTestEnv(t)
	defer env.Cleanup()

	t.Run("bash environment variable formatting", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_formatting.json")
		require.NoError(t, err, "Failed to read kvstore_formatting.json")
		tmpFile, err := os.CreateTemp("", "test_formatting_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "bash")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_bash_formatting.bash")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("fish environment variable formatting", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_formatting.json")
		require.NoError(t, err, "Failed to read kvstore_formatting.json")
		tmpFile, err := os.CreateTemp("", "test_formatting_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "fish")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_fish_formatting.fish")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("bash alias formatting", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_formatting.json")
		require.NoError(t, err, "Failed to read kvstore_formatting.json")
		tmpFile, err := os.CreateTemp("", "test_formatting_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "bash")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_bash_formatting.bash")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})

	t.Run("fish alias formatting", func(t *testing.T) {
		// Load test data from file
		kvData, err := os.ReadFile("testdata/kvstore_formatting.json")
		require.NoError(t, err, "Failed to read kvstore_formatting.json")
		tmpFile, err := os.CreateTemp("", "test_formatting_*.json")
		require.NoError(t, err, "Failed to create temp file")
		defer os.Remove(tmpFile.Name())
		_, err = tmpFile.Write(kvData)
		require.NoError(t, err, "Failed to write test key-value store")
		tmpFile.Close()

		// Run antidot init
		antidotPath := tests.GetAntidotPath(t)
		initCmd := exec.Command(antidotPath, "init", "--shell", "fish")

		// Set the ANTIDOT_STATE_FILE environment variable
		initCmd.Env = append(os.Environ(), "ANTIDOT_STATE_FILE="+tmpFile.Name())

		output, err := initCmd.Output()
		require.NoError(t, err, "antidot init failed")

		outputStr := string(output)
		normalizedOutput := normalizeOutput(outputStr)
		normalizedOutput = normalizePaths(normalizedOutput)

		// Compare with snapshot
		expectedOutput := loadSnapshot(t, "init_fish_formatting.fish")
		expectedNormalized := normalizeOutput(expectedOutput)
		expectedNormalized = normalizePaths(expectedNormalized)
		require.Equal(t, expectedNormalized, normalizedOutput, "Output should match snapshot")
	})
}
