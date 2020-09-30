package triage

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"

	"github.com/leastauthority/fleece/cmd/fleece/config"
)

var (
	CmdTriage = &cobra.Command{
		Use:   "triage <pkg> <fuzz function>",
		Short: "test discovered crashing inputs and summarize",
		Args:  cobra.ExactArgs(2),
		RunE:  runTriage,
	}

	crashLimit           int
	pattern              string
	skipPattern          string
	skipPatternDelimiter string
	safe, verbose        bool
)

func init() {
	CmdTriage.Flags().IntVar(&crashLimit, "crash-limit", 1000, "maximun number of failing inputs before stopping and showing the summary")
	//CmdTriage.Flags().StringVarP(&pattern, "pattern", "p", "", "string pattern to match in test output")
	CmdTriage.Flags().StringVarP(&skipPattern, "skip", "s", "", "skips inputs that have recorded outputs which match the skip pattern")
	// TODO: regex instead?  -- "...then you have two problems."
	CmdTriage.Flags().StringVarP(&skipPatternDelimiter, "skip-delimiter", "d", "", "if provided, used as delimiter to split skip pattern for matching multiple patterns (default: \"\")")
	CmdTriage.Flags().BoolVarP(&safe, "safe", "S", true, "if true, skips crashers with recorded outputs that timed-out or ran out of memory (default: true)")
	CmdTriage.Flags().BoolVarP(&verbose, "verbose", "v", false, "if true, logs each skip")
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]

	fleeceDir, err := config.GetFleeceDir()
	if err != nil {
		return err
	}

	testArgs := []string{
		"test",
		"-tags", "gofuzz",
		"-v",
		"-run", fmt.Sprintf("Test%s", fuzzFuncName),
		pkgName, "-args",
		"-crash-limit", fmt.Sprintf("%d", crashLimit),
		"-fleece-dir", fleeceDir,
		"-skip", skipPattern,
		"-skip-delimiter", skipPatternDelimiter,
	}
	if !safe {
		testArgs = append(testArgs, "--safe", "false")
	}
	if verbose {
		testArgs = append(testArgs, "-verbose")
	}

	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	// NB: doesn't matter if the test fails.
	_ = testCmd.Run()
	return nil
}
