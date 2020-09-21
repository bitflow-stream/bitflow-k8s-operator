package config

import "time"

var AllConfigKes = []string{
	"external.source.label",
	"resource.limit.node.annotation",
	"resource.limit.slots",
	"resource.limit.slots.grow",
	"resource.limit",
	"extra.env",
	"delete.grace.period",
	"pod.spawn.period",
	"reconcile.heartbeat",
	"scheduler",
}

func (config *Config) GetStandaloneSourceLabel() string {
	return config.GetStringParam("external.source.label", "bitflow-nodename")
}

func (config *Config) GetDefaultNodeResourceLimit() float64 {
	return config.GetFloatParam("resource.limit", 0.1)
}

func (config *Config) GetResourceLimitAnnotation() string {
	return config.GetStringParam("resource.limit.node.annotation", "bitflow-resource-limit")
}

func (config *Config) GetInitialResourceBufferSize() int {
	return config.GetIntParam("resource.limit.slots", 5)
}

func (config *Config) GetResourceBufferIncrementFactor() float64 {
	return config.GetFloatParam("resource.limit.slots.grow", 2.0)
}

func (config *Config) GetPodEnvVars() map[string]string {
	return config.GetStringMapParam("extra.env", map[string]string{})
}

func (config *Config) GetDeleteGracePeriod() time.Duration {
	return config.GetDurationParam("delete.grace.period", 20*time.Second)
}

func (config *Config) GetSpawnPeriod() time.Duration {
	return config.GetDurationParam("pod.spawn.period", 3*time.Second)
}

func (config *Config) GetReconcileHeartbeat() time.Duration {
	return config.GetDurationParam("reconcile.heartbeat", 30*time.Second)
}

func (config *Config) GetSchedulerName() string {
	return config.GetStringParam("scheduler", "least-occupied")
}
