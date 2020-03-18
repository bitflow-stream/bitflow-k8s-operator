package common

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetKubernetesClient(kubernetesConfig string) (client.Client, error) {
	kubeConf, err := GetKubernetesConfig(kubernetesConfig)
	if err != nil {
		return nil, err
	}
	return client.New(kubeConf, client.Options{Scheme: registerScheme()})
}

func GetKubernetesConfig(kubernetesConfig string) (*rest.Config, error) {
	if kubernetesConfig == "" {
		log.Println("Loading in-cluster Kubernetes Client configuration")
		return rest.InClusterConfig()
	} else {
		log.Println("Loading Kubernetes Client configuration from file:", kubernetesConfig)
		return clientcmd.BuildConfigFromFlags("", kubernetesConfig)
	}
}

func registerScheme() *runtime.Scheme {
	s := scheme.Scheme
	step := &bitflowv1.BitflowStep{}
	source := &bitflowv1.BitflowSource{}
	stepList := &bitflowv1.BitflowStepList{}
	sourceList := &bitflowv1.BitflowSourceList{}
	s.AddKnownTypes(bitflowv1.SchemeGroupVersion, step, source, stepList, sourceList)
	metav1.AddToGroupVersion(s, bitflowv1.SchemeGroupVersion)
	return s
}
