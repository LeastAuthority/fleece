#! /usr/bin/env bash
set -e

# TODO: replace this shell script
# $1 - fuzzer package
# $2 - name of the fuzzer function (will be used for the `workdir` argument to `go-fuzz`).
# `--build` - (optional) builds fuzzer package before being run.
# Additional args following `--` are passed directly to `go-fuzz`.

if [[ $# -lt 2 ]]; then
  echo "usage: entrypoint.sh <fuzzer package path> <fuzzer func name> [-b|--build] [-- [go-fuzz arg[, ...]]]"
fi

pkg=$1
shift
name=$1
shift
workdir=./lafuzz/workdirs/${name}
bin=${pkg:2}-fuzz.zip

while [[ $# -gt 0 ]]; do
  case $1 in
  -b | --build)
    go-fuzz-build ${pkg}
    mv ${pkg}-fuzz.zip ${workdir}/${bin}
    shift
    ;;
  --)
    shift
    rest_args=$@
    break
    ;;
  *)
    shift
    ;;
  esac
done

rest_args=$@

go-fuzz -bin=${workdir}/${bin}-fuzz.zip -func=${name} -workdir=${workdir} $rest_args
