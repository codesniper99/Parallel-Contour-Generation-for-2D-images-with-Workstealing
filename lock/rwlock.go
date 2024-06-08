// Package lock provides an implementation of a read-write lock
// that uses condition variables and mutexes.
package lock

import (
	"sync"
	"sync/atomic"
)

const maxConcurReaders = 32

type RwLock struct {
	mu          sync.Mutex
	readerCount atomic.Int32
	writerCount int
	writeCond   *sync.Cond
}

func NewRwLock() *RwLock {
	return &RwLock{
		writeCond: sync.NewCond(&sync.Mutex{}),
	}
}

func (rwlock *RwLock) Lock() {
	rwlock.mu.Lock()

	for rwlock.writerCount > 1 || rwlock.readerCount.Load() != 0 {
		//	fmt.Println("Actually waiting for Lock")
		rwlock.writeCond.Wait()
		//	fmt.Println("Passed out")
	}
	//fmt.Println("Got the lock, Wait for readers to end")

	//fmt.Println(rwlock.writerCount, " ", rwlock.readerCount)

	//fmt.Println("Now I can increment my writer safely")
	rwlock.writerCount++
}

func (rwlock *RwLock) Unlock() {
	rwlock.writerCount--
	//fmt.Println("broadcast")
	rwlock.mu.Unlock()
	rwlock.writeCond.Broadcast()
}

func (rwlock *RwLock) RLock() {
	rwlock.mu.Lock()
	//fmt.Println("The value of RLock ", rwlock.writerCount, rwlock.readerCount.Load())
	for rwlock.writerCount > 0 || rwlock.readerCount.Load() >= maxConcurReaders {
		//fmt.Println("inside ", rwlock.writerCount)
		rwlock.writeCond.Wait()
	}

	rwlock.readerCount.Store(rwlock.readerCount.Load() + 1)
	//fmt.Println("wow got it", rwlock.writerCount, rwlock.readerCount.Load())
}

func (rwlock *RwLock) RUnlock() {
	rwlock.readerCount.Store(rwlock.readerCount.Load() - 1)
	//fmt.Println("Reached Runlocks broadcast")
	rwlock.writeCond.Broadcast()
	//fmt.Println("Broadcast says:", rwlock.writerCount, rwlock.readerCount.Load())
	rwlock.mu.Unlock()
}
