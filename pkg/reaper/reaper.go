package reaper

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Reaper struct {
	lock sync.Mutex
}

func New() *Reaper {
	return &Reaper{}
}

func (r *Reaper) Lock() {
	r.lock.Lock()
}

func (r *Reaper) Unlock() {
	r.lock.Unlock()
}

func (r *Reaper) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGCHLD)
	defer signal.Reset(syscall.SIGCHLD)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.reapall()
		case <-sigchan:
			r.reapall()
		}
	}
}

func (r *Reaper) reapall() {
	r.Lock()
	defer r.Unlock()

	for {
		if pid, _ := syscall.Wait4(-1, nil, syscall.WNOHANG, nil); pid <= 0 {
			return
		}
	}
}
