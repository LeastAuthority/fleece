# lafuzz
## Value proposition
Simple to use and extensible tool that facilitates a fuzz -> triage -> debug workflow with mininal impact on the target repo.

## How to use
#### Setup and help
```bash
# Install lafuzz CLI binary:
GO111MODULE=off go get -u github.com/leastauthority/lafuzz/cmd/lafuzz

# Navigate to the root of the repo you'll be fuzzing:
cd ./path/to/repo/root

# Init lafuzz files in repo (see `lafuzz -h` for help):
lafuzz init

# Init local fuzzing environment
lafuzz env init

OR do both at the same time
lafuzz init --env
```
```
A tool to manage fuzzing environment

Usage:
  lafuzz [command]

Available Commands:
  env         manage local fuzzing environment
  fuzz        start a fuzzing container for the specified fuzz function (blocking)
  help        Help about any command
  init        initialize lafuzz into repo
  triage      tests known crashing inputs and prints a summary.

Flags:
      --config string   config file (default is $(pwd)/.lafuzz.yaml)
  -h, --help            help for lafuzz

Use "lafuzz [command] --help" for more information about a command.
```

#### Fuzzing
Now fuzz functions need to be defined.
The current convention is that they should be in or around the same package as the code they exercise.
Conventionally the `gofuzz` build tag is used to exclude this code from normal builds and tests.

Fuzz functions follow the `go-fuzz` signature:
```golang
# /path/to/mypkg/my_fuzz.go

//+build gofuzz

func FuzzMyFunc(data []byte) int
```
_(see: [go-fuzz readme](https://github.com/dvyukov/go-fuzz/blob/master/README.md) for more details)_

To run your fuzz functions:
```bash
lafuzz fuzz ./path/to/mypkg FuzzMyFunc [--procs <n>]
# see lafuzz fuzz --help
```

#### Triaging
Once you've discovered some crashing inputs you can look through their stack traces to debug them.
It's possible that over the course of fuzzing you may discover many thousands of crashing inputs.
The convention here is to tackle these one at a time.

To debug, the convention here is to add a "triage test" which will run the crashing inputs back through the fuzz function that produced them to see if they're still crashing.
This also has the benefit of acting as a regression test when combined with a with a simple assertion at all inputs are no longer crashing, and retaining the crashers in version control.

Here's the triage test corresponding with our example above:
```golang
# /path/to/mypkg/my_fuzz_test.go

//+build gofuzz

func TestFuzzMyFunc(t *testing.T) {
	_, panics, _ := lafuzz.
		MustNewCrasherIterator(FuzzMyFunc).
		TestFailingLimit(t, 1000)

	require.Zero(t, panics)
}
```

Use the `triage` command to look at the next crashing input and its stack trace, as well as to get a summary of where you are in the overall triage process.
With this information you should be able to debug the issue and re-run the triage test to see if that input is still crashing.
```bash
lafuzz triage ./path/to/mypkg FuzzMyFunc
# see lafuzz triage --help
```

## Contributing
#### Updating docker files (and/or other assets)

_(see: https://github.com/go-bindata/go-bindata)_
```bash
# Install gobindata
GO111MODULE=off go get github.com/go-bindata/go-bindata/...

# In lafuzz repo root
go-bindata -pkg fuzzing -ignore docker\\.go -o docker/docker.go ./docker/...
```