package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"time"

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
	var build string
	if buildCorpus {
		build = "-b"
	}

	name := fmt.Sprintf("%s_%s", filepath.Base(pkgPath), fuzzFuncName)
	runArgs := []string{
		"--rm", "-d",
		"--name", name,
		"--entrypoint", "/go-fuzz.sh",
		"-v", fmt.Sprintf("%s:/tmp/fuzzing", repoRoot),
		"go-fuzz", pkgPath, fuzzFuncName, build, "--", "-procs=1",
	}

	if err := runContainer(repoRoot, runArgs); err != nil {
		return err
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, os.Kill)

	// Wait for interrupt / kill signal
	_ = <-sigC

	// NB: false hope that progress is happening
	fmt.Printf("Shutting down gracefully...")
	go func() {
		for range time.Tick(1 * time.Second) {
			fmt.Print(".")
		}
	}()

	if err := stopContainer(name); err != nil {
		return fmt.Errorf("error encountered while stopping container: %w", err)
	}
	return nil
}

func runContainer(cwd string, args []string) error {
	args = append([]string{"run"}, args...)
	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Dir = cwd
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}

func stopContainer(name string) error {
	cmd := exec.Command("docker", "stop", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
