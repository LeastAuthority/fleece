package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetFleeceDir() (string, error) {
	fleeceDir, err := filepath.Abs(filepath.Join(
		viper.GetString(RepoRoot),
		viper.GetString(FleeceDir),
	))
	if err != nil {
		return "", fmt.Errorf("unable to get path to fleece dir: %w", err)
	}
	return fleeceDir, nil
}
