package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

func GetRepoRoot() (string, error) {
	repoRoot, err := filepath.Abs(viper.GetString(RepoRoot))
	if err != nil {
		return "", fmt.Errorf("unable to get path to repo root: %w", err)
	}
	return repoRoot, nil
}

func GetFleeceDir() (string, error) {
	repoRoot, err := GetRepoRoot()
	if err != nil {
		return "", fmt.Errorf("unable to get path to fleece dir: %w", err)
	}
	return filepath.Join(repoRoot, viper.GetString(FleeceDir)), nil
}

func GetRelativeFleeceDir() (string) {
	return filepath.Join(viper.GetString(RepoRoot), viper.GetString(FleeceDir))
}