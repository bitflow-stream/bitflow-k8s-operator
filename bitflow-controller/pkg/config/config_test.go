package config

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const (
	configName = "test-configuration"

	boolValue   = true
	floatValue  = 6.0
	intValue    = 3
	stringValue = "HelloWorld"

	updateBoolValue   = false
	updateFloatValue  = 20.1
	updateIntValue    = 101
	updateStringValue = "101"

	boolTest   = "helloBool"
	floatTest  = "helloFloat32"
	intTest    = "helloInt"
	stringTest = "helloString"

	errorInt   = "notParsableInt"
	errorFloat = "notParsableFloat"
	errorBool  = "notParsableBool"
)

type ConfigTestSuite struct {
	common.AbstractTestSuite
}

func TestConfig(t *testing.T) {
	new(ConfigTestSuite).Run(t)
}

func (s *ConfigTestSuite) UpdateConfig(client client.Client, name string) {
	conf := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{
		Name:      name,
		Namespace: common.TestNamespace,
	}, conf)
	s.NoError(err)

	conf.Data[boolTest] = strconv.FormatBool(updateBoolValue)
	conf.Data[floatTest] = strconv.FormatFloat(updateFloatValue, 'f', -1, 64)
	conf.Data[intTest] = strconv.Itoa(updateIntValue)
	conf.Data[stringTest] = updateStringValue
	s.NoError(client.Update(context.TODO(), conf))
}

func (s *ConfigTestSuite) GetConfigMap(name string) *corev1.ConfigMap {
	data := make(map[string]string)
	data[boolTest] = strconv.FormatBool(boolValue)
	data[floatTest] = strconv.FormatFloat(floatValue, 'f', -1, 64)
	data[intTest] = strconv.Itoa(intValue)
	data[stringTest] = stringValue
	data[errorFloat] = errorFloat
	data[errorInt] = errorInt
	data[errorBool] = errorBool

	var config corev1.ConfigMap
	config.Name = name
	config.Namespace = common.TestNamespace
	config.Data = data
	return &config
}

func (s *ConfigTestSuite) TestConfig() {
	configMap := s.GetConfigMap(configName)
	cl := fake.NewFakeClient(configMap)
	conf := NewConfig(cl, common.TestNamespace, configName)
	s.testSavedValues(conf)
	s.testDefaultValues(conf)
	s.testErrorValues(conf)

	s.UpdateConfig(cl, configName)
	s.testUpdatedValues(conf)
}

func (s *ConfigTestSuite) TestConfigCache() {
	configMap := s.GetConfigMap(configName)
	cl := fake.NewFakeClient(configMap)
	conf := NewConfigWithCache(cl, common.TestNamespace, configName, 2*time.Second)
	s.testSavedValues(conf)
	s.testDefaultValues(conf)
	s.testErrorValues(conf)

	s.UpdateConfig(cl, configName)
	s.testSavedValues(conf)
	time.Sleep(3 * time.Second)
	s.testUpdatedValues(conf)
}

func (s *ConfigTestSuite) TestConfigValidate() {
	configMap := s.GetConfigMap(configName)
	cl := fake.NewFakeClient(configMap)
	conf := NewConfigWithCache(cl, common.TestNamespace, configName, 2*time.Second)
	s.NoError(conf.Validate())

	confFail := NewConfig(cl, common.TestNamespace, "doesNotExist")
	s.Error(confFail.Validate())

	parameterMap, err := conf.GetParameterMap()
	s.NoError(err)
	s.Len(parameterMap, 7)
}

func (s *ConfigTestSuite) testSavedValues(conf *Config) {
	s.Equal(boolValue, conf.GetBoolParam(boolTest, false))
	s.Equal(floatValue, conf.GetFloatParam(floatTest, 40.4))
	s.Equal(intValue, conf.GetIntParam(intTest, 404))
	s.Equal(stringValue, conf.GetStringParam(stringTest, "404"))
}

func (s *ConfigTestSuite) testDefaultValues(conf *Config) {
	s.Equal(false, conf.GetBoolParam("doesNotExist", false))
	s.Equal(40.4, conf.GetFloatParam("doesNotExist", 40.4))
	s.Equal(404, conf.GetIntParam("doesNotExist", 404))
	s.Equal("404", conf.GetStringParam("doesNotExist", "404"))
}

func (s *ConfigTestSuite) testErrorValues(conf *Config) {
	s.Equal(false, conf.GetBoolParam(errorBool, false))
	s.Equal(40.4, conf.GetFloatParam(errorFloat, 40.4))
	s.Equal(404, conf.GetIntParam(errorInt, 404))
}

func (s *ConfigTestSuite) testUpdatedValues(conf *Config) {
	s.Equal(updateBoolValue, conf.GetBoolParam(boolTest, false))
	s.Equal(updateFloatValue, conf.GetFloatParam(floatTest, 40.4))
	s.Equal(updateIntValue, conf.GetIntParam(intTest, 404))
	s.Equal(updateStringValue, conf.GetStringParam(stringTest, "404"))
}
