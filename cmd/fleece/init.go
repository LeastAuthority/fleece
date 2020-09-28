package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/fleece/cmd/fleece/config"
	"github.com/leastauthority/fleece/cmd/fleece/env"
)

const (
	defaultRoot = "fleece"
	gitIgnoreLines = `*.zip`
)

var (
	cmdInit = &cobra.Command{
		// NB: disabled until hacks that break support are removed.
		//Use:   "init [output-dir]",
		Use:   "init",
		Short: "Initialize fleece into a repo",
		Long:  "Copies supporting files into output-dir (default: $(pwd)/fleece) and adds config file (default: .fleece.yaml)",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runInit,
	}

	initEnv bool
)

func init() {
	cmdInit.Flags().BoolVarP(&initEnv, "env", "e", false, "if provided, also runs the equivalent of \"fleece env init\"")
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

	if err := env.RestoreBindata(outputRoot); err != nil {
		return err
	}

	// TODO: flags for these
	// NB: repo root is expected to be the parent of outputRoot.
	viper.Set(config.RepoRoot, filepath.Dir(outputRoot))
	viper.Set(config.FleeceDir, outputRoot)
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	fmt.Println("Repo initialized for use with fleece.")
	if initEnv {
		if err := env.CmdInit.RunE(cmd, nil); err != nil {
			return err
		}
	} else {
		fmt.Println("Next run `fleece env init`!")
	}
	return nil
}

func makeAllWorkdirsDir(outputRoot string) error {
	workdirs := filepath.Join(outputRoot, "workdirs")

	// NB: might need to be more permissive.
	if err := os.MkdirAll(workdirs, 0755); err != nil {
		return err
	}

	writeGitKeep(workdirs)
	writeGitIgnore(workdirs)

	return nil
}

func writeGitKeep(dir string) error {
	gitkeep := filepath.Join(dir, ".gitkeep")
	if err := ioutil.WriteFile(gitkeep, nil, 0644); err != nil {
		return err
	}
	return nil
}

func writeGitIgnore(dir string) error {
	gitignore := filepath.Join(dir, ".gitignore")
	if err := ioutil.WriteFile(gitignore, []byte(gitIgnoreLines), 0644); err != nil {
		return err
	}
	return nil
}
