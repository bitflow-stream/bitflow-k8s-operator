package config

import (
	"context"
	"strconv"
	"time"

	"github.com/antongulenko/golib"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type cacheEntry struct {
	value           string
	updateTimestamp time.Time
}

type Config struct {
	req           client.Client
	logger        *log.Entry
	namespace     string
	configMapName string

	cache            map[string]cacheEntry
	invalidatePeriod time.Duration

	loggedMissingKeys map[string]bool
}

func NewConfig(client client.Client, namespace, configMapName string) *Config {
	conf := &Config{
		req:               client,
		logger:            newConfigLogger(namespace, configMapName),
		namespace:         namespace,
		configMapName:     configMapName,
		loggedMissingKeys: make(map[string]bool),
	}
	conf.logger.Info("Using direct ConfigMap access")
	return conf
}

func NewConfigWithCache(client client.Client, namespace, configMapName string, invalidatePeriod time.Duration) *Config {
	conf := &Config{
		req:               client,
		logger:            newConfigLogger(namespace, configMapName),
		namespace:         namespace,
		configMapName:     configMapName,
		invalidatePeriod:  invalidatePeriod,
		cache:             make(map[string]cacheEntry),
		loggedMissingKeys: make(map[string]bool),
	}
	conf.logger.Infof("Using cached ConfigMap access (caching period %v)", invalidatePeriod)
	return conf
}

func newConfigLogger(namespace, configMapName string) *log.Entry {
	return log.WithField("config-map", configMapName).WithField("namespace", namespace)
}

func (config *Config) loadConfigMap() (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := config.req.Get(context.TODO(), types.NamespacedName{
		Name:      config.configMapName,
		Namespace: config.namespace,
	}, configMap)
	return configMap, err
}

func (config *Config) Validate() error {
	_, err := config.loadConfigMap()
	return err
}

func (config *Config) SilentlyUseDefaults(keys ...string) {
	for _, key := range keys {
		config.loggedMissingKeys[key] = true
	}
}

func (config *Config) SilentlyUseAllDefaults() {
	config.SilentlyUseDefaults(AllConfigKes...)
}

func (config *Config) getStringParam(key string) (string, bool, error) {
	if config.cache != nil {
		if value, ok := config.cache[key]; ok {
			if time.Now().Sub(value.updateTimestamp) < config.invalidatePeriod {
				return value.value, true, nil
			}
		}
	}
	configMap, err := config.loadConfigMap()
	if err != nil {
		return "", false, err
	}
	if value, ok := configMap.Data[key]; ok {
		if config.cache != nil {
			config.cache[key] = cacheEntry{value, time.Now()}
		}
		return value, true, nil
	}
	return "", false, nil
}

func (config *Config) getParam(key string, defaultValue interface{}, parseValue func(value string) (interface{}, error)) interface{} {
	strValue, found, err := config.getStringParam(key)
	if err != nil {
		config.logger.Warnf("Failed to load config value for key '%v', using default value '%v'. Error: %v", key, defaultValue, err)
	} else if !found {
		if !config.loggedMissingKeys[key] {
			config.logger.Warnf("Missing config value for key '%v', using default value: %v", key, defaultValue)
			config.loggedMissingKeys[key] = true
		}
	} else {
		if result, err := parseValue(strValue); err != nil {
			config.logger.Warnf("Failed to parse config value for key '%v' = '%v', using default value '%v'. Error: %v", key, strValue, defaultValue, err)
		} else {
			return result
		}
	}
	return defaultValue
}

func (config *Config) GetParameterMap() (map[string]string, error) {
	conf, err := config.loadConfigMap()
	if err != nil {
		return nil, err
	}
	return conf.Data, nil
}

func (config *Config) GetStringParam(key string, defaultValue string) string {
	return config.getParam(key, defaultValue, func(strValue string) (interface{}, error) {
		return strValue, nil
	}).(string)
}

func (config *Config) GetIntParam(key string, defaultValue int) int {
	return config.getParam(key, defaultValue, func(strValue string) (interface{}, error) {
		return strconv.Atoi(strValue)
	}).(int)
}

func (config *Config) GetFloatParam(key string, defaultValue float64) float64 {
	return config.getParam(key, defaultValue, func(strVal string) (interface{}, error) {
		return strconv.ParseFloat(strVal, 64)
	}).(float64)
}

func (config *Config) GetDurationParam(key string, defaultValue time.Duration) time.Duration {
	return config.getParam(key, defaultValue, func(strVal string) (interface{}, error) {
		return time.ParseDuration(strVal)
	}).(time.Duration)
}

func (config *Config) GetBoolParam(key string, defaultValue bool) bool {
	return config.getParam(key, defaultValue, func(strVal string) (interface{}, error) {
		return strconv.ParseBool(strVal)
	}).(bool)
}

func (config *Config) GetStringMapParam(key string, defaultValue map[string]string) map[string]string {
	return config.getParam(key, defaultValue, func(strVal string) (interface{}, error) {
		res := golib.ParseMap(strVal)
		return res, nil
	}).(map[string]string)
}

func (config *Config) GetStringSliceParam(key string, defaultValue []string) []string {
	return config.getParam(key, defaultValue, func(strVal string) (interface{}, error) {
		res := golib.ParseSlice(strVal)
		return res, nil
	}).([]string)
}
