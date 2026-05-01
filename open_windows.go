//go:build windows

package main

import (
	"os/exec"
	"syscall"
)

func openFile(path string) error {
	cmd := exec.Command("cmd", "/C", "start", "", path)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
