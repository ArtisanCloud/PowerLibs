package locker

import (
	"fmt"
	"sync"
	"testing"
)

func Test_MutexLocked(t *testing.T) {

	m := sync.Mutex{}
	fmt.Println("m locked =", MutexLocked(&m))
	m.Lock()
	fmt.Println("m locked =", MutexLocked(&m))
	m.Unlock()
	fmt.Println("m locked =", MutexLocked(&m))

}

func Test_RWMutexWriteAndReadLocked(t *testing.T) {

	rw := sync.RWMutex{}
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.Lock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.Unlock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.RLock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.RLock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.RUnlock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))
	rw.RUnlock()
	fmt.Println("rw write locked =", RWMutexWriteLocked(&rw), " read locked =", RWMutexReadLocked(&rw))

}
