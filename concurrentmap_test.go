package setqueue

import (
    "fmt"
    "reflect"
    "runtime/debug"
    "testing"
    "math/rand"
)

func testConcurrentMap(
    t *testing.T,
    newConcurrentMap func() ConcurrentMap,
    genKey func() interface{},
    genVal func() interface{},
    keyKind reflect.Kind,
    valueKind reflect.Kind,
) {
    mapType := fmt.Sprintf("ConcurrentMap<keyType=%s, valueType=%s>", keyKind, valueKind)
    defer func() {
        if err := recover(); err != nil {
            debug.PrintStack()
            t.Errorf("Fatal Error: %s: %s\n", mapType, err)
        }
    }()
    t.Logf("Starting Test %s ...", mapType)

    // Basic
    cmap := newConcurrentMap()
    expectedLen := 0
    if cmap.Len() != expectedLen {
        t.Errorf("Error: The length of %s value %d is not %d!\n",
            mapType, cmap.Len(), expectedLen)
        t.FailNow()
    }
    expectedLen = 5
    testMap := make(map[interface{}]interface{}, expectedLen)
    var invalidKey interface{}
    for i := 0; i < expectedLen; i++ {
        key := genKey()
        testMap[key] = genVal()
        if invalidKey == nil {
            invalidKey = key
        }
    }
    for key, value := range testMap {
        oldValue, ok := cmap.Put(key, value)
        if !ok {
            t.Errorf("Error: Put (%v, %v) to %s value %d is failing!\n",
                key, value, mapType, cmap)
            t.FailNow()
        }
        if oldValue != nil {
            t.Errorf("Error: Already had a (%v, %v) in %s value %d!\n",
                key, value, mapType, cmap)
            t.FailNow()
        }
        t.Logf("Put (%v, %v) to the %s value %v.",
            key, value, mapType, cmap)
    }
    if cmap.Len() != expectedLen {
        t.Errorf("Error: The length of %s value %d is not %d!\n",
            mapType, cmap.Len(), expectedLen)
        t.FailNow()
    }
    for key, value := range testMap {
        contains := cmap.Contains(key)
        if !contains {
            t.Errorf("Error: The %s value %v do not contains %v!",
                mapType, cmap, key)
            t.FailNow()
        }
        actualValue := cmap.Get(key)
        if actualValue == nil {
            t.Errorf("Error: The %s value %v do not contains %v!",
                mapType, cmap, key)
            t.FailNow()
        }
        t.Logf("The %s value %v contains key %v.", mapType, cmap, key)
        if actualValue != value {
            t.Errorf("Error: The value of %s value %v with key %v do not equals %v!\n",
                mapType, actualValue, key, value)
            t.FailNow()
        }
        t.Logf("The value of %s value %v to key %v is %v.",
            mapType, cmap, key, actualValue)
    }
    oldValue := cmap.Remove(invalidKey)
    if oldValue == nil {
        t.Errorf("Error: Remove %v from %s value %d is failing!\n",
            invalidKey, mapType, cmap)
        t.FailNow()
    }
    t.Logf("Removed (%v, %v) from the %s value %v.",
        invalidKey, oldValue, mapType, cmap)
    delete(testMap, invalidKey)

    // Type
    actualValueType := cmap.ValueType()
    if actualValueType == nil {
        t.Errorf("Error: The value type of %s value is nil!\n",
            mapType)
        t.FailNow()
    }
    actualValueKind := actualValueType.Kind()
    if actualValueKind != valueKind {
        t.Errorf("Error: The value type of %s value %s is not %s!\n",
            mapType, actualValueKind, valueKind)
        t.FailNow()
    }
    t.Logf("The value type of %s value %v is %s.",
        mapType, cmap, actualValueKind)
    actualKeyKind := cmap.KeyType().Kind()
    if actualKeyKind != valueKind {
        t.Errorf("Error: The key type of %s value %s is not %s!\n",
            mapType, actualKeyKind, keyKind)
        t.FailNow()
    }
    t.Logf("The key type of %s value %v is %s.",
        mapType, actualKeyKind, keyKind)

    // Export
    keys := cmap.Keys()
    values := cmap.Values()
    pairs := cmap.ToMap()
    for key, value := range testMap {
        var hasKey bool
        for _, k := range keys {
            if k == key {
                hasKey = true
            }
        }
        if !hasKey {
            t.Errorf("Error: The keys of %s value %v do not contains %v!\n",
                mapType, cmap, key)
            t.FailNow()
        }
        var hasValue bool
        for _, v := range values {
            if v == value {
                hasValue = true
            }
        }
        if !hasValue {
            t.Errorf("Error: The values of %s value %v do not contains %v!\n",
                mapType, cmap, value)
            t.FailNow()
        }
        var hasPair bool
        for k, v := range pairs {
            if k == key && v == value {
                hasPair = true
            }
        }
        if !hasPair {
            t.Errorf("Error: The values of %s value %v do not contains (%v, %v)!\n",
                mapType, cmap, key, value)
            t.FailNow()
        }
    }

    // Clear
    cmap.Clear()
    if cmap.Len() != 0 {
        t.Errorf("Error: Clear %s value %d is failing!\n",
            mapType, cmap)
        t.FailNow()
    }
    t.Logf("The %s value % has been cleared.", mapType, cmap)
}

