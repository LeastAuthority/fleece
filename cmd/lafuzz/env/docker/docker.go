package docker

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/lafuzz/cmd/config"
)

var (
	CmdDocker = &cobra.Command{
		Use:   "docker",
		Short: "manage docker image used for fuzzing",
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	CmdDocker.AddCommand(CmdBuild)
}

func getDockerDir() string {
	return filepath.Join(viper.GetString(config.OutputRoot), dockerDirName)
}

func getDockerfilePath() string {
	return filepath.Join(getDockerDir(), defaultDockerfileName)
}
