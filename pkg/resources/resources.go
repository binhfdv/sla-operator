package resources

import (
	// "strings"

	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "volcano.sh/apis/pkg/apis/batch/v1alpha1"
	v1alpha1event "volcano.sh/apis/pkg/apis/bus/v1alpha1"

	// v1beta1 "volcano.sh/apis/pkg/apis/scheduling/v1beta1"

	slaoperatorv1alpha1 "github.com/binhfdv/sla-operator/api/v1alpha1"
)

func Age(pod *corev1.Pod) uint {
	diff := uint(time.Now().Sub(pod.Status.StartTime.Time).Seconds())
	return diff //uint(diff / 24)
}

func getLabels(clientResource *slaoperatorv1alpha1.Slaml) map[string]string {
	return map[string]string{
		"app": clientResource.Spec.Name,
	}
}

func getAnnotations(clientResource *slaoperatorv1alpha1.Slaml) map[string]string {
	return map[string]string{
		"sla-waiting-time": "1m",
	}
}

func getPlugins(clientResource *slaoperatorv1alpha1.Slaml) map[string][]string {
	return map[string][]string{
		"pytorch": {"--master=master", "--worker=worker", "--port=23456"},
		// "gang": {"1"},
	}
}

func getTask(clientResource *slaoperatorv1alpha1.Slaml) []v1alpha1.TaskSpec {
	var tasks []v1alpha1.TaskSpec
	for _, task := range clientResource.Spec.Tasks {
		temp := v1alpha1.TaskSpec{
			Name:     task.TaskName,
			Replicas: task.ContainerReplicas,
			// Policies: []v1alpha1.LifecyclePolicy{
			// 	{
			// 		Event:  v1alpha1event.TaskCompletedEvent,
			// 		Action: v1alpha1event.CompleteJobAction,
			// 	},
			// },
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            task.TaskName,
							Image:           task.ContainerRegistry + "/" + task.ContainerImage + ":" + task.ContainerTag,
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
					RestartPolicy: corev1.RestartPolicyOnFailure,
				},
			},
		}
		tasks = append(tasks, temp)
	}
	return tasks
}

// func CreateQueue(clientResource *slaoperatorv1alpha1.Slaml) *v1beta1.Queue {
// 	return &v1beta1.Queue{

// 	}
// }

// func CreatePodGroup(clientResource *slaoperatorv1alpha1.Slaml) *v1beta1.PodGroup {
// 	return &v1beta1.PodGroup{

// 	}
// }

func CreateJobPod(clientResource *slaoperatorv1alpha1.Slaml) *v1alpha1.Job {
	// fmt.Println(getTask(clientResource))
	return &v1alpha1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:        clientResource.Spec.Name,
			Namespace:   clientResource.Namespace,
			Labels:      getLabels(clientResource),
			Annotations: getAnnotations(clientResource),
		},
		Spec: v1alpha1.JobSpec{
			SchedulerName: "volcano",
			MinAvailable:  1,
			Policies: []v1alpha1.LifecyclePolicy{
				{
					Event:  v1alpha1event.TaskCompletedEvent,
					Action: v1alpha1event.CompleteJobAction,
				},
			},
			Plugins: getPlugins(clientResource),
			Tasks:   getTask(clientResource),
		},
	}
}

// func EnforceDeadline(clientResource *slaoperatorv1alpha1.Slaml) *v1alpha1.Job {

// 	return &v1alpha1.Job{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:        clientResource.Spec.Name,
// 			Namespace:   clientResource.Namespace,
// 			Labels:      getLabels(clientResource),
// 			Annotations: getAnnotations(clientResource),
// 		},
// 		Spec: v1alpha1.JobSpec{
// 			SchedulerName: "volcano",
// 			MinAvailable:  1,
// 			Plugins:       getPlugins(clientResource),
// 			Tasks:         getTask(clientResource),
// 		},
// 	}
// }

// if  strings.ToLower(clientResource.Spec.VolcanoKind) == "job" {
// 	return betav1.Job{

// 	}
// } else if  strings.ToLower(clientResource.Spec.VolcanoKind) == "podgroup" {
// 	fmt.Println(num, "has 1 digit")
// } else {
// 	fmt.Println(num, "has multiple digits")
// }
