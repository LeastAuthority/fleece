package internal

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"golang.org/x/tools/go/packages"
)

type fuzzFunc func([]byte) int

func GetModPath(importPath string) (string, error) {
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

func GetFuncName(f interface{}) string {
	val := reflect.ValueOf(f)
	addr := val.Pointer()
	return filepath.Ext(runtime.FuncForPC(addr).Name())[1:]
}
