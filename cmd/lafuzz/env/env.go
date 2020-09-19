package env

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/leastauthority/lafuzz/cmd/lafuzz/env/docker"
)

var (
	CmdEnv = &cobra.Command{
		Use:   "env",
		Args: cobra.ExactArgs(0),
		Short: "manage local fuzzing environment",
	}
	CmdInit = &cobra.Command{
		Use: "init",
		Short: "initialize local fuzzing environment for first time use",
		Args: cobra.ExactArgs(0),
		RunE: runInit,
	}
)

func init() {
	CmdEnv.AddCommand(CmdInit)
	CmdEnv.AddCommand(docker.CmdDocker)
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Println("Initializing local fuzzing environment")
	return docker.CmdBuild.RunE(cmd, args)
}