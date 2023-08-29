/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	// "reflect"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "volcano.sh/apis/pkg/apis/batch/v1alpha1"

	slaoperatorv1alpha1 "github.com/binhfdv/sla-operator/api/v1alpha1"
	"github.com/binhfdv/sla-operator/pkg/resources"
	// "github.com/binhfdv/sla-operator/pkg/rest"
)

// SlamlReconciler reconciles a Slaml object
type SlamlReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=sla-operator.dcn.com,resources=slamls,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sla-operator.dcn.com,resources=slamls/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sla-operator.dcn.com,resources=slamls/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Slaml object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile

// Reference: https://book.kubebuilder.io/cronjob-tutorial/controller-implementation

func (r *SlamlReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)

	// TODO(user): your logic here
	var clientResource = &slaoperatorv1alpha1.Slaml{}

	// 1: Load the slaml jobs by name
	if err := r.Get(ctx, req.NamespacedName, clientResource); err != nil {
		log.Error(err, "unable to fetch client")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// 2: List all active jobs, and update the status
	var childClientResource v1alpha1.JobList
	if err := r.List(ctx, &childClientResource, client.InNamespace(req.Namespace), client.MatchingFields{jobOwnerKey: req.Name}); err != nil {
		log.Error(err, "unable to list child Jobs")
		return ctrl.Result{}, err
	}

	// find the active list of jobs
	var pendingJobs []*v1alpha1.Job
	var runningJobs []*v1alpha1.Job
	var successfulJobs []*v1alpha1.Job
	var failedJobs []*v1alpha1.Job
	// var mostRecentTime *time.Time // find the last run so we can update the status
	// isJobFinished
	isJobFinished := func(job *v1alpha1.Job) (bool, v1alpha1.JobPhase) {
		for _, c := range job.Status.Conditions {
			if c.Status == v1alpha1.Completed {
				return true, c.Status
			} else if c.Status == v1alpha1.Pending || c.Status == v1alpha1.Running {
				return false, c.Status
			}
		}

		return false, ""
	}
	// getScheduledTimeForJob

	for i, job := range childClientResource.Items {
		_, finishedType := isJobFinished(&job)
		switch finishedType {
		case "": // ongoing
			failedJobs = append(failedJobs, &childClientResource.Items[i])
		case v1alpha1.Pending:
			pendingJobs = append(pendingJobs, &childClientResource.Items[i])
		case v1alpha1.Running:
			runningJobs = append(runningJobs, &childClientResource.Items[i])
		case v1alpha1.Completed:
			successfulJobs = append(successfulJobs, &childClientResource.Items[i])
		}

		// We'll store the launch time in an annotation, so we'll reconstitute that from
		// the active jobs themselves.
		// scheduledTimeForJob, err := getScheduledTimeForJob(&job)
		// if err != nil {
		// 	log.Error(err, "unable to parse schedule time for child job", "job", &job)
		// 	continue
		// }
		// if scheduledTimeForJob != nil {
		// 	if mostRecentTime == nil {
		// 		mostRecentTime = scheduledTimeForJob
		// 	} else if mostRecentTime.Before(*scheduledTimeForJob) {
		// 		mostRecentTime = scheduledTimeForJob
		// 	}
		// }
	}

	log.V(1).Info("job count", "running jobs", len(runningJobs), "pending jobs", len(pendingJobs), "successful jobs", len(successfulJobs), "failed jobs", len(failedJobs))

	// if mostRecentTime != nil {
	// 	cronJob.Status.LastScheduleTime = &metav1.Time{Time: *mostRecentTime}
	// } else {
	// 	cronJob.Status.LastScheduleTime = nil
	// }
	// cronJob.Status.Active = nil
	// for _, activeJob := range activeJobs {
	// 	jobRef, err := ref.GetReference(r.Scheme, activeJob)
	// 	if err != nil {
	// 		log.Error(err, "unable to make reference to active job", "job", activeJob)
	// 		continue
	// 	}
	// 	cronJob.Status.Active = append(cronJob.Status.Active, *jobRef)
	// }

	// clientResourceOld := clientResource.DeepCopy()

	if clientResource.Status.ClientStatus == "" {
		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusPending
	}

	switch clientResource.Status.ClientStatus {
	case slaoperatorv1alpha1.StatusPending:
		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusRunning

		err := r.Status().Update(context.TODO(), clientResource)
		if err != nil {
			log.Error(err, "failed to update client status")
			return ctrl.Result{}, err
		} else {
			log.Info("updated client status: " + clientResource.Status.ClientStatus)
			return ctrl.Result{Requeue: true}, nil
		}

	case slaoperatorv1alpha1.StatusRunning:
		pod := resources.CreateJobPod(clientResource)
		startPoint := time.Now()
		query := &corev1.Pod{}
		// fmt.Println("pod: \n\n", pod)
		log.Info("HERE 1\n")
		err := r.Client.Get(ctx, client.ObjectKey{Namespace: pod.Namespace, Name: pod.ObjectMeta.Name}, query)
		if err != nil && errors.IsNotFound(err) {
			if clientResource.Status.LastPodName == "" {
				err = ctrl.SetControllerReference(clientResource, pod, r.Scheme)
				if err != nil {
					log.Info("ERROR 1")
					return ctrl.Result{}, err
				}

				err = r.Create(context.TODO(), pod)
				if err != nil {
					log.Info("ERROR 2")
					return ctrl.Result{}, err
				}

				log.Info("pod created successfully", "name", pod.Name)
				log.Info("response time: ", time.Since(startPoint))
				return ctrl.Result{}, nil
			} else {
				clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusCleaning
			}
		} else if err != nil {
			log.Error(err, "cannot get pod")
			return ctrl.Result{}, err
		} else if query.Status.Phase == corev1.PodFailed || query.Status.Phase == corev1.PodSucceeded {
			log.Info("container terminated", "reason", query.Status.Reason, "message", query.Status.Message)

			clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusCleaning
		} else if query.Status.Phase == corev1.PodRunning {
			// 		if clientResource.Status.LastPodName != clientResource.Spec.ContainerImage+clientResource.Spec.ContainerTag {
			// 			if query.Status.ContainerStatuses[0].Ready {
			// 				log.Info("Trying to bind to: " + query.Status.PodIP)

			// 				if !rest.GetClient(clientResource, query.Status.PodIP) {
			// 					if rest.BindClient(clientResource, query.Status.PodIP) {
			// 						log.Info("Client" + clientResource.Spec.ClientId + " is binded to pod " + query.ObjectMeta.GetName() + ".")
			// 						clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusCleaning
			// 					} else {
			// 						log.Info("Client not added.")
			// 					}
			// 				} else {
			// 					log.Info("Client binded already.")
			// 				}
			// 			} else {
			// 				log.Info("Container not ready, reschedule bind")
			// 				return ctrl.Result{Requeue: true}, err
			// 			}

			// 			log.Info("Client last pod name: " + clientResource.Status.LastPodName)
			// 			log.Info("Pod is running.")
			// 		}
		} else if query.Status.Phase == corev1.PodPending {
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{Requeue: true}, err
		}

	// 	if !reflect.DeepEqual(clientResourceOld.Status, clientResource.Status) {
	// 		err = r.Status().Update(context.TODO(), clientResource)
	// 		if err != nil {
	// 			log.Error(err, "failed to update client status from running")
	// 			return ctrl.Result{}, err
	// 		} else {
	// 			log.Info("updated client status RUNNING -> " + clientResource.Status.ClientStatus)
	// 			return ctrl.Result{Requeue: true}, nil
	// 		}
	// 	}
	// case slaoperatorv1alpha1.StatusCleaning:
	// 	query := &corev1.Pod{}
	// 	HasClients := rest.HasClients(clientResource, query.Status.PodIP)

	// 	err := r.Client.Get(ctx, client.ObjectKey{Namespace: clientResource.Namespace, Name: clientResource.Status.LastPodName}, query)
	// 	if err == nil && clientResource.ObjectMeta.DeletionTimestamp.IsZero() {
	// 		if !HasClients {
	// 			err = r.Delete(context.TODO(), query)
	// 			if err != nil {
	// 				log.Error(err, "Failed to remove old pod")
	// 				return ctrl.Result{}, err
	// 			} else {
	// 				log.Info("Old pod removed")
	// 				return ctrl.Result{Requeue: true}, nil
	// 			}
	// 		}
	// 	}

	// 	if clientResource.Status.LastPodName != clientResource.Spec.ContainerImage+clientResource.Spec.ContainerTag {
	// 		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusRunning
	// 		clientResource.Status.LastPodName = clientResource.Spec.ContainerImage + clientResource.Spec.ContainerTag
	// 	} else {
	// 		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusPending
	// 		clientResource.Status.LastPodName = ""
	// 	}

	// 	if !reflect.DeepEqual(clientResourceOld.Status, clientResource.Status) {
	// 		err = r.Status().Update(context.TODO(), clientResource)
	// 		if err != nil {
	// 			log.Error(err, "failed to update client status from cleaning")
	// 			return ctrl.Result{}, err
	// 		} else {
	// 			log.Info("updated client status CLEANING -> " + clientResource.Status.ClientStatus)
	// 			return ctrl.Result{Requeue: true}, nil
	// 		}
	// 	}
	default:
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SlamlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slaoperatorv1alpha1.Slaml{}).
		Complete(r)
}
