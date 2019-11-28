package thanos

import (
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"

	log "github.com/sirupsen/logrus"
)

type Thanos struct {
	client      kubernetes.Interface
	logger      log.FieldLogger
}

func NewThanos(client kubernetes.Interface, logger log.FieldLogger) *Thanos {
	return &Thanos{
		client:      client,
		logger:      logger.WithField("thanos", "DeletePod"),
	}
}

func (t *Thanos) Kill(victimPod v1.Pod) error {
	t.logger.WithFields(log.Fields{
		"namespace": victimPod.Namespace,
		"name":      victimPod.Name,
	}).Debug("calling deletePod endpoint")

	return t.client.CoreV1().Pods(victimPod.Namespace).Delete(victimPod.Name,nil)
}
