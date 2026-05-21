//go:build !darwin && !windows && !linux

package hotkey

import (
	"fmt"
	"runtime"
)

func Instruction() string {
	return fmt.Sprintf("Платформа %s/%s не поддержана.", runtime.GOOS, runtime.GOARCH)
}

func Listen(_ chan<- Event) error {
	return fmt.Errorf("платформа %s/%s не поддержана", runtime.GOOS, runtime.GOARCH)
}
