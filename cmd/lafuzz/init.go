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

var (
	cmdInit = &cobra.Command{
		Use:   "init [output dir]",
		Short: "Initialize lafuzz into output dir (default: `$(pwd)/lafuzz`). Also builds go-fuzz docker image.",
		RunE:  runInit,
	}

	noCache bool
)

func init() {
	rootCmd.AddCommand(cmdInit, cmdTriage)
	cmdInit.Flags().BoolVar(&noCache, "no-cache", false, "passes --no-cache to docker build")
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

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
	}
	// NB: repo root is expected to be the parent of outputRoot.
	viper.Set(config.RepoRoot, filepath.Dir(outputRoot))
	if err := viper.SafeWriteConfig(); err != nil {
		return err
	}

	contextDir := filepath.Join(outputRoot, "docker")
	dockerfile := filepath.Join(contextDir, "go-fuzz.dockerfile")
	var dockerArgs []string
	if noCache {
		dockerArgs = append(dockerArgs, "--no-cache")
	}
	return buildDocker(contextDir, dockerfile, dockerArgs...)
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

// TODO: doesn't detect fs changes
func buildDocker(contextDir, dockerfile string, additionalArgs ...string) error {
	argsStr := fmt.Sprintf("build -t go-fuzz -f %s %s", dockerfile, contextDir)
	args := append(strings.Split(argsStr, " "), additionalArgs...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
