package hotkey

type Kind int

const (
	Sound Kind = iota
	Exit
)

type Event struct {
	Kind Kind
	Down bool
}

func enqueue(events chan<- Event, event Event) {
	if events == nil {
		return
	}

	select {
	case events <- event:
	default:
	}
}
