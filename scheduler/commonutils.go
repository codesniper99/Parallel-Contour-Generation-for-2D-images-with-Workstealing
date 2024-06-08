package scheduler

import (
	messagepackage "proj3/messagePackage"
	"sync/atomic"
)

type node struct {
	effect messagepackage.Message
}

type TASLock struct {
	value int32
}

func (tasLock *TASLock) lock() {
	for {
		if atomic.SwapInt32(&tasLock.value, 1) == 0 {
			return
		}
	}
}

func (tasLock *TASLock) unlock() {
	atomic.StoreInt32(&tasLock.value, 0)
}

type Queue struct {
	items []node
	lock  TASLock
	size  int
}

func (q *Queue) Enqueue(item node) {
	q.lock.lock()
	defer q.lock.unlock()
	q.items = append(q.items, item)
	q.size = q.size + 1
}

func (q *Queue) Dequeue() node {
	q.lock.lock()
	defer q.lock.unlock()
	if len(q.items) == 0 {
		tmp := node{}
		return tmp
	}
	item := q.items[0]
	q.items = q.items[1:]
	q.size = q.size - 1
	return item
}

func CreateImages(config messagepackage.Config) {

	effects := ReadEffects()

	for _, effect := range effects {
		executeSequentialEffect(config, effect)
	}
}

func CreateImagesTaskQueueAndRun(config messagepackage.Config) {

	effects := ReadEffects()

	taskQueue := Queue{}

	for _, effect := range effects {
		var tmp node
		tmp.effect = effect
		taskQueue.Enqueue(tmp)
	}
	if config.Mode == "s" {
		CreateImages(config)
	} else if config.Mode == "p" {
		ParallelWorkStealingTaskQueueRun(config, taskQueue)
	} else if config.Mode == "chunk" {
		ParallelChunksTaskQueueRun(config, taskQueue)
	}

}
