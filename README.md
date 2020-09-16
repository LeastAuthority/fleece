### How to use
##### Setup
```bash
# Install lafuzz CLI binary:
GO111MODULE=off go get -u github.com/leastauthority/lafuzz/cmd/lafuzz

# Navigate to the root of the repo you'll be fuzzing:
cd ./path/to/repo/root

# Init lafuzz files in repo (see `lafuzz -h` for help):
lafuzz init 
```

##### Fuzzing
```bash
lafuzz fuzz <package path> <fuzz function name>
```

##### Triaging
```bash
lafuzz triage <package path> <fuzz function name>
```

### Contributing
##### Updating docker files (and/or other assets)

_(see: https://github.com/go-bindata/go-bindata)_
```bash
# Install gobindata
GO111MODULE=off go get github.com/go-bindata/go-bindata/...

# In lafuzz repo root
go-bindata -o fuzzing/docker.go ./docker/...
```