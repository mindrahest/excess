package schedule

type baseScheduler struct {
	runQueue  chan job
	completed chan result
	jobRunner func(*job)
}

// jobs is a slice of jobs ordered according to some scheduling policies.
type jobs []job

// run starts executing the scheduled jobs from the run queue.
func (s *baseScheduler) run() {
	//TODO(student) Implement run loop
	for job := range s.runQueue {
		s.jobRunner(&job)

		s.completed <- result{job, job.scheduled}

	}
	close(s.completed)
}

// results returns the channel of results.
// This is primarily used for testing.
func (s *baseScheduler) results() chan result {
	return s.completed
}
