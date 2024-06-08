package scheduler

// _____________ bounded queueue
import (
	"math/rand"
	messagepackage "proj3/messagePackage"
	"sync/atomic"
)

type AtomicStampedReferenceInt struct {
	val   int32
	stamp int32
}
type BoundedDEQue struct {
	tasks  []*messagepackage.Message
	bottom atomic.Int32
	top    AtomicStampedReferenceInt
}

var checker []int

func (q *BoundedDEQue) size() int32 {
	if q.bottom.Load() >= q.top.val {
		return q.bottom.Load() - q.top.val
	} else {
		return 0
	}
}
func (asr *AtomicStampedReferenceInt) compareAndSet(expectVal int32, newVal int32, expectStamp int32, newStamp int32) bool {
	return atomic.CompareAndSwapInt32(&asr.val, expectVal, newVal) && atomic.CompareAndSwapInt32(&asr.stamp, expectStamp, newStamp)
}

func (asr *AtomicStampedReferenceInt) set(newVal int32, newStamp int32) {
	asr.val = newVal
	asr.stamp = newStamp
}

func NewBoundedDEQue(capacity int) *BoundedDEQue {
	return &BoundedDEQue{
		tasks:  make([]*messagepackage.Message, capacity),
		bottom: atomic.Int32{},
		top: AtomicStampedReferenceInt{
			val:   int32(0),
			stamp: int32(0),
		},
	}
}

func (q *BoundedDEQue) PushBottom(r *messagepackage.Message) {
	q.tasks[q.bottom.Load()] = r
	q.bottom.Add(1)
}

func (q *BoundedDEQue) isEmpty() bool {
	return q.top.val < q.bottom.Load()
}

type WorkStealingThread struct {
	Queue []*BoundedDEQue
	id    int
}

func NewWorkStealingThread(bigQueue []*BoundedDEQue, id int) *WorkStealingThread {
	return &WorkStealingThread{Queue: bigQueue, id: id}
}

func (w *WorkStealingThread) Run(config messagepackage.Config, jobs *atomic.Int32, totalJobs int32) {

	me := w.id
	for jobs.Load() < totalJobs || w.Queue[me].size() > 0 {
		task := w.Queue[me].popBottom()
		if task != nil {
			ProcessParallelWorkStealingChunk(*task, config.Threshold, config.ThreadCount)
			jobs.Add(1)
		}
		size := w.Queue[me].size()
		if size == 0 {
			victim := rand.Intn(config.ThreadCount)

			theft := w.Queue[victim].popTop()
			if theft != nil {
				w.Queue[me].PushBottom(theft)
			}
		}

	}
}

type RecursiveAction struct {
	// Your RecursiveAction struct fields go here
	val int
}

func (q *BoundedDEQue) popTop() *messagepackage.Message {

	oldTop := q.top.val
	newTop := oldTop + 1
	oldStamp := q.top.stamp
	newStamp := oldStamp + 1
	if q.bottom.Load() <= oldTop {
		return nil
	}
	r := q.tasks[oldTop]
	if q.top.compareAndSet(oldTop, newTop, oldStamp, newStamp) {
		return r
	}
	return nil
}

func (q *BoundedDEQue) popBottom() *messagepackage.Message {
	if q.bottom.Load() == 0 {
		return nil
	}
	newBottom := q.bottom.Add(-1)
	r := q.tasks[newBottom]
	oldTop := q.top.val
	newTop := 0
	oldStamp := q.top.stamp
	newStamp := oldStamp + 1

	if newBottom > oldTop {
		return r
	}

	if newBottom == oldTop {
		q.bottom.Store(0)
		if q.top.compareAndSet(oldTop, int32(newTop), oldStamp, newStamp) {
			return r
		}
	}
	q.top.set(int32(newTop), newStamp)
	return nil
}
