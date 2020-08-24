package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/antongulenko/golib"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/controller"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/version"
	"github.com/operator-framework/operator-sdk/pkg/k8sutil"
	kubemetrics "github.com/operator-framework/operator-sdk/pkg/kube-metrics"
	"github.com/operator-framework/operator-sdk/pkg/leader"
	"github.com/operator-framework/operator-sdk/pkg/metrics"
	"github.com/operator-framework/operator-sdk/pkg/restmapper"
	sdkVersion "github.com/operator-framework/operator-sdk/version"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
)

type Main struct {
	lockName            string
	metricsHost         string
	metricsPort         int
	operatorMetricsPort int
}

func main() {
	var mgr Main
	os.Exit(mgr.executeMain())
}

func (main *Main) executeMain() int {
	flag.StringVar(&main.lockName, "lock-name", "bitflow-operator-lock", "Cluster-unique name of a lock that is acquired when starting")
	flag.StringVar(&main.metricsHost, "metrics-host", "0.0.0.0", "Binding host for serving metrics API")
	flag.IntVar(&main.metricsPort, "metrics-port", 8383, "Binding port for serving metrics API")
	flag.IntVar(&main.operatorMetricsPort, "operator-metrics-port", 8686, "Binding port for serving operator-specific metrics API")

	golib.RegisterLogFlags()
	flag.Parse()
	golib.ConfigureLogging()
	log.Infoln("Bitflow operator version:", version.Version)
	log.Infof("Go Version: %s, Go OS/Arch: %s/%s, Version of operator-sdk: %v", runtime.Version(), runtime.GOOS, runtime.GOARCH, sdkVersion.Version)

	namespace, err := k8sutil.GetWatchNamespace()
	if err != nil {
		log.Errorln("Failed to get watch namespace:", err)
		return 1
	}

	// Get a config to talk to the API server
	cfg, err := config.GetConfig()
	if err != nil {
		log.Errorln("Failed to load Kube config:", err)
		return 1
	}

	// Become the leader before proceeding
	ctx := context.TODO()
	log.Infof("Acquiring leader lock '%v'...", main.lockName)
	if err := leader.Become(ctx, main.lockName); err != nil {
		log.Errorf("Failed to acquire leader lock '%v': %v", main.lockName, err)
		return 1
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(cfg, manager.Options{
		Namespace:          namespace,
		MapperProvider:     restmapper.NewDynamicRESTMapper,
		MetricsBindAddress: fmt.Sprintf("%s:%d", main.metricsHost, main.metricsPort),
	})
	if err != nil {
		log.Errorln("Failed to create manager:", err)
		return 1
	}

	log.Infoln("Registering components")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Errorln("Failed to register Kubernetes schemes:", err)
		return 1
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr, namespace); err != nil {
		log.Errorln("Failed to register manager functions:", err)
		return 1
	}

	// Serve default metrics and metrics about custom resource objects
	if err := main.serveCRMetrics(cfg); err != nil {
		log.Warnln("Failed to generate and serve custom resource metrics:", err)
	}
	if err := main.serveMetrics(namespace, ctx, cfg); err != nil {
		log.Warnln("Failed to serve metrics:", err)
	}

	log.Infoln("Starting the operator")
	if err := mgr.Start(signals.SetupSignalHandler()); err != nil {
		log.Errorln("Operator failed:", err)
		return 1
	}
	return 0
}

// serveCRMetrics gets the Operator/CustomResource GVKs and generates metrics based on those types.
// It serves those metrics on "http://metricsHost:8383".
func (main *Main) serveCRMetrics(cfg *rest.Config) error {
	// Below function returns filtered operator/CustomResource specific GVKs.
	// For more control override the below GVK list with your own custom logic.
	filteredGVK, err := k8sutil.GetGVKsFromAddToScheme(apis.AddToScheme)
	if err != nil {
		return err
	}
	// Get the namespace the operator is currently deployed in.
	operatorNs, err := k8sutil.GetOperatorNamespace()
	if err != nil {
		return err
	}
	// To generate metrics in other namespaces, add the values below.
	ns := []string{operatorNs}
	// Generate and serve custom resource specific metrics.
	err = kubemetrics.GenerateAndServeCRMetrics(cfg, ns, filteredGVK, main.metricsHost, int32(main.operatorMetricsPort))
	if err != nil {
		return err
	}
	return nil
}

func (main *Main) serveMetrics(namespace string, ctx context.Context, cfg *rest.Config) error {
	// Put to the below struct any other metrics ports you want to expose.
	servicePorts := []v1.ServicePort{
		{Port: int32(main.metricsPort), Name: metrics.OperatorPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: int32(main.metricsPort)}},
		{Port: int32(main.operatorMetricsPort), Name: metrics.CRPortName, Protocol: v1.ProtocolTCP, TargetPort: intstr.IntOrString{Type: intstr.Int, IntVal: int32(main.operatorMetricsPort)}},
	}
	// Create Service object to expose the metrics port(s).
	service, err := metrics.CreateMetricsService(ctx, cfg, servicePorts)
	if err != nil {
		return fmt.Errorf("Failed to create metrics service: %v", err)
	}

	// CreateServiceMonitors will automatically create the prometheus-operator ServiceMonitor resources
	// necessary to configure Prometheus to scrape metrics from this operator.
	services := []*v1.Service{service}
	_, err = metrics.CreateServiceMonitors(cfg, namespace, services)
	if err != nil {
		// If this operator is deployed to a cluster without the prometheus-operator running, it will return
		// ErrServiceMonitorNotPresent, which can be used to safely skip ServiceMonitor creation.
		if err == metrics.ErrServiceMonitorNotPresent {
			err = fmt.Errorf("Install Prometheus operator in your cluster to create ServiceMonitor objects. Error: %v", err)
		} else {
			err = fmt.Errorf("Could not create ServiceMonitor object: %v", err)
		}
	}
	return err
}
