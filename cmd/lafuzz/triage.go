package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var cmdTriage = &cobra.Command{
	Use:   "triage <pkg> <fuzz function>",
	Short: "tests known crashing inputs and prints a summary.",
	Args:  cobra.ExactArgs(2),
	RunE:  runTriage,
}

func init() {
	rootCmd.AddCommand(cmdTriage)
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]
	testArgsStr := fmt.Sprintf("test -tags gofuzz -v -run Test%s %s", fuzzFuncName, pkgName)
	testArgs := strings.Split(testArgsStr, " ")

	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	return testCmd.Run()
}
