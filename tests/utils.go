package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/doron-cohen/antidot/internal/utils"
	"github.com/stretchr/testify/require"
)

// TestEnv manages test environment setup and cleanup
type TestEnv struct {
	OriginalAppName       string
	OriginalXdgDataHome   string
	OriginalXdgConfigHome string
	OriginalXdgCacheHome  string
	OriginalXdgStateHome  string
	OriginalXdgRuntimeDir string
}

// SetupTestEnv creates a test environment with isolated directories
func SetupTestEnv(t *testing.T) *TestEnv {
	// Create a temporary directory for test data
	tmpDir := t.TempDir()

	// Override the app directories to use our temp directory
	originalAppName := utils.AppDirs.AppName
	utils.AppDirs.AppName = "antidot_test"

	// Set all XDG variables to our temp directory to ensure complete isolation
	originalXdgDataHome := os.Getenv("XDG_DATA_HOME")
	originalXdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
	originalXdgCacheHome := os.Getenv("XDG_CACHE_HOME")
	originalXdgStateHome := os.Getenv("XDG_STATE_HOME")
	originalXdgRuntimeDir := os.Getenv("XDG_RUNTIME_DIR")

	os.Setenv("XDG_DATA_HOME", tmpDir)
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	os.Setenv("XDG_CACHE_HOME", tmpDir)
	os.Setenv("XDG_STATE_HOME", tmpDir)
	os.Setenv("XDG_RUNTIME_DIR", tmpDir)

	return &TestEnv{
		OriginalAppName:       originalAppName,
		OriginalXdgDataHome:   originalXdgDataHome,
		OriginalXdgConfigHome: originalXdgConfigHome,
		OriginalXdgCacheHome:  originalXdgCacheHome,
		OriginalXdgStateHome:  originalXdgStateHome,
		OriginalXdgRuntimeDir: originalXdgRuntimeDir,
	}
}

// Cleanup restores the original environment
func (te *TestEnv) Cleanup() {
	utils.AppDirs.AppName = te.OriginalAppName
	os.Setenv("XDG_DATA_HOME", te.OriginalXdgDataHome)
	os.Setenv("XDG_CONFIG_HOME", te.OriginalXdgConfigHome)
	os.Setenv("XDG_CACHE_HOME", te.OriginalXdgCacheHome)
	os.Setenv("XDG_STATE_HOME", te.OriginalXdgStateHome)
	os.Setenv("XDG_RUNTIME_DIR", te.OriginalXdgRuntimeDir)
}

// GetAntidotPath returns the absolute path to the antidot binary
func GetAntidotPath(t *testing.T) string {
	antidotPath, err := filepath.Abs("../../antidot")
	require.NoError(t, err, "Failed to get absolute path to antidot")
	return antidotPath
}
