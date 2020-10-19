package fuzzing

import (
	"path/filepath"
	"reflect"
	"runtime"
)

type Func func(data []byte) int

type Env struct {
	FleeceDir string
}

func NewEnv(fleeceDir string) *Env {
	return &Env{
		FleeceDir: fleeceDir,
	}
}

func NewLocalEnv() (*Env, error) {
	fleeceDir, err := config.GetFleeceDir()
	if err != nil {
		return nil, err
	}

	return &Env{
		FleeceDir: fleeceDir,
	}, nil
}

func (e *Env) GetWorkdirs() string {
	return filepath.Join(e.FleeceDir, "workdirs")
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

