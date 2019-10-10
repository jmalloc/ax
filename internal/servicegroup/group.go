package servicegroup

import (
	"context"
	"sync"

	"go.uber.org/multierr"
)

// A Group is a set of subtasks that occur as part of a parent task which runs
// until an error occurs or its context is canceled.
type Group struct {
	// ctx is derived from the context passed to NewGroup.
	// It is passed to the tasks started in this group.
	ctx context.Context

	// cancel cancels ctx. It is called when one of the goroutines in the group
	// returns a non-nil error.
	cancel func()

	// done is closed when an error has occurred (including the context being
	// canceled) and all existing tasks have returned.
	done chan struct{}

	// wg keeps track of the tasks that are currently being executed.
	// wg.Wait() is only called after the context is canceled. wgM protects new
	// against no goroutines being started after the context has been canceled.
	wgM sync.RWMutex
	wg  sync.WaitGroup

	// err is the set of non-nil errors returned by the tasks.
	errM sync.RWMutex
	err  error
}

// NewGroup returns a Group that is bound to ctx.
// No new tasks can be started once the group's context has been canceled.
func NewGroup(ctx context.Context) *Group {
	ctx, cancel := context.WithCancel(ctx)

	g := &Group{
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}

	go g.wait()

	return g
}

// Wait blocks until an error occurs or the group's context is canceled.
func (g *Group) Wait() error {
	<-g.done

	g.errM.RLock()
	defer g.errM.RUnlock()

	return g.err
}

// Go calls the given function in a new goroutine.
//
// It returns a non-nil error if the group's context has been canceled, or
// another task has returned a non-nil error.
//
// The context passed to fn is canceled when the group's context is canceled,
// or some other task returns a non-nil error.
func (g *Group) Go(
	fn func(context.Context) error,
) error {
	select {
	case <-g.ctx.Done():
		return g.ctx.Err()
	default:
	}

	g.wgM.RLock()
	defer g.wgM.RUnlock()

	select {
	case <-g.ctx.Done():
		return g.ctx.Err()
	default:
		g.wg.Add(1)
		go g.execute(fn)
		return nil
	}
}

// execute calls fn, and captures the returned error.
// if fn returns a non-nil error, g.ctx is canceled.
func (g *Group) execute(
	fn func(context.Context) error,
) {
	defer g.wg.Done()

	err := fn(g.ctx)
	if err == nil {
		return
	}

	g.errM.Lock()
	defer g.errM.Unlock()

	if g.err == nil || err != context.Canceled {
		g.err = multierr.Append(g.err, err)
	}

	g.cancel()
}

// wait closes g.done after an error occurs and all tasks return.
func (g *Group) wait() {
	<-g.ctx.Done()

	g.wgM.Lock()
	defer g.wgM.Unlock()

	g.wg.Wait()
	close(g.done)
}
