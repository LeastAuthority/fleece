package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/lafuzz/cmd/config"
	"github.com/leastauthority/lafuzz/cmd/lafuzz/env"
	"github.com/leastauthority/lafuzz/docker"
)

const (
	defaultRoot = "lafuzz"
)

var (
	cmdInit = &cobra.Command{
		Use:   "init [output-dir]",
		Short: "initialize lafuzz into a repo",
		Long:  "copies supporting files into output-dir (default: $(pwd)/lafuzz) and adds config file (default: .lafuzz.yaml)",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runInit,
	}

	initEnv bool
)

func init() {
	cmdInit.Flags().BoolVarP(&initEnv, "env", "e", false, "if provided, also runs the equivalent of \"lafuzz env init\"")
}

func runInit(cmd *cobra.Command, args []string) error {
	var outputRoot string
	if len(args) > 0 {
		outputRoot = args[0]
	} else {
		outputRoot = filepath.Join(".", defaultRoot)
	}

	if err := makeAllWorkdirsDir(outputRoot); err != nil {
		return err
	}

	if err := docker.RestoreAssets(outputRoot, "docker"); err != nil {
		return err
	}

	// NB: repo root is expected to be the parent of outputRoot.
	viper.Set(config.RepoRoot, filepath.Dir(outputRoot))
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	fmt.Println("Repo initialized for use with lafuzz.")
	if initEnv {
		if err := env.CmdInit.RunE(cmd, nil); err != nil {
			return err
		}
	} else {
		fmt.Println("Next run `lafuzz env init`!")
	}
	return nil
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
