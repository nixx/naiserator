package horizontalpodautoscaler

import (
	"github.com/nais/liberator/pkg/apis/nais.io/v1"
	autoscaling "k8s.io/api/autoscaling/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/nais/naiserator/pkg/resourcecreator/resource"
	"github.com/nais/naiserator/pkg/util"
)

type Source interface {
	resource.Source
	GetReplicas() *nais_io_v1.Replicas
}

func Create(source Source, ast *resource.Ast) {
	replicas := source.GetReplicas()

	if (*replicas.Max) <= 0 {
		return
	}
	if replicas.DisableAutoScaling || *replicas.Min == *replicas.Max {
		return
	}

	hpa := &autoscaling.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: autoscaling.SchemeGroupVersion.Identifier(),
		},
		ObjectMeta: resource.CreateObjectMeta(source),
		Spec: autoscaling.HorizontalPodAutoscalerSpec{
			ScaleTargetRef: autoscaling.CrossVersionObjectReference{
				APIVersion: "apps/v1",
				Kind:       "Deployment",
				Name:       source.GetName(),
			},
			Metrics: []autoscaling.MetricSpec{
				{
					Type: autoscaling.ResourceMetricSourceType,
					Resource: &autoscaling.ResourceMetricSource{
						Name: "cpu",
						Target: autoscaling.MetricTarget{
							Type:               autoscaling.UtilizationMetricType,
							AverageUtilization: util.Int32p(int32(replicas.CpuThresholdPercentage)),
						},
					},
				},
			},
			MinReplicas: util.Int32p(int32(*replicas.Min)),
			MaxReplicas: int32(*replicas.Max),
		},
	}
	ast.AppendOperation(resource.OperationCreateOrUpdate, hpa)
}
