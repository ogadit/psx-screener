//go:build !windows

package main

import "fmt"

func openFile(_ string) error {
	return fmt.Errorf("not supported on this platform")
}
