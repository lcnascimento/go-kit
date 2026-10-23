package outbox

import "sync"

type state int

const (
	idle state = iota
	running
	stopped
)

type lifecycle struct {
	mu    sync.Mutex
	state state
}

func (l *lifecycle) start() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.state != idle {
		return ErrAlreadyStarted
	}

	l.state = running

	return nil
}

func (l *lifecycle) stop() (previous state) {
	l.mu.Lock()
	defer l.mu.Unlock()

	previous = l.state
	l.state = stopped

	return previous
}

func (l *lifecycle) is(s state) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.state == s
}
