package triage

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var (
	CmdTriage = &cobra.Command{
		Use:   "triage",
		Short: "Browse, test and report on crashing in/outputs",
	}

	cmdTest = &cobra.Command{
		Use:   "test <pkg> <fuzz function>",
		Short: "Test crashers and summarize",
		Args:  cobra.ExactArgs(2),
		RunE: runTriage,
	}
	//cmdReport

	crashLimit int
	pattern string
)

func init() {
	CmdTriage.AddCommand(cmdTest)

	cmdTest.Flags().IntVar(&crashLimit, "crash-limit", 1000, "maximun number of failing inputs before stopping and showing the summary")
	cmdTest.Flags().StringVarP(&pattern, "pattern", "p", "", "string pattern to match in test output")
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]

	// count crashers (iteratively)
	//	 - passing | failing | total/limit
	//   - printable
	// call/recover target
	// compare panic msgs
	// show "current" crasher
	//   - printable
	// skip to next crasher
	// skip with pattern
	// return to previous crasher
	// find/count with pattern
	// unique
	//   - printable

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
