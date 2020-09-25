package schedule

import (
	"fmt"
	"time"
)

// strideScheduler may be defined as an alias for rrScheduler; it has the same fields.
type strideScheduler struct {
	baseScheduler
	quantum time.Duration
}

// newStrideScheduler returns a stride scheduler.
// With this scheduler, jobs are executed similar to round robin,
// but with exact proportions determined by how many tickets each job is assigned.
func newStrideScheduler(quantum time.Duration) *strideScheduler {
	return &strideScheduler{
		quantum: quantum,
		baseScheduler: baseScheduler{
			runQueue:  make(chan job, queueSize),
			completed: make(chan result, queueSize),
			jobRunner: func(job *job) {
				job.run(job.scheduled)
			},
		},
	}
}

// schedule schedules the provided jobs according to a stride scheduler's order.
// The task with the lowest pass is scheduled to run first.
func (s *strideScheduler) schedule(theJobs jobs) {

	running := true

	for running {
		running = false

		for i := 0; i < len(theJobs); i++ {
			it := minPass(theJobs)

			if theJobs[it].remaining >= s.quantum {

				theJobs[it].scheduled = s.quantum
				theJobs[it].remaining -= s.quantum

				s.runQueue <- theJobs[it]
				fmt.Println(theJobs[it])
				fmt.Println("quantum")
				if theJobs[it].remaining > 0 {
					running = true
				}

			} else if theJobs[it].remaining > 0 {
				s.runQueue <- theJobs[it]
				fmt.Println(theJobs[it])
				fmt.Println("remaining")
			}
			if theJobs[it].remaining == 0 {
				theJobs = append(theJobs[:i], theJobs[i+1:]...)
			}

		}

	}
	close(s.runQueue)
}

var lowestPass int

// minPass returns the index of the job with the lowest pass value.

func minPass(theJobs jobs) int {

	lowestIndex := 0

	lowestStride := 0

	for i := range theJobs {

		if lowestStride == 0 {
			lowestStride = theJobs[i].stride
		}
		if theJobs[i].pass < lowestPass {
			lowestPass = theJobs[i].pass
			lowestIndex = i
			continue
		}

		if theJobs[i].pass == lowestPass {

			if theJobs[i].stride < lowestStride && theJobs[i].remaining > 0 {

				lowestPass = theJobs[i].pass
				lowestStride = theJobs[i].stride
				lowestIndex = i

			}

		}

	}

	theJobs[lowestIndex].pass += theJobs[lowestIndex].stride

	//TODO(student) implement minPass and use it from schedule()
	// You need to keep track of both the lowest pass value and its index.

	return lowestIndex
}

func (s *strideScheduler) run() {
	//TODO(student) Implement run loop
	for job := range s.runQueue {
		//job.run(s.quantum)
		s.jobRunner(&job)
		s.completed <- result{job, time.Since(job.start)}

	}
	close(s.completed)
}

func (s *strideScheduler) results() chan result {
	return s.completed
}

/* if theJobs[i].pass == theJobs[i+1].pass && theJobs[i].pass == theJobs[i+2].pass {
println("2nd if")
lowestIndex = i
theJobs[i].pass += theJobs[i].stride */
