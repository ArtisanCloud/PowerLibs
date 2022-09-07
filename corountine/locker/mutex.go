package locker

import (
	"reflect"
	"sync"
)

const mutexLocked = 1

func MutexLocked(m *sync.Mutex) bool {
	state := reflect.ValueOf(m).Elem().FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func RWMutexWriteLocked(rw *sync.RWMutex) bool {
	// RWMutex has a "w" sync.Mutex field for write lock
	state := reflect.ValueOf(rw).Elem().FieldByName("w").FieldByName("state")
	return state.Int()&mutexLocked == mutexLocked
}

func RWMutexReadLocked(rw *sync.RWMutex) bool {
	return reflect.ValueOf(rw).Elem().FieldByName("readerCount").Int() > 0
}
