package controller

import (
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/controller/bitflow"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// AddToManager adds all Controllers to the Manager
func AddToManager(m manager.Manager, watchNamespace string) error {
	return bitflow.Add(m, watchNamespace)
}
