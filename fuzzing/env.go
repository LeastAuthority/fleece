package fuzzing

import (
	"path/filepath"
	"reflect"
	"runtime"
)

type Func func(data []byte) int

type Env struct {
	outputRoot string
}

func NewEnv(outputRoot string) *Env {
	return &Env{
		outputRoot: outputRoot,
	}
}

func (e *Env) GetWorkdirs() string {
	return filepath.Join(e.outputRoot, "workdirs")
}

func (e *Env) GetCrasherDir(fuzzFunc Func) string {
	name := GetFuncName(fuzzFunc)
	return filepath.Join(e.GetWorkdirs(), name, "crashers")
}

func GetFuncName(f Func) string {
	val := reflect.ValueOf(f)
	addr := val.Pointer()
	return filepath.Ext(runtime.FuncForPC(addr).Name())[1:]
}

