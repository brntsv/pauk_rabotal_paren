//go:build windows

package sound

import (
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

var (
	windowsWinMM         = syscall.NewLazyDLL("winmm.dll")
	windowsMCISendString = windowsWinMM.NewProc("mciSendStringW")
)

func Play(soundPath string) {
	go playWindows(soundPath)
}

func playWindows(soundPath string) {
	alias := fmt.Sprintf("pauk_%d", time.Now().UnixNano())
	openCommand := fmt.Sprintf(`open "%s" type mpegvideo alias %s`, soundPath, alias)

	if err := mciSendString(openCommand); err != nil {
		log.Printf("не удалось открыть звук %s: %v", soundPath, err)
		return
	}

	done := make(chan struct{})
	defer func() {
		close(done)
		_ = mciSendString("close " + alias)
	}()

	go func() {
		select {
		case <-time.After(30 * time.Second):
			_ = mciSendString("stop " + alias)
			_ = mciSendString("close " + alias)
		case <-done:
		}
	}()

	if err := mciSendString("play " + alias + " wait"); err != nil {
		log.Printf("ошибка воспроизведения %s: %v", soundPath, err)
	}
}

func mciSendString(command string) error {
	commandPtr, err := syscall.UTF16PtrFromString(command)
	if err != nil {
		return err
	}

	ret, _, _ := windowsMCISendString.Call(
		uintptr(unsafe.Pointer(commandPtr)),
		0,
		0,
		0,
	)
	if ret != 0 {
		return fmt.Errorf("mciSendStringW вернул код %d для команды %q", ret, command)
	}

	return nil
}
