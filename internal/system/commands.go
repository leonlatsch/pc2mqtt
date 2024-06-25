package system

import (
	"errors"
	"os/exec"
	"runtime"
)

const (
	WINDOWS = "windows"
	MACOS   = "darwin"
	LINUX   = "linux"
)

// Shutdown

var shutdownSystemCommands = map[string]exec.Cmd{
	WINDOWS: *exec.Command("shutdown", "/s"),
	MACOS:   *exec.Command("shutdown", "-h", "now"),
	LINUX:   *exec.Command("poweroff"),
}

func GetShutdownCommand() (*exec.Cmd, error) {
	cmd, ok := shutdownSystemCommands[runtime.GOOS]

	if !ok {
		return nil, errors.New(runtime.GOOS + " does not support shutdown")
	}

	return &cmd, nil
}

// REBOOT

var rebootSystemCommands = map[string]exec.Cmd{
	WINDOWS: *exec.Command("shutdown", "/r"),
	MACOS:   *exec.Command("reboot"),
	LINUX:   *exec.Command("reboot"),
}

func GetRebootCommand() (*exec.Cmd, error) {
	cmd, ok := rebootSystemCommands[runtime.GOOS]

	if !ok {
		return nil, errors.New(runtime.GOOS + " does not support reboot")
	}

	return &cmd, nil
}
