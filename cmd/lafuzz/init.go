package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/lafuzz/cmd/config"
	"github.com/leastauthority/lafuzz/docker"
)

const (
	defaultRoot = "lafuzz"
)

var cmdInit = &cobra.Command{
	Use:   "init [output dir]",
	Short: "Initialize lafuzz into output dir; defaults to `$(pwd)/lafuzz`.",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(cmdInit, cmdTriage)
}

func runInit(cmd *cobra.Command, args []string) error {
	var outputRoot string
	if len(args) > 0 {
		outputRoot = args[0]
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		outputRoot = filepath.Join(pwd, defaultRoot)
	}

	if err := makeAllWorkdirsDir(outputRoot); err != nil {
		return err
	}

	if err := docker.RestoreAssets(outputRoot, "docker"); err != nil {
		return err
	}

	viper.Set(config.RepoRoot, outputRoot)
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	contextDir := filepath.Join(outputRoot, "docker")
	dockerfile := filepath.Join(contextDir, "go-fuzz.dockerfile")
	return buildDocker(contextDir, dockerfile)
}

func makeAllWorkdirsDir(outputRoot string) error {
	workdirs := filepath.Join(outputRoot, "workdirs")

	// NB: might need to be more permissive.
	if err := os.MkdirAll(workdirs, 0755); err != nil {
		return err
	}

	gitkeep := filepath.Join(workdirs, ".gitkeep")
	if err := ioutil.WriteFile(gitkeep, nil, 0644); err != nil {
		return err
	}
	return nil
}

func buildDocker(contextDir, dockerfile string) error {
	argsStr := fmt.Sprintf("build -t go-fuzz -f %s %s", dockerfile, contextDir)
	args := strings.Split(argsStr, " ")
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
