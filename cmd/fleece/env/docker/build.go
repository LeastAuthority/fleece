package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	dockerDirName         = "docker"
	defaultDockerfileName = "go-fuzz.dockerfile"
)

var (
	CmdBuild = &cobra.Command{
		Use:   "build",
		Short: "builds docker image used by the \"fuzz\" command",
		Args:  cobra.ExactArgs(0),
		RunE:  runBuild,
	}

	noCache bool
)

func init() {
	CmdBuild.Flags().BoolVar(&noCache, "no-cache", false, "passes --no-cache to docker build")
}

func runBuild(cmd *cobra.Command, args []string) error {
	var dockerArgs []string
	if noCache {
		dockerArgs = append(dockerArgs, "--no-cache")
	}
	return buildDocker(getDockerDir(), getDockerfilePath(), dockerArgs...)
}

func buildDocker(contextDir, dockerfile string, additionalArgs ...string) error {
	argsStr := fmt.Sprintf("build -t go-fuzz -f %s %s", dockerfile, contextDir)
	args := append(strings.Split(argsStr, " "), additionalArgs...)
	cmd := exec.Command("docker", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
