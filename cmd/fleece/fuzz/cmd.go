package fuzz

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spf13/cobra"

	"github.com/leastauthority/fleece/docker"
)

var (
	CmdFuzz = &cobra.Command{
		Use:   "fuzz <pkg> <fuzz function>",
		Short: "Run go-fuzz against a fuzz function",
		Args:  cobra.ExactArgs(2),
		RunE:  runFuzz,
	}

	buildBin bool
	procs    int
)

func init() {
	CmdFuzz.Flags().BoolVarP(&buildBin, "build", "b", false, "if true, rebuilds test binary before running (default: false)")
	CmdFuzz.Flags().IntVarP(&procs, "procs", "p", 1, "number of processors to use (passed to go-fuzz's -procs flag)")
}

func runFuzz(cmd *cobra.Command, args []string) error {
	pkgPath := args[0]
	fuzzFuncName := args[1]

	fuzzCfg := docker.FuzzConfig{
		FuncName: fuzzFuncName,
		Build:    buildBin,
		Procs:    procs,
	}
	// TODO: use docker engine api!
	containerName, err := docker.RunGoFuzz(pkgPath, fuzzCfg)
	if err != nil {
		return err
	}

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, os.Kill)

	// TODO: use docker engine api
	if err := docker.LogGoFuzz(pkgPath, fuzzCfg); err != nil {
		return err
	}

	// NB: false hope that progress is happening
	fmt.Printf("Shutting down gracefully...")
	go func() {
		for range time.Tick(1 * time.Second) {
			fmt.Print(".")
		}
	}()

	// TODO: use docker engine api!
	if err := docker.StopContainer(containerName); err != nil {
		return fmt.Errorf("error encountered while stopping container: %w", err)
	}
	return nil
}
