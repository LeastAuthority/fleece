package docker

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/leastauthority/fleece/cmd/fleece/config"
)

type FuzzConfig struct {
	FuncName string
	Build    bool
	Procs    int
}

// TODO: use docker engine api
func RunGoFuzz(pkgPath string, cfg FuzzConfig) (string, error) {
	var build string
	if cfg.Build {
		build = "-b"
	}

	name := ContainerName(pkgPath, cfg.FuncName)
	workdir, err := GetGuestWorkdir(cfg.FuncName)
	if err != nil {
		return "", err
	}

	cmd := []string{
		pkgPath, cfg.FuncName, build,
		"--", "-procs", fmt.Sprint(cfg.Procs), "-workdir", workdir,
	}
	runArgs := []string{
		"--rm", "-d",
		"--name", name,
		"--entrypoint", "/go-fuzz.sh",
	}
	return name, RunGoFuzzContainer(cmd, runArgs...)
}

// TODO: use docker engine api
func RunGoFuzzContainer(cmd []string, runArgs ...string)  error {
	repoRoot, err := config.GetRepoRoot()
	if err != nil {
		return err
	}

	args := append([]string{
		"-v", fmt.Sprintf("%s:/tmp/fuzzing", repoRoot),
	}, runArgs...)
	args = append(args, "go-fuzz")
	args = append(args, cmd...)

	return RunContainer(repoRoot, args)
}

// TODO: use docker engine api
func LogGoFuzz(pkgPath string, cfg FuzzConfig) error {
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, os.Interrupt, os.Kill)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		name := ContainerName(pkgPath, cfg.FuncName)
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
	return nil
}
