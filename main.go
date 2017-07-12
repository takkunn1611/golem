package main

import (
	"flag"
	"fmt"
	"os"

	luar "layeh.com/gopher-luar"

	lua "github.com/yuin/gopher-lua"
)

func newLState() *lua.LState {
	L := lua.NewState()
	L.PreloadModule(
		"golem",
		func(L *lua.LState) int {
			mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})
			L.SetField(mod, "worker", luar.New(L, func(name string, parent chan lua.LValue) chan lua.LValue {
				worker := make(chan lua.LValue)

				go func() {
					workerL := newLState()

					workerL.SetGlobal("parent", lua.LChannel(parent))

					if err := workerL.DoFile(name); err != nil {
						close(worker)
					}

					workerFn := workerL.GetGlobal("worker")

					for msg := range worker {
						if err := workerL.CallByParam(lua.P{
							Fn:      workerFn,
							NRet:    0,
							Protect: true,
						}, msg); err != nil {
							close(worker)
						}
					}
				}()

				return worker
			}))
			L.Push(mod)
			return 1
		},
	)
	return L
}

const (
	exitOk  = 0
	exitErr = 1
)

var (
	// Name of this program
	Name = "golem"
	// Version of this program
	Version = "0.0.1"
)

func rootCmd(args []string) int {
	var (
		version    bool
		currentDir string
	)
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.BoolVar(&version, "version", false, "show version")
	flags.StringVar(&currentDir, "C", ".", "current directory")
	flags.SetOutput(os.Stderr)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", Name)
		flags.PrintDefaults()
	}

	if err := flags.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return exitOk
		}

		return exitErr
	}

	if version {
		fmt.Fprintf(os.Stderr, "%s version %s\n", Name, Version)
		return exitOk
	}

	handleErr := func(f string, args ...interface{}) int {
		fmt.Fprintf(os.Stderr, Name+": "+f, args...)
		return exitErr
	}

	os.Chdir(currentDir)

	L := newLState()
	defer L.Close()

	if err := L.DoFile("index.lua"); err != nil {
		return handleErr("%v\n", err)
	}

	if err := L.CallByParam(lua.P{
		Fn:      L.GetGlobal("main"),
		NRet:    0,
		Protect: true,
	}); err != nil {
		return handleErr("%v\n", err)
	}

	return exitOk
}

func main() {
	os.Exit(rootCmd(os.Args))
}
