// copied directly from https://github.com/moby/moby/blob/master/pkg/reexec/reexec.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

var registeredInitializers = make(map[string]func())

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registered under name %q", name))
	}

	registeredInitializers[name] = initializer
}

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.
func Init() bool {
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()

		return true
	}
	return false
}

func naiveSelf() string {
	name := os.Args[0]
	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err == nil {
			return lp
		}
	}
	// handle conversion of relative paths to absolute
	if absName, err := filepath.Abs(name); err == nil {
		return absName
	}
	// if we couldn't get absolute name, return original
	// (NOTE: Go only errors on Abs() if os.Getwd fails)
	return name
}

// Taken directly from https://github.com/moby/moby/blob/master/pkg/reexec/command_linux.go

// Self returns the path to the current process's binary.
// Returns "/proc/self/exe".
func Self() string {
	return "/proc/self/exe"
}

// Command returns *exec.Cmd which has Path as current binary. Also it setting
// SysProcAttr.Pdeathsig to SIGTERM.
// This will use the in-memory version (/proc/self/exe) of the current binary,
// it is thus safe to delete or replace the on-disk binary (os.Args[0]).
func Command(args ...string) *exec.Cmd {
	return &exec.Cmd{
		Path: Self(),
		Args: args,
		SysProcAttr: &syscall.SysProcAttr{
			Pdeathsig: unix.SIGTERM,
		},
	}
}
