package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	cmdTriage = &cobra.Command{
		Use:   "triage <pkg> <fuzz function>",
		Short: "Test crashers and summarize",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runTriage(cmd, args); err != nil {
				return err
			}
			return nil
		},
	}

	crashLimit int
)

func init() {
	cmdTriage.Flags().IntVar(&crashLimit, "crash-limit", 1000, "maximun number of failing inputs before stopping and showing the summary")
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]
	testArgsStr := fmt.Sprintf(
		"test -tags gofuzz -v -run Test%s %s -args -crash-limit %d",
		fuzzFuncName, pkgName, crashLimit)
	testArgs := strings.Split(testArgsStr, " ")

	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	// NB: doesn't matter if the test fails.
	_ = testCmd.Run()
	return nil
}
