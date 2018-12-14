package setqueue

import (
    "reflect"
    "sync"
    "bytes"
    "fmt"
)

type ConcurrentMap interface {
    GenericMap
}

type MyConcurrentMap struct {
    m         map[interface{}]interface{}
    keyType   reflect.Type
    valueType reflect.Type
    mutex     sync.Mutex
}

func (cmap *MyConcurrentMap) Get(key interface{}) interface{} {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    return cmap.m[key]
}

func (cmap *MyConcurrentMap) isAcceptablePair(k, v interface{}) bool {
    if k == nil || reflect.TypeOf(k) != cmap.keyType {
        return false
    }
    if v == nil || reflect.TypeOf(v) != cmap.valueType {
        return false
    }
    return true
}

func (cmap *MyConcurrentMap) Put(key interface{}, value interface{}) (interface{}, bool) {
    if !cmap.isAcceptablePair(key, value) {
        return nil, false
    }
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    oldValue := cmap.m[key]
    cmap.m[key] = value
    return oldValue, true
}

func (cmap *MyConcurrentMap) Remove(key interface{}) interface{} {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    oldValue := cmap.m[key]
    delete(cmap.m, key)
    return oldValue
}

func (cmap *MyConcurrentMap) Clear() {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    cmap.m = make(map[interface{}]interface{})
}

func (cmap *MyConcurrentMap) Len() int {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    return len(cmap.m)
}

func (cmap *MyConcurrentMap) Contains(key interface{}) bool {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    _, ok := cmap.m[key]
    return ok
}

func (cmap *MyConcurrentMap) Keys() []interface{} {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    initialLen := len(cmap.m)
    keys := make([]interface{}, initialLen)
    index := 0
    for k := range cmap.m {
        keys[index] = k
        index++
    }
    return keys
}

func (cmap *MyConcurrentMap) Values() []interface{} {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    initialLen := len(cmap.m)
    values := make([]interface{}, initialLen)
    index := 0
    for _, v := range cmap.m {
        values[index] = v
        index++
    }
    return values
}

func (cmap *MyConcurrentMap) ToMap() map[interface{}]interface{} {
    cmap.mutex.Lock()
    defer cmap.mutex.Unlock()
    replica := make(map[interface{}]interface{})
    for k, v := range cmap.m {
        replica[k] = v
    }
    return replica
}

func (cmap *MyConcurrentMap) KeyType() reflect.Type {
    return cmap.keyType
}

func (cmap *MyConcurrentMap) ValueType() reflect.Type {
    return cmap.valueType
}

func (cmap *MyConcurrentMap) String() string {
    var buf bytes.Buffer
    buf.WriteString("ConcurrentMap<")
    buf.WriteString(cmap.keyType.Kind().String())
    buf.WriteString(",")
    buf.WriteString(cmap.valueType.Kind().String())
    buf.WriteString(">{")
    first := true
    for k, v := range cmap.m {
        if first {
            first = false
        } else {
            buf.WriteString(" ")
        }
        buf.WriteString(fmt.Sprintf("%v", k))
        buf.WriteString(":")
        buf.WriteString(fmt.Sprintf("%v", v))
    }
    buf.WriteString("}")
    return buf.String()
}

func NewConcurrentMap(keyType, valueType reflect.Type) ConcurrentMap {
    return &MyConcurrentMap{
        keyType:   keyType,
        valueType: valueType,
        m:         make(map[interface{}]interface{}),
    }
}
