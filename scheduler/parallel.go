package scheduler

import (
	"math"
	"proj3/customwg"
	messagepackage "proj3/messagePackage"
	"sync/atomic"
)

type MessageCopy struct {
	InPath  string
	OutPath string
	Effects []string
}

func ProcessParallelWorkStealingChunk(msg messagepackage.Message, threshold int, threadCount int) {
	executeParallelWorkStealingEffects(msg, threshold, threadCount)
}

func ParallelWorkStealingTaskQueueRun(config messagepackage.Config, taskQueue Queue) {

	numWorkers := config.ThreadCount
	workers := make([]*WorkStealingThread, numWorkers)
	x := float64(taskQueue.size) / float64(numWorkers)

	numsPerWorker := math.Ceil(x)
	queueSize := 2 * numsPerWorker
	bigWorkQueue := make([]*BoundedDEQue, numWorkers)

	for i := 0; i < numWorkers; i++ {
		tmpQueue := NewBoundedDEQue(int(queueSize))
		bigWorkQueue[i] = tmpQueue
		workers[i] = NewWorkStealingThread(bigWorkQueue, i)
	}

	totalJobs := int32(taskQueue.size)
	//fmt.Println("Total jobs = ", totalJobs)
	var wg customwg.CustomWaitGroup = *customwg.NewWaitGroup()
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		for j := 0; j < int(numsPerWorker); j++ {
			if taskQueue.size > 0 {
				tmp := taskQueue.Dequeue().effect
				msg := &messagepackage.Message{InPath: tmp.InPath, OutPath: tmp.OutPath, Effects: tmp.Effects}
				workers[i].Queue[i].PushBottom(msg)
			} else {
				break
			}
		}
		if taskQueue.size <= 0 {
			break
		}
	}

	var jobs atomic.Int32
	jobs.Store(0)
	for i := 0; i < numWorkers; i++ {
		go func(workerIndex int) {
			defer wg.Done()

			workers[workerIndex].Run(config, &jobs, totalJobs)
		}(i)
	}

	wg.Wait()
}

func ParallelChunksTaskQueueRun(config messagepackage.Config, taskQueue Queue) {

	numWorkers := config.ThreadCount

	numberOfImages := len(taskQueue.items)
	for i := 0; i < numberOfImages; i++ {
		frontNode := taskQueue.Dequeue()
		executeParallelChunkEffects(frontNode.effect, config.Threshold, numWorkers)
	}
}
