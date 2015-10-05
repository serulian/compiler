// Copyright 2015 The Serulian Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package compilerutil

import (
	"container/list"
	"sync"
)

// orderedWorkQueue defines a queue for processing work in a dependency-preserving fashion.
//
// Note that tasks are run concurrently and only limited by their explicit dependencies and
// the go scheduler. Care must be taken to ensure that concurrent data structures are used for
// any shared state (which itself should be minimized if at all possible)
type orderedWorkQueue struct {
	running bool
	result  *orderedWorkQueueResult // Result of the run.

	jobs        *list.List           // The jobs. Only to be accessed by Enqueue and runQueue
	keyCounter  map[interface{}]int  // The map of keys to their count. Only to be accessed by Enqueue and runQueue.
	runningJobs map[interface{}]bool // The running jobs, by key. Only to be accessed by runQueue.

	terminated  chan bool
	workTracker sync.WaitGroup
}

// A work item in the ordered work queue.
type orderedWorkQueueItem struct {
	key            interface{}   // The job key.
	value          interface{}   // The job's extra data value.
	dependencyKeys []interface{} // If non-empty, the keys of the jobs that must be complete before this job run.
	performer      performFn     // The function to execute.
}

// orderedWorkQueueResult represents a result of a run of an orderedWorkQueue
type orderedWorkQueueResult struct {
	Status   bool          // Status of the work performed. Will be the combination of the return values of all the jobs.
	HasCycle bool          // Whether a cycle was detected.
	Cycle    []interface{} // If a cycle was detected, the keys of the remaining jobs when it was detected.
}

// Queue returns a new ordered work queue.
func Queue() *orderedWorkQueue {
	return &orderedWorkQueue{
		running:    false,
		jobs:       list.New(),
		terminated: make(chan bool),

		keyCounter:  map[interface{}]int{},
		runningJobs: map[interface{}]bool{},
	}
}

// performFn performs work over the given enqueued key and value. The function should return
// whether the work completed successfully.
type performFn func(key interface{}, value interface{}) bool

// Enqueue adds a work item to be performed. The key argument gives a unique name to the job, preventing any other
// jobs with the same name from being run concurrently. The value is extra data to be passed to the job.
// dependencyKeys specify the keys of other jobs that must *all* complete before this job begins.
func (q *orderedWorkQueue) Enqueue(key interface{}, value interface{}, performer performFn, dependencyKeys ...interface{}) {
	if q.running {
		panic("Cannot enqueue to an already running work queue")
	}

	item := orderedWorkQueueItem{
		key:            key,
		value:          value,
		dependencyKeys: dependencyKeys,
		performer:      performer,
	}

	q.jobs.PushBack(item)
	q.keyCounter[key] = q.keyCounter[key] + 1
}

// Run performs all the work in the queue, blocking until completion and returns true only if
// ALL functions run returned true.
func (q *orderedWorkQueue) Run() *orderedWorkQueueResult {
	// Mark the queue as running.
	q.running = true
	q.result = &orderedWorkQueueResult{true, false, make([]interface{}, 0)}

	// Run the internal gorountine.
	go q.runQueue()

	// Wait for the run to terminate.
	<-q.terminated
	return q.result
}

// performJob performs the specified job. This is called under a goroutine.
func (q *orderedWorkQueue) performJob(job orderedWorkQueueItem) {
	if !job.performer(job.key, job.value) {
		q.result.Status = false
	}

	// Tell runQueue that we're ready for more work.
	q.workTracker.Done()
}

// runQueue performs the internal running of the queue.
func (q *orderedWorkQueue) runQueue() {
	for {
		if q.jobs.Len() == 0 {
			// All done!
			q.terminated <- true
			return
		}

		// For each job in the queue, find any without any outstanding dependencies and start a single instance
		// of each key.
		var jobsStarted = make([]orderedWorkQueueItem, 0)
		var next *list.Element

	itemLoop:
		for e := q.jobs.Front(); e != nil; e = next {
			job := e.Value.(orderedWorkQueueItem)
			next = e.Next()

			// Check to ensure that there are no jobs with keys matching the dependency keys.
			for _, dependency := range job.dependencyKeys {
				if _, exists := q.keyCounter[dependency]; exists {
					continue itemLoop
				}
			}

			// Check to ensure that a job with the same key isn't already running.
			if _, exists := q.runningJobs[job.key]; exists {
				continue itemLoop
			}

			// Add the job to the running list.
			q.runningJobs[job.key] = true

			// Remove the job from the items list.
			q.jobs.Remove(e)

			// Spin off a worker to run the job.
			q.workTracker.Add(1)
			jobsStarted = append(jobsStarted, job)

			go q.performJob(job)
		}

		// If we have reached this point and no jobs have been started, then we have a circular dependency.
		if len(jobsStarted) == 0 {
			// Nothing more to do.
			q.result.HasCycle = true
			q.result.Cycle = make([]interface{}, q.jobs.Len())

			var index = 0
			for e := q.jobs.Front(); e != nil; e = e.Next() {
				q.result.Cycle[index] = e.Value.(orderedWorkQueueItem).key
				index++
			}

			q.terminated <- true
			return
		}

		// Wait for the jobs to finish.
		q.workTracker.Wait()

		// Remove the run jobs from the key counter map.
		for _, jobRun := range jobsStarted {
			count, _ := q.keyCounter[jobRun.key]
			if count == 1 {
				delete(q.keyCounter, jobRun.key)
			} else {
				q.keyCounter[jobRun.key] = count - 1
			}
		}

		// Reset the running jobs list.
		q.runningJobs = map[interface{}]bool{}
	}
}
