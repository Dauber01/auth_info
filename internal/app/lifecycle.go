package app

import (
	"errors"
	"fmt"
	"sync"
)

type closeHook struct {
	name  string
	close func() error
}

// Lifecycle owns initialized resources. Providers register in dependency order before publication.
type Lifecycle struct {
	hooks []closeHook
	once  sync.Once
	err   error
}

// Add registers a resource after successful acquisition, before the application is published.
func (l *Lifecycle) Add(name string, close func() error) {
	l.hooks = append(l.hooks, closeHook{name, close})
}

// Close releases every resource in reverse order and returns the same result on repeated calls.
func (l *Lifecycle) Close() error {
	l.once.Do(func() {
		for i := len(l.hooks) - 1; i >= 0; i-- {
			hook := l.hooks[i]
			if err := hook.close(); err != nil {
				l.err = errors.Join(l.err, fmt.Errorf("close %s: %w", hook.name, err))
			}
		}
	})
	return l.err
}
