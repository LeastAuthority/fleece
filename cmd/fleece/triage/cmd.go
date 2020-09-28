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

	crashLimit int
	pattern    string
)

func init() {
	CmdTriage.Flags().IntVar(&crashLimit, "crash-limit", 1000, "maximun number of failing inputs before stopping and showing the summary")
	//CmdTriage.Flags().StringVarP(&pattern, "pattern", "p", "", "string pattern to match in test output")
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
	}
	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	// NB: doesn't matter if the test fails.
	_ = testCmd.Run()
	return nil
}
