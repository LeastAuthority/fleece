package fuzzing

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// CrasherIterator is an iterator for go-fuzz "crashers" located in the
//	respective fuzz function's working directory.
type CrasherIterator struct {
	i          int
	infos      []os.FileInfo
	fuzzFunc   Func
	crasherDir string
}

// NewCrasherIteratorFor returns an iterator for crashers that lazily loads	their inputs and outputs.
func NewCrasherIterator(env *Env, fuzzFunc Func) (*CrasherIterator, error) {
	crasherDir := env.GetCrasherDir(fuzzFunc)
	crasherInfos, err := ioutil.ReadDir(crasherDir)
	if err != nil {
		return nil, err
	}
	return &CrasherIterator{
		infos:      crasherInfos,
		crasherDir: crasherDir,
		fuzzFunc:   fuzzFunc,
	}, nil
}

// MustNewCrasherIterator returns an iterator for crashers but panics if an error occurs.
func MustNewCrasherIterator(env *Env, fuzzFunc Func) *CrasherIterator {
	iter, err := NewCrasherIterator(env, fuzzFunc)
	if err != nil {
		panic(err)
	}
	return iter
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
			done = iter.i == len(iter.infos)-1
			continue
		}

		input, err := ioutil.ReadFile(filepath.Join(iter.crasherDir, name))
		if err != nil {
			return nil, done, err
		}

		//output, err := ioutil.ReadFile(filepath.Join(iter.crasherDir, name) + ".output")
		//if err != nil {
		//	return nil, done, err
		//}

		next = &Crasher{
			Name:     name,
			Input:    input,
			//Output:   string(output),
			FuzzFunc: iter.fuzzFunc,
		}

		iter.i++
		done = iter.i == len(iter.infos)-1
		break
	}
	return next, done, nil
}

// TestFailingLimit tests each crasher's input against its respective fuzz
//	function until it sees `limit` failing inputs
func (iter CrasherIterator) TestFailingLimit(t *testing.T, env *Env, limit int) (_ *Crasher, panics int, total int) {
	crasherIterator, err := NewCrasherIterator(env, iter.fuzzFunc)
	require.NoError(t, err)

	// TODO: parallelize
	var done, didPanic bool
	var firstCrasher, crasher *Crasher
	var firstPanicMsg string
	for panics < limit {
		crasher, done, err = crasherIterator.Next()
		require.NoError(t, err)
		if done {
			break
		}

		didPanic = false
		crasher.Test(func(panicMsg string) {
			didPanic = true
			if firstCrasher == nil {
				firstPanicMsg = "panic: " + panicMsg
				firstCrasher = crasher
			}
		})
		if didPanic {
			panics++
		}
		total++
	}

	if firstCrasher != nil {
		fmt.Printf("Current panicking crasher: %s\n", firstCrasher.Name)
		fmt.Println("")
		fmt.Printf("Current panic message:\n%s\n", firstPanicMsg)
		fmt.Println("")
	}

	fmt.Printf("Crasher summary:\n===============\n")
	fmt.Printf("- passing: %d\n", total-panics)
	fmt.Printf("- failing: %d\n", panics)
	fmt.Printf("- total tested: %d\n", total)
	return firstCrasher, panics, total
}
