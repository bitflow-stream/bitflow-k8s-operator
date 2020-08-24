package scheduler

import "errors"

type EqualDistributionScheduler struct {
	NodeNames []string
	PodNames  []string
}

func (eds EqualDistributionScheduler) Schedule() (bool, map[string]string, error) {
	if err := eds.validate(eds); err != nil {
		return false, nil, err
	}

	m := make(map[string]string)

	if len(eds.PodNames) == 0 {
		return true, m, nil
	}

	var nodeIndex = 0
	for _, podName := range eds.PodNames {
		if nodeIndex >= len(eds.NodeNames) {
			nodeIndex = 0
		}

		m[podName] = eds.NodeNames[nodeIndex]

		nodeIndex++
	}

	return true, m, nil
}

func (eds EqualDistributionScheduler) validate(scheduler EqualDistributionScheduler) error {
	if len(scheduler.NodeNames) == 0 {
		return errors.New("no nodes in scheduler")
	}
	for _, name := range scheduler.NodeNames {
		if name == "" {
			return errors.New("empty name in NodeNames")
		}
	}
	if len(scheduler.PodNames) == 0 {
		return errors.New("no pods in scheduler")
	}
	for _, name := range scheduler.PodNames {
		if name == "" {
			return errors.New("empty name in PodNames")
		}
	}
	return nil
}
