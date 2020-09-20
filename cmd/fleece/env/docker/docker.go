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
	absoluteOutputRoot, err := filepath.Abs(viper.GetString(config.RepoRoot))
	if err != nil {
		panic(err)
	}
	//fmt.Printf("abs: %s\n", absoluteOutputRoot)
	return filepath.Join(absoluteOutputRoot, dockerDirName)
}

func getDockerfilePath() string {
	return filepath.Join(getDockerDir(), defaultDockerfileName)
}
