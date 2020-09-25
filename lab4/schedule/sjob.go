package schedule

import "time"

// newSJob creates a job entry for stride scheduling.
func newSJob(id, tickets int, estimated time.Duration) job {
	const numerator = 10000
	//TODO(student) finish the sjob constructor

	return job{
		id:        id,
		estimated: estimated,
		scheduled: estimated,
		remaining: estimated,
		doJob:     time.Sleep,
		tickets:   tickets,
		pass:      0,
		stride:    numerator / tickets,
	}
}
