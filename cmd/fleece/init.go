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

const gitIgnoreLines = `*.zip`

var (
	cmdInit = &cobra.Command{
		Use:   "init",
		Short: "Initialize fleece into a repo",
		Long:  "Copies supporting files into output-dir (default: $(pwd)/fleece) and adds config file (default: .fleece.yaml)",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runInit,
	}

	initEnv   bool
	repoRoot  string
	fleeceDir string
)

func init() {
	cmdInit.Flags().BoolVarP(&initEnv, "env", "e", false, "if provided, also runs the equivalent of \"fleece env init\"")
	cmdInit.Flags().StringVar(&repoRoot, "repo-root", ".", "path to the repo/module root relative to the config file (default: .)")
	cmdInit.Flags().StringVar(&fleeceDir, "fleece-dir", "fleece", "path to the fleece dir relative to the repo/module root (default: fleece)")
}

func runInit(cmd *cobra.Command, args []string) error {
	viper.Set(config.RepoRoot, repoRoot)
	viper.Set(config.FleeceDir, fleeceDir)
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	if err := makeAllWorkdirsDir(fleeceDir); err != nil {
		return err
	}

	if err := env.RestoreBindata(fleeceDir); err != nil {
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
