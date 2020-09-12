package config

import "time"

var AllConfigKes = []string{
	"external.source.node.label",
	"resource.limit.annotation",
	"resource.buffer.init",
	"resource.buffer.factor",
	"resource.limit",
	"extra.env",
	"delete.grace.period",
	"state.validation.period",
	"state.validation.heartbeat",
	"schedulers",
}

func (config *Config) GetStandaloneSourceLabel() string {
	return config.GetStringParam("external.source.node.label", "nodename")
}

func (config *Config) GetResourceLimitAnnotation() string {
	return config.GetStringParam("resource.limit.annotation", "bitflow-resource-limit")
}

func (config *Config) GetInitialResourceBufferSize() int {
	return config.GetIntParam("resource.buffer.init", 2)
}

func (config *Config) GetResourceBufferIncrementFactor() float64 {
	return config.GetFloatParam("resource.buffer.factor", 2.0)
}

func (config *Config) GetDefaultNodeResourceLimit() float64 {
	return config.GetFloatParam("resource.limit", 0.1)
}

func (config *Config) GetPodEnvVars() map[string]string {
	return config.GetStringMapParam("extra.env", map[string]string{})
}

func (config *Config) GetDeleteGracePeriod() time.Duration {
	return config.GetDurationParam("delete.grace.period", 30*time.Second)
}

func (config *Config) GetSpawnPeriod() time.Duration {
	// TODO rename parameter (requires changing many config files)
	return config.GetDurationParam("state.validation.period", 0)
}

func (config *Config) GetReconcileHeartbeat() time.Duration {
	// TODO rename parameter (requires changing many config files)
	return config.GetDurationParam("state.validation.heartbeat", 2*time.Minute)
}

func (config *Config) GetDefaultScheduler() string {
	// TODO rename parameter (requires changing many config files)
	return config.GetStringParam("schedulers", "least-occupied")
}
