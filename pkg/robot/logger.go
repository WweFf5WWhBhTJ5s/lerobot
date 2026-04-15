package robot

import (
	"fmt"
	"io"
	"os"
	"time"
)

// Logger writes robot events to an io.Writer in a structured text format.
type Logger struct {
	out io.Writer
}

// NewLogger creates a Logger that writes to out.
// If out is nil, os.Stdout is used.
func NewLogger(out io.Writer) *Logger {
	if out == nil {
		out = os.Stdout
	}
	return &Logger{out: out}
}

// Handle implements the notifier handler signature and logs the event.
func (l *Logger) Handle(e Event) {
	timestamp := time.Now().UTC().Format(time.RFC3339)
	line := fmt.Sprintf("%s [%s] robot=%s", timestamp, e.Type, e.RobotID)
	if e.Status != nil {
		line += fmt.Sprintf(" status=%s", e.Status)
	}
	fmt.Fprintln(l.out, line)
}

// Subscribe registers the logger as a handler on the given Notifier
// and returns an unsubscribe function.
func (l *Logger) Subscribe(n *Notifier) func() {
	return n.Subscribe(l.Handle)
}
