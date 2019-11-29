package kubethanos

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
	"kubethanos/thanos"
	"math/rand"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/tools/reference"
)

type KubeThanos struct {
	Client           kubernetes.Interface
	Namespaces       labels.Selector
	IncludedPodNames *regexp.Regexp
	ExcludedPodNames *regexp.Regexp
	Thanos           *thanos.Thanos
	PercentageToKill float64
	DryRun           bool
	EventRecorder    record.EventRecorder
}

var logger = log.StandardLogger()

var podNotFound = "pod not found"
var errPodNotFound = errors.New(podNotFound)

func New(client kubernetes.Interface, namespaces labels.Selector, includedPodNames, excludedPodNames *regexp.Regexp, percentageToKill float64, dryRun bool, thanos *thanos.Thanos) *KubeThanos {
	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: client.CoreV1().Events(v1.NamespaceAll)})
	recorder := broadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "thanos"})

	return &KubeThanos{
		Client:           client,
		Namespaces:       namespaces,
		IncludedPodNames: includedPodNames,
		ExcludedPodNames: excludedPodNames,
		DryRun:           dryRun,
		PercentageToKill: percentageToKill,
		Thanos:           thanos,
		EventRecorder:    recorder,
	}
}

func (kubeThanos *KubeThanos) Run(ctx context.Context, next <-chan time.Time) {
	for {
		var err = kubeThanos.KillPods()
		if err != nil {
			logger.WithField("err", err).Error("failed to kill pod")
		}

		select {
		case <-next:
		case <-ctx.Done():
			return
		}
	}
}

func (kubeThanos *KubeThanos) KillPods() error {
	podsToKill, err := kubeThanos.SelectPodsToKill()

	if err != nil {
		return err
	}

	if err == errPodNotFound {
		logger.Debug(podNotFound)
		return nil
	}

	var result *multierror.Error
	for _, victim := range podsToKill {
		err = kubeThanos.DeletePod(victim)
		if err != nil {
			logger.WithFields(log.Fields{"err": err,
				"pod": victim}).Error("Failed to delete pod...")
		}
	}

	return result.ErrorOrNil()
}

func (kubeThanos *KubeThanos) SelectPodsToKill() ([]v1.Pod, error) {
	pods, err := kubeThanos.SelectCandidatePods()
	if err != nil {
		return []v1.Pod{}, err
	}

	if len(pods) == 0 {
		return []v1.Pod{}, errPodNotFound
	}

	logger.WithFields(log.Fields{
		"size": len(pods),
	}).Info("Total pods:")

	pods = RandomPodSlice(pods, kubeThanos.PercentageToKill)

	logger.WithFields(log.Fields{
		"size": len(pods),
	}).Info("Pods to kill:")

	return pods, nil
}

func (kubeThanos *KubeThanos) SelectCandidatePods() ([]v1.Pod, error) {
	listOptions := metav1.ListOptions{LabelSelector: ""} // get all labels

	allPods, err := kubeThanos.Client.CoreV1().Pods(kubeThanos.Namespaces.String()).List(listOptions)
	if err != nil {
		return nil, err
	}

	filteredPods, err := filterByNamespaces(allPods.Items, kubeThanos.Namespaces)
	if err != nil {
		return nil, err
	}

	filteredPods = filterTerminatingPods(filteredPods)

	return filteredPods, nil
}

func (kubeThanos *KubeThanos) DeletePod(pod v1.Pod) error {
	logger.WithFields(log.Fields{
		"namespace": pod.Namespace,
		"name":      pod.Name,
	}).Info("terminating pod")

	if kubeThanos.DryRun {
		return nil
	}

	err := kubeThanos.Thanos.Kill(pod)
	if err != nil {
		return err
	}

	ref, err := reference.GetReference(scheme.Scheme, &pod)
	if err != nil {
		return err
	}

	kubeThanos.EventRecorder.Event(ref, v1.EventTypeNormal, "Killing", "Pod was killed by kubethanos to restore balance.")

	return nil
}

func filterByNamespaces(pods []v1.Pod, namespaces labels.Selector) ([]v1.Pod, error) {
	if namespaces.Empty() {
		return pods, nil
	}

	requirements, _ := namespaces.Requirements()
	var includeRequirements []labels.Requirement
	var excludeRequirements []labels.Requirement

	for _, req := range requirements {
		switch req.Operator() {
		case selection.Exists:
			includeRequirements = append(includeRequirements, req)
		case selection.DoesNotExist:
			excludeRequirements = append(excludeRequirements, req)
		default:
			return nil, fmt.Errorf("unsupported operator: %s", req.Operator())
		}
	}

	var filteredPods []v1.Pod

	for _, pod := range pods {
		included := len(includeRequirements) == 0

		selector := labels.Set{pod.Namespace: ""}

		for _, req := range includeRequirements {
			if req.Matches(selector) {
				included = true
				break
			}
		}

		for _, req := range excludeRequirements {
			if !req.Matches(selector) {
				included = false
				break
			}
		}

		if included {
			filteredPods = append(filteredPods, pod)
		}
	}

	return filteredPods, nil
}

func filterTerminatingPods(pods []v1.Pod) []v1.Pod {
	var filteredList []v1.Pod
	for _, pod := range pods {
		if pod.DeletionTimestamp != nil {
			continue
		}
		filteredList = append(filteredList, pod)
	}
	return filteredList
}

func RandomPodSlice(pods []v1.Pod, percentageToKill float64) []v1.Pod {
	count := int(float64(len(pods)) * percentageToKill)

	rand.Shuffle(len(pods), func(i, j int) { pods[i], pods[j] = pods[j], pods[i] })
	res := pods[0:count]
	return res
}
