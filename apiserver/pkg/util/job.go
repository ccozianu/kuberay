package util

import (
	api "github.com/ray-project/kuberay/proto/go_client"
	rayv1api "github.com/ray-project/kuberay/ray-operator/apis/ray/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RayJob struct {
	*rayv1api.RayJob
}

const rayJobDefaultVersion = "1.13"

// NewRayJob creates a RayJob.
func NewRayJob(apiJob *api.RayJob, computeTemplateMap map[string]*api.ComputeTemplate) (*RayJob, error) {
	rayJob := &rayv1api.RayJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        apiJob.Name,
			Namespace:   apiJob.Namespace,
			Labels:      buildRayJobLabels(apiJob),
			Annotations: buildRayJobAnnotations(apiJob),
		},
		Spec: rayv1api.RayJobSpec{
			Entrypoint:               apiJob.Entrypoint,
			Metadata:                 apiJob.Metadata,
			RuntimeEnvYAML:           apiJob.RuntimeEnv,
			ShutdownAfterJobFinishes: apiJob.ShutdownAfterJobFinishes,
			TTLSecondsAfterFinished:  &apiJob.TtlSecondsAfterFinished,
			JobId:                    apiJob.JobId,
			RayClusterSpec:           nil,
			ClusterSelector:          apiJob.ClusterSelector,
		},
	}
	if apiJob.ClusterSpec != nil {
		clusterSpec, err := buildRayClusterSpec(rayJobDefaultVersion, nil, apiJob.ClusterSpec, computeTemplateMap)
		if err != nil {
			return nil, err
		}
		rayJob.Spec.RayClusterSpec = clusterSpec
	}
	return &RayJob{rayJob}, nil
}

func (j *RayJob) Get() *rayv1api.RayJob {
	return j.RayJob
}

func buildRayJobLabels(job *api.RayJob) map[string]string {
	labels := map[string]string{}
	labels[RayClusterNameLabelKey] = job.Name
	labels[RayClusterUserLabelKey] = job.User
	labels[KubernetesApplicationNameLabelKey] = ApplicationName
	labels[KubernetesManagedByLabelKey] = ComponentName
	return labels
}

func buildRayJobAnnotations(job *api.RayJob) map[string]string {
	return job.Metadata
}
