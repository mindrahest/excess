package schedule

import (
	"fmt"
	"time"
)

type rrScheduler struct {
	baseScheduler
	quantum time.Duration
}

// newRRScheduler returns a Round Robin scheduler with the time slice, quantum.
func newRRScheduler(quantum time.Duration) *rrScheduler {
	fmt.Println(quantum)
	return &rrScheduler{
		quantum: quantum,
		baseScheduler: baseScheduler{
			runQueue:  make(chan job, queueSize),
			completed: make(chan result, queueSize),
			jobRunner: func(job *job) {
				job.run(0)
			},
		},
	}
}

// schedule schedules the provided jobs in round robin order.
func (s *rrScheduler) schedule(theJobs jobs) {
	running := true
	for running {
		running = false
		for i := range theJobs {
			fmt.Println(theJobs[i])
			if theJobs[i].remaining >= s.quantum {
				theJobs[i].scheduled = s.quantum
				theJobs[i].remaining -= s.quantum

				s.runQueue <- theJobs[i]

				if theJobs[i].remaining > 0 {
					running = true
				}

			} else if theJobs[i].remaining > 0 {
				s.runQueue <- theJobs[i]

			}

		}
	}
	close(s.runQueue)
}

//////
func (s *rrScheduler) run() {
	//TODO(student) Implement run loop
	for job := range s.runQueue {
		job.run(s.quantum)
		//s.jobRunner(&job)
		s.completed <- result{job, time.Since(job.start)}

	}
	close(s.completed)
}

func (s *rrScheduler) results() chan result {
	return s.completed
}

///////
