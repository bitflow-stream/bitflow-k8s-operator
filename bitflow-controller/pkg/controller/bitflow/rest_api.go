package bitflow

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/antongulenko/golib"
	bitflowv1 "github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/apis/bitflow/v1"
	"github.com/bitflow-stream/bitflow-k8s-operator/bitflow-controller/pkg/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	RestApiDataSourcesPath     = "/dataSources"
	RestApiMatchingSourcesPath = "/matchingSources"
)

func (r *BitflowReconciler) startRestApi(listenEndpoint string) {
	ginTask := golib.NewGinTaskWithHandler(listenEndpoint, &golib.GinLogHandler{Logger: golib.Log, Handler: r.handleGinRequestLog})
	r.SetupRestInterface(ginTask.Engine)
	// TODO there is no error handling for the REST API server
	ginTask.Start(nil)
}

func (r *BitflowReconciler) SetupRestInterface(engine *gin.Engine) {
	engine.GET("/health", r.handleHealth)
	engine.GET("/ip", r.handleIpRequest)
	engine.GET("/pods", r.handleRespawningPods)
	if r.statistic != nil {
		engine.GET("/statistics", r.handleStatistics)
	}
	engine.GET(fmt.Sprintf("%v/:podName", RestApiDataSourcesPath), r.handleDataSources)
	engine.GET(RestApiMatchingSourcesPath, r.handleAllMatchingSources)
	engine.GET(fmt.Sprintf("%v/:stepName", RestApiMatchingSourcesPath), r.handleMatchingSources)
}

func (r *BitflowReconciler) handleGinRequestLog(ctx *gin.Context) (log.Level, string, bool) {
	if len(ctx.Errors) == 0 && (ctx.Request.URL.Path == "/health" || strings.HasPrefix(ctx.Request.URL.Path, RestApiDataSourcesPath)) {
		// Hide health requests from the default logs, as they occur frequently
		return log.DebugLevel, "", true
	}
	return log.InfoLevel, "", false
}

func (r *BitflowReconciler) handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{})
}

func (r *BitflowReconciler) handleIpRequest(ctx *gin.Context) {
	// Miniature "what-is-my-ip" service to reliably obtain IPs reachable from the orchestrator
	host, _, err := net.SplitHostPort(ctx.Request.RemoteAddr)
	if err == nil {
		_, err = ctx.Writer.WriteString(host)
	}
	if err != nil {
		panic(err) // Will be caught and logged by the Gin middleware
	}
}

func (r *BitflowReconciler) handleRespawningPods(ctx *gin.Context) {
	podsMap := r.pods.ListRespawningPods()
	ctx.JSON(http.StatusOK, podsMap)
}

func (r *BitflowReconciler) handleStatistics(ctx *gin.Context) {
	data := r.statistic.GetData()
	ctx.JSON(http.StatusOK, data)
}

func (r *BitflowReconciler) handleDataSources(ctx *gin.Context) {
	podName := ctx.Param("podName")
	sources, err := r.listInputSourcesForPod(podName, r.namespace)
	if err != nil {
		_ = ctx.Error(err)
		ctx.Status(http.StatusInternalServerError)
		return
	}
	if len(sources) == 0 {
		ctx.JSON(http.StatusNotFound, []string{})
		return
	}
	sourceUrls := make([]string, len(sources))
	for i, source := range sources {
		sourceUrls[i] = source.Spec.URL
	}
	ctx.JSON(http.StatusOK, sourceUrls)
}

func (r *BitflowReconciler) handleMatchingSources(ctx *gin.Context) {
	stepName := ctx.Param("stepName")
	step, err := common.GetStep(r.client, stepName, r.namespace)
	if err != nil {
		ReplyJSON(ctx, nil, err)
	} else {
		sources, err := r.getMatchingSourceNames(step)
		ReplyJSON(ctx, sources, err)
	}
}

func (r *BitflowReconciler) handleAllMatchingSources(ctx *gin.Context) {
	allSources := make(map[string][]string)
	steps, err := common.GetSteps(r.client, r.namespace)
	if err != nil {
		ReplyJSON(ctx, nil, err)
		return
	}
	for _, step := range steps {
		sources, err := r.getMatchingSourceNames(step)
		if err != nil {
			ReplyJSON(ctx, nil, err)
			return
		}
		allSources[step.Name] = sources
	}
	ReplyJSON(ctx, allSources, err)
}

func (r *BitflowReconciler) getMatchingSourceNames(step *bitflowv1.BitflowStep) ([]string, error) {
	sourceObjects, err := r.listMatchingSources(step)
	if err != nil {
		return nil, fmt.Errorf("Failed to list matching data sources for %v: %v", step, err)
	}
	sources := make([]string, len(sourceObjects))
	for i, sourceObj := range sourceObjects {
		sources[i] = sourceObj.Name
	}
	return sources, nil
}

func (r *BitflowReconciler) listInputSourcesForPod(podName, namespace string) ([]*bitflowv1.BitflowSource, error) {
	var pod corev1.Pod
	err := r.client.Get(context.TODO(), client.ObjectKey{Name: podName, Namespace: namespace}, &pod)
	if err != nil {
		if errors.IsNotFound(err) {
			log.WithField("pod", podName).Debugln("Failed to fetch input sources for pod:", err)
			return nil, nil
		}
		return nil, fmt.Errorf("Error fetching pod '%v': %v", podName, err)
	}
	stepName := pod.Labels[bitflowv1.LabelStepName]
	if stepName == "" {
		return nil, fmt.Errorf("pod '%v' has no valid '%v' label", podName, bitflowv1.LabelStepName)
	}
	step, err := common.GetStep(r.client, stepName, namespace)
	if err != nil {
		if errors.IsNotFound(err) {
			log.WithField("pod", podName).WithField("step", stepName).Debugln("Failed to fetch step:", err)
			return nil, nil
		}
		return nil, err
	}

	if step.Type() == bitflowv1.StepTypeOneToOne {
		sourceName := pod.Labels[bitflowv1.PodLabelOneToOneSourceName]
		if sourceName == "" {
			return nil, fmt.Errorf("pod '%v' has no valid '%v' label", podName, bitflowv1.PodLabelOneToOneSourceName)
		}
		source, err := common.GetSource(r.client, sourceName, namespace)
		if err != nil {
			if errors.IsNotFound(err) {
				log.WithField("pod", podName).WithField("source", sourceName).Debugln("Failed to fetch one-to-one input source for pod:", err)
				return nil, nil
			}
			return nil, fmt.Errorf("Error fetching source '%v': %v", sourceName, err)
		}
		return []*bitflowv1.BitflowSource{source}, nil
	} else {
		return r.listMatchingSources(step)
	}
}

func ReplyError(c *gin.Context, statusCode int, err error) {
	c.JSON(statusCode, map[string]string{"error": err.Error()})
}

func ReplyJSON(c *gin.Context, obj interface{}, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		if errors.IsNotFound(err) {
			status = http.StatusNotFound
		}
		ReplyError(c, status, err)
	} else if obj != nil {
		c.JSON(http.StatusOK, obj)
	} else {
		c.Status(http.StatusOK)
	}
}
