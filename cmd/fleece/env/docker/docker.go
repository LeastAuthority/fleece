package docker

import (
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/leastauthority/fleece/cmd/config"
)

var (
	CmdDocker = &cobra.Command{
		Use:   "docker",
		Short: "Manage docker image used for fuzzing",
		Args:  cobra.ExactArgs(0),
	}
)

func init() {
	CmdDocker.AddCommand(CmdBuild)
}

func getDockerDir() string {
	absoluteRepoRoot, err := filepath.Abs(viper.GetString(config.RepoRoot))
	if err != nil {
		panic(err)
	}
	//fmt.Printf("abs: %s\n", absoluteRepoRoot)
	return filepath.Join(absoluteRepoRoot, "fleece", dockerDirName)
}

func getDockerfilePath() string {
	return filepath.Join(getDockerDir(), defaultDockerfileName)
}
