package main

import (
	"github.com/leastauthority/fleece/cmd/flags"
	"github.com/leastauthority/fleece/cmd/fleece/env"
	"github.com/leastauthority/fleece/cmd/fleece/fuzz"
	"github.com/leastauthority/fleece/cmd/fleece/triage"
	"github.com/leastauthority/fleece/cmd/fleece/update"
	"github.com/spf13/cobra"
)

func init() {
	cmds := []*cobra.Command{
		cmdInit,
		update.CmdUpdate,
		fuzz.CmdFuzz,
		triage.CmdTriage,
		env.CmdEnv,
	}
	rootCmd.AddCommand(cmds...)
	rootCmd.PersistentFlags().BoolVarP(&flags.Interactive, "interactive", "i", false, "if true, a term-ui is presented where available")
}

func main() {
	Execute()
}
