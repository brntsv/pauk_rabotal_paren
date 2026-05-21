//go:build darwin

package hotkey

/*
#cgo darwin LDFLAGS: -framework ApplicationServices -framework CoreFoundation
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>

extern void goHandleKeyboardEvent(unsigned long long eventType, long long keyCode, int down);

static CFMachPortRef eventTap;

static CGEventRef keyboardEventCallback(CGEventTapProxy proxy, CGEventType type, CGEventRef event, void *refcon) {
	if (type == kCGEventTapDisabledByTimeout || type == kCGEventTapDisabledByUserInput) {
		if (eventTap != NULL) {
			CGEventTapEnable(eventTap, true);
		}
		return event;
	}

	if (type != kCGEventKeyDown && type != kCGEventFlagsChanged) {
		return event;
	}

	long long keyCode = CGEventGetIntegerValueField(event, kCGKeyboardEventKeycode);
	int down = 1;

	if (type == kCGEventFlagsChanged) {
		CGEventFlags flags = CGEventGetFlags(event);
		down = (flags & kCGEventFlagMaskCommand) != 0;
	}

	goHandleKeyboardEvent((unsigned long long)type, keyCode, down);
	return event;
}

static int runDarwinKeyboardListener(void) {
	CGEventMask mask = CGEventMaskBit(kCGEventKeyDown) | CGEventMaskBit(kCGEventFlagsChanged);
	eventTap = CGEventTapCreate(
		kCGSessionEventTap,
		kCGHeadInsertEventTap,
		kCGEventTapOptionListenOnly,
		mask,
		keyboardEventCallback,
		NULL
	);

	if (eventTap == NULL) {
		return 0;
	}

	CFRunLoopSourceRef source = CFMachPortCreateRunLoopSource(kCFAllocatorDefault, eventTap, 0);
	CFRunLoopAddSource(CFRunLoopGetCurrent(), source, kCFRunLoopCommonModes);
	CGEventTapEnable(eventTap, true);
	CFRunLoopRun();

	CFRelease(source);
	CFRelease(eventTap);
	return 1;
}
*/
import "C"

import (
	"fmt"
	"runtime"
)

const (
	darwinKeyRightCommand = 54
	darwinKeyEscape       = 53
)

var darwinEvents chan<- Event

func Instruction() string {
	return "Слушаю клавиатуру. Нажми правый Command для звука, Esc для выхода."
}

func Listen(events chan<- Event) error {
	darwinEvents = events

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if C.runDarwinKeyboardListener() == 0 {
		return fmt.Errorf("не удалось запустить глобальный хук клавиатуры; проверь разрешения Accessibility для терминала")
	}

	return nil
}

//export goHandleKeyboardEvent
func goHandleKeyboardEvent(_ C.ulonglong, keyCode C.longlong, down C.int) {
	eventDown := down != 0

	switch int64(keyCode) {
	case darwinKeyRightCommand:
		enqueue(darwinEvents, Event{Kind: Sound, Down: eventDown})
	case darwinKeyEscape:
		enqueue(darwinEvents, Event{Kind: Exit, Down: eventDown})
	}
}
