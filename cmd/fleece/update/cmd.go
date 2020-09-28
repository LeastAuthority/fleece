package update

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/fleece/cmd/fleece/config"
	"github.com/leastauthority/fleece/cmd/fleece/env"
)

var (
	CmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "update fleece CLI binary and optionally, repo files",
		Args:  cobra.MaximumNArgs(1),
		RunE:  runUpdate,
	}

	updateRepo, noUpdateCLI bool
)

func init() {
	CmdUpdate.Flags().BoolVar(&updateRepo, "repo", false, "if true, updates repo files (restore bindata); prompts before overwrite)")
	CmdUpdate.Flags().BoolVar(&noUpdateCLI, "no-cli", false, "if true, doesn't update the CLI before updating repo files (implies --repo)")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	if !noUpdateCLI {
		if err := updateCLI(); err != nil {
			return err
		}
	}

	if updateRepo || noUpdateCLI {
		if err := updateFiles(); err != nil {
			return err
		}
	}
	return nil
}

func updateCLI() error {
	goArgs := strings.Split("get -u github.com/leastauthority/fleece/cmd/fleece", " ")
	goCmd := exec.Command("go", goArgs...)
	goCmd.Env = append(os.Environ(), "GO111MODULE=off")
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	return goCmd.Run()
}

func updateFiles() error {
	return env.RestoreBindata(viper.GetString(config.FleeceDir))
}
