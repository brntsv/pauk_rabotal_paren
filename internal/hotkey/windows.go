//go:build windows

package hotkey

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

const (
	windowsWHKeyboardLL = 13

	windowsWMKeyDown    = 0x0100
	windowsWMKeyUp      = 0x0101
	windowsWMSysKeyDown = 0x0104
	windowsWMSysKeyUp   = 0x0105

	windowsLLKHFExtended = 0x01

	windowsVKEscape = 0x1B
	windowsVKMenu   = 0x12
	windowsVKRAlt   = 0xA5
)

type windowsKBDLLHookStruct struct {
	vkCode      uint32
	scanCode    uint32
	flags       uint32
	time        uint32
	dwExtraInfo uintptr
}

type windowsPoint struct {
	x int32
	y int32
}

type windowsMsg struct {
	hwnd    uintptr
	message uint32
	wParam  uintptr
	lParam  uintptr
	time    uint32
	pt      windowsPoint
}

var (
	windowsUser32 = syscall.NewLazyDLL("user32.dll")

	windowsSetWindowsHookEx   = windowsUser32.NewProc("SetWindowsHookExW")
	windowsCallNextHookEx     = windowsUser32.NewProc("CallNextHookEx")
	windowsUnhookWindowsHook  = windowsUser32.NewProc("UnhookWindowsHookEx")
	windowsGetMessage         = windowsUser32.NewProc("GetMessageW")
	windowsEvents             chan<- Event
	windowsKeyboardHookHandle uintptr
	windowsKeyboardProc       = syscall.NewCallback(handleWindowsKeyboardEvent)
)

func Instruction() string {
	return "Слушаю клавиатуру. Нажми правый Alt для звука, Esc для выхода."
}

func Listen(events chan<- Event) error {
	windowsEvents = events

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	hook, _, err := windowsSetWindowsHookEx.Call(
		windowsWHKeyboardLL,
		windowsKeyboardProc,
		0,
		0,
	)
	if hook == 0 {
		return wrapWindowsError("не удалось запустить глобальный хук клавиатуры", err)
	}
	windowsKeyboardHookHandle = hook
	defer windowsUnhookWindowsHook.Call(hook)

	var msg windowsMsg
	for {
		ret, _, err := windowsGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		switch int32(ret) {
		case -1:
			return wrapWindowsError("ошибка в цикле сообщений Windows", err)
		case 0:
			return nil
		}
	}
}

func handleWindowsKeyboardEvent(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		keyEvent := (*windowsKBDLLHookStruct)(unsafe.Pointer(lParam))
		down := wParam == windowsWMKeyDown || wParam == windowsWMSysKeyDown
		up := wParam == windowsWMKeyUp || wParam == windowsWMSysKeyUp

		if down || up {
			if isWindowsRightAlt(keyEvent) {
				enqueue(windowsEvents, Event{Kind: Sound, Down: down})
			}

			switch keyEvent.vkCode {
			case windowsVKEscape:
				enqueue(windowsEvents, Event{Kind: Exit, Down: down})
			}
		}
	}

	ret, _, _ := windowsCallNextHookEx.Call(
		windowsKeyboardHookHandle,
		uintptr(nCode),
		wParam,
		lParam,
	)
	return ret
}

func isWindowsRightAlt(event *windowsKBDLLHookStruct) bool {
	return event.vkCode == windowsVKRAlt ||
		event.vkCode == windowsVKMenu && event.flags&windowsLLKHFExtended != 0
}

func wrapWindowsError(message string, err error) error {
	if errno, ok := err.(syscall.Errno); ok && errno == 0 {
		return fmt.Errorf("%s", message)
	}
	return fmt.Errorf("%s: %w", message, err)
}
