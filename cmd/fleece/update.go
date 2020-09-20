package main

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cmdUpdate = &cobra.Command{
		Use:   "update",
		Short: "update fleece CLI binary using \"go get -u\"",
		RunE:  runUpdate,
	}
)

func runUpdate(cmd *cobra.Command, agrs []string) error {
	goArgs := strings.Split("get -u github.com/leastauthority/fleece/cmd/fleece", " ")
	goCmd := exec.Command("go", goArgs...)
	goCmd.Env = append(os.Environ(), "GO111MODULE=off")
	goCmd.Stderr = os.Stderr
	goCmd.Stdout = os.Stdout
	return goCmd.Run()
}
