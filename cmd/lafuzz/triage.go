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
		Short: "tests known crashing inputs and prints a summary.",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runTriage(cmd, args); err != nil {
				fmt.Printf("%+v\n", err)
				return err
			}
			return nil
		},
	}

	failureLimit     int
)

func init() {
	cmdTriage.Flags().IntVar(&failureLimit, "failure-failureLimit", 1000, "maximun number of failing inputs before stopping and showing the summary")
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]
	testArgsStr := fmt.Sprintf(
		"test -tags gofuzz -v -run Test%s %s -args -failingLimit %d",
		fuzzFuncName, pkgName, failureLimit)
	testArgs := strings.Split(testArgsStr, " ")

	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	return testCmd.Run()
}
