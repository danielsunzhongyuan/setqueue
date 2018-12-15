package setqueue

import (
    "reflect"
    "errors"
    "log"
)

type SetMessageQueue interface {
    GenericMessageQueue
}

type MySetQueue struct {
    data chan interface{}
    m    ConcurrentMap
}

func (sq *MySetQueue) Add(message interface{}) error {
    if sq.m.Contains(message) {
        return KEY_EXISTS
    }
    data, ok := sq.m.Put(message, true)
    if !ok {
        log.Println("Failed inserting a message to queue:", data)
        return KEY_EXISTS
    }

    sq.data <- message
    return nil
}

func (sq *MySetQueue) Get(handler Handler) error {
    message := <- sq.data
    sq.m.Remove(message)
    if err := handler(message); err != nil {
        log.Println("process message", message, "failed, then put it back")
        sq.Add(message)
        return PROCESS_FAILED
    }
    return nil
}

func (sq *MySetQueue) Close() {
    close(sq.data)
    sq.m.Clear()
}

func NewSetQueue(keyType, valueType reflect.Type, cacheNumber int) SetMessageQueue {
    return &MySetQueue{
        data: make(chan interface{}, cacheNumber),
        m:    NewConcurrentMap(keyType, valueType),
    }
}
