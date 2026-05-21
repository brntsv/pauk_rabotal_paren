package app

import (
	"fmt"
	"os"
	"path/filepath"

	"pauk_rabotal_paren/internal/hotkey"
	"pauk_rabotal_paren/internal/sound"
)

const soundFile = "pauk_rabotal.mp3"

// Run starts the global keyboard listener and blocks until the app exits.
func Run() error {
	soundPath, err := resolveSoundPath(soundFile)
	if err != nil {
		return err
	}

	events := make(chan hotkey.Event, 32)
	go handleKeyboardEvents(events, soundPath)

	fmt.Println(hotkey.Instruction())

	if err := hotkey.Listen(events); err != nil {
		return fmt.Errorf("слушать клавиатуру: %w", err)
	}

	return nil
}

func handleKeyboardEvents(events <-chan hotkey.Event, soundPath string) {
	soundHotkeyDown := false

	for event := range events {
		switch event.Kind {
		case hotkey.Sound:
			if event.Down && !soundHotkeyDown {
				soundHotkeyDown = true
				sound.Play(soundPath)
			}
			if !event.Down {
				soundHotkeyDown = false
			}
		case hotkey.Exit:
			if event.Down {
				fmt.Println("Выход по Esc")
				os.Exit(0)
			}
		}
	}
}

func resolveSoundPath(name string) (string, error) {
	workingDir, _ := os.Getwd()
	executablePath, _ := os.Executable()

	return resolveSoundPathFrom(name, workingDir, executablePath)
}

func resolveSoundPathFrom(name string, workingDir string, executablePath string) (string, error) {
	candidates := make([]string, 0, 4)

	if workingDir != "" {
		candidates = append(candidates, filepath.Join(workingDir, name))
	}

	if executablePath != "" {
		executableDir := filepath.Dir(executablePath)
		candidates = append(candidates, filepath.Join(executableDir, name))
		candidates = append(candidates, filepath.Join(projectRootForBuildDir(executableDir), name))
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && !info.IsDir() {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("файл звука %q не найден рядом с рабочей директорией, бинарником или корнем проекта", name)
}

func projectRootForBuildDir(executableDir string) string {
	parentDir := filepath.Dir(executableDir)
	if filepath.Base(parentDir) != "target" {
		return executableDir
	}

	switch filepath.Base(executableDir) {
	case "debug", "release":
		return filepath.Dir(parentDir)
	default:
		return executableDir
	}
}
