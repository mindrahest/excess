package schedule

import "sort"

type sjfScheduler struct {
	baseScheduler
}

// newSJFScheduler returns a shortest job first scheduler.
// With this scheduler, jobs are executed in the order of shortest job first.
func newSJFScheduler() *sjfScheduler {
	return &sjfScheduler{
		baseScheduler: baseScheduler{
			runQueue:  make(chan job, queueSize),
			completed: make(chan result, queueSize),
			jobRunner: func(job *job) {
				job.run(0)
			},
		},
	}
}

// schedule schedules the provided jobs according to SJF order.
// The tasks with the lowest estimate is scheduled to run first.
func (s *sjfScheduler) schedule(inJobs jobs) {
	sort.Slice(inJobs, func(i, j int) bool {
		return inJobs[i].estimated < inJobs[j].estimated
	})
	//	if inJobs[i].estimated > inJobs[j].estimated {
	//		inJobs[i], inJobs[j] = inJobs[j], inJobs[i]
	//	}

	for _, job := range inJobs {

		s.runQueue <- job
	}

	close(s.runQueue)
}
