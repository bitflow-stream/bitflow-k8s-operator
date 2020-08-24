package scheduler

type Scheduler interface {
	Schedule() (bool, map[string]string, error)
}
