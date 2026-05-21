//go:build linux

package hotkey

import "fmt"

func Instruction() string {
	return "Linux backend пока не реализован: глобальные хоткеи зависят от X11, Wayland или evdev."
}

func Listen(_ chan<- Event) error {
	return fmt.Errorf("Linux пока не поддержан: глобальные хоткеи зависят от X11, Wayland или evdev и требуют отдельного backend")
}
