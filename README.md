# Fleece
## Value proposition
Simple to use and extensible tool that facilitates a fuzz -> triage -> debug workflow with mininal impact on the target repo.

## How to use
#### Setup and help
```bash
# Install fleece CLI binary:
GO111MODULE=off go get -u github.com/leastauthority/fleece/cmd/fleece

# Navigate to the root of the repo you'll be fuzzing:
cd ./path/to/repo/root

# Init fleece files in repo:
fleece init

# Init local fuzzing environment
fleece env init

# OR do both at the same time
# fleece init --env
```
```
Fuzz -> Triage -> Debug -> Repeat

Usage:
  fleece [command]

Available Commands:
  env         manage local fuzzing environment
  fuzz        run go-fuzz against a fuzz function
  help        Help about any command
  init        initialize fleece into a repo
  triage      test crashers and summarize
  update      update fleece CLI binary using "go get -u"

Flags:
      --config string   config file (default is $(pwd)/.fleece.yaml)
  -h, --help            help for fleece

Use "fleece [command] --help" for more information about a command.
```

#### Fuzzing
Now fuzz functions need to be defined.
The current convention is that they should be in or around the same package as the code they exercise.
Use the `gofuzz` build tag to exclude this code from normal builds and tests.

Fuzz functions follow the `go-fuzz` signature:
```golang
# see ./example/example_fuzz.go

//+build gofuzz

package example

func FuzzBuggyFunc(data []byte) int {
    // ...
}
```
_(see: [go-fuzz readme](https://github.com/dvyukov/go-fuzz/blob/master/README.md) for more details)_

To run your fuzz functions:
```bash
fleece fuzz ./example FuzzBuggyFunc --procs 1
# see fleece fuzz --help
```

#### Triaging
Once you've discovered some crashing inputs you can look through their stack traces to debug them.
It's possible that over the course of fuzzing you may discover many thousands of crashing inputs.
We're going to tackle these one at a time.
It's possible that multiple inputs have a common bug.

To debug, add a "triage test" which will run the crashing inputs back through the fuzz function that produced them to see if they're still crashing.
This also has the benefit of acting as a regression test when combined with a with a simple assertion at all inputs are no longer crashing, and retaining the crashers in version control.
Again, it's probably best to se the `gofuzz` build tag to exclude this code from normal builds and tests.

Here's the triage test corresponding with our example above:
_(NOTE: currently fleece expects the `crash-limit` flag to be defined on triage tests!)_
```golang
# see ./example/example_fuzz_test.go

//+build gofuzz

package example

var crashLimit int

func init() {
	flag.IntVar(&crashLimit, "crash-limit", 1000, "number of crashing inputs to test before stopping")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestFuzzBuggyFunc(t *testing.T) {
	_, panics, _ := fuzzing.
		MustNewCrasherIterator(FuzzBuggyFunc).
		TestFailingLimit(t, crashLimit)

	require.Zero(t, panics)
}
```

Use the `triage` command to look at the next crashing input and its stack trace, as well as to get a summary of where you are in the overall triage process.
With this information you should be able to debug the issue and re-run the triage test to see if that input is still crashing.
```bash
fleece triage ./example FuzzBuggyFunc
# see fleece triage --help
```

## Contributing
#### Updating bindata files

_(see: https://github.com/go-bindata/go-bindata)_
```bash
# Install gobindata
GO111MODULE=off go get github.com/go-bindata/go-bindata/...

# In fleece repo root
go-bindata -pkg bindata -o bindata/bindata.go -ignore=bindata\\.go -prefix=bindata ./bindata/...
```


#### Feature Wishlist
- Internals
  + [ ] Parallelize triaging!
  + [ ] Switch to docker engine api
- Generators
  + [ ] fuzz functions
  + [ ] triage tests
- Reporting
  + [ ] Summary (per fuzz function)
  + [ ] Unique crashing outputs (per fuzz function)
  + [ ] Issue tracking integration
