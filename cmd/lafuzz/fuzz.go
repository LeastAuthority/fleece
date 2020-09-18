package main

import (
	"context"
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
		Use:   "fuzz <pkg> <fuzz function>",
		Short: "start a fuzzing container for the specified fuzz function (blocking)",
		Args:  cobra.ExactArgs(2),
		RunE:  runFuzz,
	}

	buildCorpus bool
	procs       int
)

func init() {
	cmdFuzz.Flags().BoolVar(&buildCorpus, "build-corpus", false, "if true, builds corpus before running (default: false)")
	cmdFuzz.Flags().IntVar(&procs, "procs", 1, "number of processors to use (passed to go-fuzz's -procs flag)")
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

	// TODO: factor out
	name := containerName(pkgPath, fuzzFuncName)
	runArgs := []string{
		"--rm", "-d",
		"--name", name,
		"--entrypoint", "/go-fuzz.sh",
		"-v", fmt.Sprintf("%s:/tmp/fuzzing", repoRoot),
		"go-fuzz", pkgPath, fuzzFuncName, build, "--", "-procs", fmt.Sprint(procs),
	}

	if err := runContainer(repoRoot, runArgs); err != nil {
		return err
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	// TODO: use docker engine api
	go func() {
		logArgs := []string{"logs", "-f", name}
		logCmd := exec.Command("docker", logArgs...)
		logCmd.Stdout = os.Stdout
		logCmd.Stderr = os.Stderr
		if err := logCmd.Start(); err != nil {
			panic(err)
		}
		for {
			select {
			case <-ctx.Done():
				if err := logCmd.Process.Kill(); err != nil {
					panic(err)
				}
			default:
			}
		}
	}()

	// Wait for interrupt / kill signal
	_ = <-sigC
	cancel()

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

func containerName(pkgPath, fuzzFuncName string) string {
	return fmt.Sprintf("%s_%s", filepath.Base(pkgPath), fuzzFuncName)
}

// TODO: move to docker package
func runContainer(cwd string, args []string) error {
	args = append([]string{"run"}, args...)
	dockerCmd := exec.Command("docker", args...)
	dockerCmd.Dir = cwd
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	return dockerCmd.Run()
}

// TODO: move to docker package
func stopContainer(name string) error {
	cmd := exec.Command("docker", "stop", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
