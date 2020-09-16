package fuzzing

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"golang.org/x/tools/go/packages"

	"github.com/stretchr/testify/require"
)

// Fuzz constants for go-fuzz to use when returning from the Fuzz func
var (
	FuzzInteresting = 1
	FuzzNormal      = 0
	FuzzDiscard     = -1
)

// Crasher represents a go-fuzz "crasher" (an input that crashed the respective
//	fuzz function), its corresponding output (panic message), and name (input hash).
type Crasher struct {
	Name     string
	Input    []byte
	Output   string
	FuzzFunc FuzzFunc
}

// CrasherIterator is an iterator for go-fuzz "crashers" located in the
//	respective fuzz function's working directory.
type CrasherIterator struct {
	i          int
	infos      []os.FileInfo
	fuzzFunc   FuzzFunc
	crasherDir string
}

type FuzzFunc func([]byte) int
type RecoverCallback func(panicMsg string)

// NewCrasherIteratorFor returns an iterator for crashers that lazily loads	their inputs and outputs.
func NewCrasherItertor(fuzzFunc FuzzFunc) (*CrasherIterator, error) {
	name := getFuncName(fuzzFunc)

	workdir, err := getWorkdir(name)
	if err != nil {
		return nil, err
	}

	crashersDir := filepath.Join(workdir, "crashers")
	crasherInfos, err := ioutil.ReadDir(crashersDir)
	if err != nil {
		return nil, err
	}
	return &CrasherIterator{
		infos:      crasherInfos,
		crasherDir: filepath.Join(workdir, "crashers"),
		fuzzFunc:   fuzzFunc,
	}, nil
}

// MustNewCrasherIterator returns an iterator for crashers but panics if an error occurs.
func MustNewCrasherIterator(fuzzFunc FuzzFunc) *CrasherIterator {
	iter, err := NewCrasherItertor(fuzzFunc)
	if err != nil {
		panic(err)
	}
	return iter
}

// Recover is intended to be deferred. It calls the recover callback with the
//	string representation of the recovered value in the event of a panic.
func (crasher *Crasher) Recover(recoverCb RecoverCallback) {
	if r := recover(); r != nil {
		recoverCb(fmt.Sprint(r))
	}
}

// Test calls the crashers fuzz function its input and recovers from panics
//	with the passed recover callback.
func (crasher *Crasher) Test(recoverCb RecoverCallback) {
	defer crasher.Recover(recoverCb)
	crasher.FuzzFunc(crasher.Input)
}

// Next gets the next crasher and increments the iterator.
func (iter *CrasherIterator) Next() (next *Crasher, done bool, err error) {
	// NB: iter.infos contains input, output, and quoted input files.
	//	We want to skip non-input files.
	for !done {
		info := iter.infos[iter.i]
		name := info.Name()
		if info.IsDir() || filepath.Ext(name) != "" {
			iter.i++
			done = iter.i == len(iter.infos)-2
			continue
		}

		input, err := ioutil.ReadFile(filepath.Join(iter.crasherDir, name))
		if err != nil {
			return nil, done, err
		}

		output, err := ioutil.ReadFile(filepath.Join(iter.crasherDir, name) + ".output")
		if err != nil {
			return nil, done, err
		}

		next = &Crasher{
			Name:     name,
			Input:    input,
			Output:   string(output),
			FuzzFunc: iter.fuzzFunc,
		}

		iter.i++
		done = iter.i == len(iter.infos)-2
		break
	}
	return next, done, nil
}

// TestFailingLimit tests each crasher's input against its respective fuzz
//	function until it sees `limit` failing inputs
func (iter CrasherIterator) TestFailingLimit(t *testing.T, limit int) (_ *Crasher, panics int, total int) {
	crasherIterator, err := NewCrasherItertor(iter.fuzzFunc)
	require.NoError(t, err)

	var done, didPanic bool
	var firstCrasher, crasher *Crasher
	for !done && panics < limit {
		crasher, done, err = crasherIterator.Next()
		require.NoError(t, err)

		didPanic = false
		crasher.Test(func(panicMsg string) {
			didPanic = true
			if firstCrasher == nil && crasher != nil {
				firstCrasher = crasher
			}
		})
		if didPanic {
			panics++
		}
		total++
	}

	fmt.Printf("Crasher summary:\n===============\n")
	fmt.Printf("- passing: %d\n", total-panics)
	fmt.Printf("- failing: %d\n", panics)
	fmt.Printf("- total: %d\n", total)

	if firstCrasher != nil {
		fmt.Println("")
		fmt.Printf("Next panicking crasher: %s\n%s\n", firstCrasher.Name, firstCrasher.Output)
	}
	return firstCrasher, panics, total
}

func getWorkdir(name string) (string, error) {
	pkgPath := reflect.TypeOf(FuzzNormal).PkgPath()
	modPath, err := getModPath(pkgPath)
	if err != nil {
		return "", err
	}

	return filepath.Join(modPath, "fuzzing", "workdirs", name), nil
}

func getModPath(importPath string) (string, error) {
	mode := packages.NeedName | packages.NeedModule
	pkgConfig := &packages.Config{Mode: mode}
	pkgs, err := packages.Load(pkgConfig, importPath)
	if err != nil {
		return "", err
	}
	if len(pkgs) > 1 {
		return "", errors.New(fmt.Sprintf(
			"more than one package matched pattern %q",
			importPath,
		))
	}

	pkg := pkgs[0]
	if pkg.Errors != nil || len(pkg.Errors) != 0 {
		errCount := packages.PrintErrors(pkgs)
		return "", errors.New(fmt.Sprintf(
			"%d errors encountered while loading package with import path %q",
			errCount, importPath,
		))
	}

	mod := pkg.Module
	if mod == nil {
		return "", errors.New(fmt.Sprintf(
			"unable to load module for package at import path %q",
			importPath,
		))
	}

	return mod.Dir, nil
}

func getFuncName(f FuzzFunc) string {
	return filepath.Ext(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name())[1:]
}
