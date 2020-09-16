package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/lafuzz/cmd/config"
)

var (
	cmdFuzz = &cobra.Command{
		Use:  "fuzz <pkg> <fuzz function>",
		Args: cobra.ExactArgs(2),
		RunE: runFuzz,
	}

	buildCorpus bool
)

func init() {
	rootCmd.AddCommand(cmdFuzz)
	cmdFuzz.Flags().BoolVar(&buildCorpus, "build-corpus", false, "if true, builds corpus before running (default: false)")
}

// TODO: stop container on exit!
func runFuzz(cmd *cobra.Command, args []string) error {
	pkgPath := args[0]
	fuzzFuncName := args[1]
	repoRoot := viper.GetString(config.RepoRoot)
	var build  string
	if buildCorpus {
		build = "-b"
	}

	dockerArgs := strings.Split(
		fmt.Sprintf("run --rm "+
			"--entrypoint /go-fuzz.sh "+
			"-v %s:/tmp/fuzzing "+
			"go-fuzz %s %s %s -- -procs=1 ",
			repoRoot, pkgPath, fuzzFuncName, build,
		), " ")
	dockerCmd := exec.Command("docker", dockerArgs...)
	dockerCmd.Dir = repoRoot
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}
