package mrevent

import (
    "context"
    "fmt"
    "sync"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
    // eventPool - thread safe structure stores events, their handlers and functions for management
    eventPool struct {
        mutex       sync.RWMutex
        waitGroup   sync.WaitGroup
        subscribers map[string][]chan eventArgsWithContext
        bufferSize  uint32
    }

    eventArgsWithContext struct {
        ctx context.Context
        args mrcore.EventArgs
    }
)

// Make sure the eventPool conforms with the mrcore.EventPool interface
var _ mrcore.EventPool = (*eventPool)(nil)

func NewEventPool(bufferSize uint32) *eventPool {
    return &eventPool{
        subscribers: map[string][]chan eventArgsWithContext{},
        bufferSize: bufferSize,
    }
}

// Listen - Subscribe on eventName
func (e *eventPool) Listen(eventName string, handlerFunc mrcore.EventHandlerFunc) {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    subscriber := make(chan eventArgsWithContext, e.bufferSize)

    go func(ch chan eventArgsWithContext, waitGroup *sync.WaitGroup) {
        for {
            data, more := <-ch

            if !more {
                break
            }

            handlerFunc(data.ctx, data.args)
            waitGroup.Done()
        }
    }(subscriber, &e.waitGroup)

    e.subscribers[eventName] = append(e.subscribers[eventName], subscriber)
}

// Emit - Call eventName there
func (e *eventPool) Emit(eventName string, args mrcore.EventArgs) error {
    return e.EmitWithContext(context.Background(), eventName, args)
}

// EmitWithContext - Call eventName there
func (e *eventPool) EmitWithContext(ctx context.Context, eventName string, args mrcore.EventArgs) error {
    e.mutex.RLock()
    defer e.mutex.RUnlock()

    if !e.has(eventName) {
        return fmt.Errorf("subscribers to %s not found", eventName)
    }

    for _, ch := range e.subscribers[eventName] {
        e.waitGroup.Add(1)
        ch <- eventArgsWithContext{ctx, args}
    }

    return nil
}

// List - Returns a list of events that listeners are subscribed to
func (e *eventPool) List() []string {
    e.mutex.RLock()
    defer e.mutex.RUnlock()

    return e.list()
}

// Has - Checks if there are listeners for the passed eventName
func (e *eventPool) Has(eventName string) bool {
    e.mutex.RLock()
    defer e.mutex.RUnlock()

    return e.has(eventName)
}

// Remove - Removes listeners by event name
// Removing listeners closes subscribers and stops the goroutine.
//
// If you call the function without the "names" parameter, all listeners of all events will be removed.
func (e *eventPool) Remove(eventNames ...string) {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    if len(eventNames) == 0 {
        eventNames = e.list()
    }

    for _, eventName := range eventNames {
        for _, subscriber := range e.subscribers[eventName] {
            close(subscriber)
        }

        delete(e.subscribers, eventName)
    }
}

// Wait - Blocks the thread until all running events are completed
func (e *eventPool) Wait() {
    e.mutex.Lock()
    defer e.mutex.Unlock()

    e.waitGroup.Wait()
}

func (e *eventPool) list() []string {
    list := make([]string, 0, len(e.subscribers))

    for eventName := range e.subscribers {
        list = append(list, eventName)
    }

    return list
}

func (e *eventPool) has(eventName string) bool {
    _, ok := e.subscribers[eventName]

    return ok
}
