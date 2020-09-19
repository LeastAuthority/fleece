package main

import (
  "github.com/leastauthority/lafuzz/cmd/lafuzz/env"
  "github.com/spf13/cobra"
)

func init() {
  cmds := []*cobra.Command{
    cmdInit,
    cmdUpdate,
    cmdFuzz,
    cmdTriage,
    env.CmdEnv,
  }
  rootCmd.AddCommand(cmds...)
}

func main() {
  Execute()
}