func TestInt64Cmap(t *testing.T) {
    newCmap := func() ConcurrentMap {
        keyType := reflect.TypeOf(int64(2))
        valueType := keyType
        return NewConcurrentMap(keyType, valueType)
    }

    testConcurrentMap(
        t,
        newCmap,
        func() interface{} {
            return rand.Int63n(1000)
        },
        func() interface{} {
            return rand.Int63n(1000)
        },
        reflect.Int64,
        reflect.Int64,
    )
}

func TestFloat64Cmap(t *testing.T) {
    newCmap := func() ConcurrentMap {
        keyType := reflect.TypeOf(float64(2))
        valueType := keyType
        return NewConcurrentMap(keyType, valueType)
    }
    testConcurrentMap(
        t,
        newCmap,
        func() interface{} { return rand.Float64() },
        func() interface{} { return rand.Float64() },
        reflect.Float64,
        reflect.Float64,
    )
}

//func TestStringCmap(t *testing.T) {
//    newCmap := func() ConcurrentMap {
//        keyType := reflect.TypeOf(string(2))
//        valueType := keyType
//        return NewConcurrentMap(keyType, valueType)
//    }
//    testConcurrentMap(
//        t,
//        newCmap,
//        func() interface{} { return genRandString() },
//        func() interface{} { return genRandString() },
//        reflect.String,
//        reflect.String)
//}

func BenchmarkConcurrentMap(b *testing.B) {
    keyType := reflect.TypeOf(int32(2))
    valueType := keyType
    cmap := NewConcurrentMap(keyType, valueType)
    var key, value int32
    fmt.Printf("N=%d.\n", b.N)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        seed := int32(i)
        key = seed
        value = seed << 10
        b.StartTimer()
        cmap.Put(key, value)
        _ = cmap.Get(key)
        b.StopTimer()
        b.SetBytes(8)
        b.StartTimer()
    }
    ml := cmap.Len()
    b.StopTimer()
    mapType := fmt.Sprintf("ConcurrentMap<%s, %s>",
        keyType.Kind().String(), valueType.Kind().String())
    b.Logf("The length of %s value is %d.\n", mapType, ml)
    b.StartTimer()
}

func BenchmarkMap(b *testing.B) {
    keyType := reflect.TypeOf(int32(2))
    valueType := keyType
    imap := make(map[interface{}]interface{})
    var key, value int32
    fmt.Printf("N=%d.\n", b.N)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        b.StopTimer()
        seed := int32(i)
        key = seed
        value = seed << 10
        b.StartTimer()
        imap[key] = value
        b.StopTimer()
        _ = imap[key]
        b.StopTimer()
        b.SetBytes(8)
        b.StartTimer()
    }
    ml := len(imap)
    b.StopTimer()
    mapType := fmt.Sprintf("Map<%s, %s>",
        keyType.Kind().String(), valueType.Kind().String())
    b.Logf("The length of %s value is %d.\n", mapType, ml)
    b.StartTimer()
}
