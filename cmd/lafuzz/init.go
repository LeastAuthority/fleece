package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/cobra"
)

const defaultRoot = "lafuzz"

var (
	cmdInit = &cobra.Command{
		Use:  "init [output dir]",
		Short: "Initialize lafuzz into output dir; defaults to `$(pwd)/lafuzz`.",
		RunE: runInit,
	}
	//cmdFuzz = &cobra.Command{
	//	Use: "fuzz <pkg> <fuzz function>,
	//  Args: cobra.ExactArgs(2),
	//  RunE: runFuzz,
	//}
	cmdTriage = &cobra.Command{
		Use:  "triage <pkg> <fuzz function>",
		Short: "Triage tests known crashing inputs and prints a summary.",
		Args: cobra.ExactArgs(2),
		RunE: runTriage,
	}
)

func init() {
	rootCmd.AddCommand(cmdInit, cmdTriage)
}

func runInit(cmd *cobra.Command, args []string) error {
	var outputRoot, workdirs string
	if len(args) > 0 {
		outputRoot = args[0]
	} else {
		pwd, err := os.Getwd()
		if err != nil {
			return err
		}
		outputRoot = filepath.Join(pwd, defaultRoot)
	}
	workdirs = filepath.Join(outputRoot, "workdirs")

	// NB: might need to be more permissive.
	if err := os.MkdirAll(workdirs, 0755); err != nil {
		return err
	}

	gitkeep := filepath.Join(workdirs, ".gitkeep")
	if err := ioutil.WriteFile(gitkeep, nil, 0644); err != nil {
		return err
	}

	pkgPath := fmt.Sprintf("%s", getPkgPath())
	goCmd := exec.Command("go", "get", pkgPath)
	goCmd.Env = []string{"GO111MDOULE=off"}
	if stdout, err := goCmd.CombinedOutput(); err != nil {
		fmt.Print(string(stdout))
		return err
	}
	return nil
}

func runTriage(cmd *cobra.Command, args []string) error {
	pkgName := args[0]
	fuzzFuncName := args[1]
	testArgsStr := fmt.Sprintf("test -v -run Test%s %s", fuzzFuncName, pkgName)
	testArgs := strings.Split(testArgsStr, " ")

	testCmd := exec.Command("go", testArgs...)
	testCmd.Stdout = os.Stdout
	testCmd.Stderr = os.Stderr

	return testCmd.Run()
}

func getPkgPath() string {
	return reflect.ValueOf(getPkgPath).Type().PkgPath()
}
