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
	"encoding/json"
	"time"

	// "reflect"

	"fmt"

	logx "github.com/theritikchoure/logx"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1alpha1 "volcano.sh/apis/pkg/apis/batch/v1alpha1"

	slaoperatorv1alpha1 "github.com/binhfdv/sla-operator/api/v1alpha1"
	"github.com/binhfdv/sla-operator/pkg/resources"
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
	logx.ColoringEnabled = true
	// TODO(user): your logic here
	var clientResource = &slaoperatorv1alpha1.Slaml{}

	if err := r.Get(ctx, req.NamespacedName, clientResource); err != nil {
		logx.LogWithLevel("unable to fetch client", "ERROR")
		log.Error(err, "unable to fetch client")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if clientResource.Status.ClientStatus == "" {
		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusPending
		logx.LogWithLevel("update status to Pending", "INFO")
		log.Info("")
	}

	switch clientResource.Status.ClientStatus {
	case slaoperatorv1alpha1.StatusPending:
		clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusRunning

		err := r.Status().Update(context.TODO(), clientResource)
		if err != nil {
			logx.LogWithLevel("failed to update client status", "ERROR")
			log.Error(err, "failed to update client status")
			return ctrl.Result{}, err
		} else {
			logx.LogWithLevel("updated client status to RUNNING", "SUCCESS")
			log.Info("updated client status: " + clientResource.Status.ClientStatus)
			return ctrl.Result{Requeue: true}, nil
		}

	case slaoperatorv1alpha1.StatusRunning:
		job := resources.CreateJobPod(clientResource)
		jsonData, err := json.MarshalIndent(job.Spec, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return ctrl.Result{}, err
		}

		fmt.Printf("Job Spec: %s\n", string(jsonData))

		// startPoint := time.Now()
		query := &corev1.Pod{}
		fmt.Println("Jod: \n\n", job)
		logx.LogWithLevel("at RUNNING stage", "INFO")

		err = r.Client.Get(ctx, client.ObjectKey{Namespace: job.Namespace, Name: job.ObjectMeta.Name}, query)
		if err != nil && errors.IsNotFound(err) {
			err = ctrl.SetControllerReference(clientResource, job, r.Scheme)
			if err != nil {
				logx.LogWithLevel("cannot SetControllerReference for the job", "ERROR")
				return ctrl.Result{}, err
			}

			err = r.Create(context.TODO(), job)
			if err != nil {
				logx.LogWithLevel("cannot create the job", "ERROR")
				return ctrl.Result{}, err
			}

			log.Info("pod created successfully", "name", job.Name)

			return ctrl.Result{}, nil
		} else {
			clientResource.Status.ClientStatus = slaoperatorv1alpha1.StatusCleaning
		}

		namespace := "machine"
		// matchLabels := resources.GetLabels(clientResource)
		// fmt.Printf("Matching labels: %s\n", matchLabels)
		// list all pods in the namespace
		var childPods corev1.PodList
		if err := r.List(ctx, &childPods, client.InNamespace(req.Namespace)); err != nil {
			log.Error(err, "unable to list child pods")
			return ctrl.Result{}, err
		}

		// print status for each pod
		fmt.Printf("Pod Status for namespace %s:\n", namespace)
		for _, pod := range childPods.Items {
			// fmt.Printf("Pod Name: %s\n", pod.Name)
			// fmt.Printf("Phase: %s\n", pod.Status.Phase)
			// fmt.Printf("Conditions: %v\n", pod.Status.Conditions)
			// fmt.Printf("Container Statuses: %v\n", pod.Status.ContainerStatuses)
			// jsonData, err := json.MarshalIndent(pod.Spec, "", "  ")
			// if err != nil {
			// 	fmt.Println("Error marshaling JSON:", err)
			// 	return ctrl.Result{}, err
			// }

			fmt.Printf("Pod Spec: %s\n", pod.Spec.Subdomain)
			age := resources.Age(&pod)
			name := pod.Name
			fmt.Printf("Pod Spec: %d\n", age)

			if err := r.deletePod(ctx, &pod, age); err != nil {
				log.Error(err, "unable to delete", "source pod", pod)
				return ctrl.Result{}, err
			}
			logx.LogWithLevel("pod "+name+" deleted", "INFO")
			fmt.Println("------------------------------")
		}
		fmt.Println("---------------------------------------e---------------------------------------------------")
		// r.updateStatus(clientResource, slaoperatorv1alpha1.StatusCleaning)

	case slaoperatorv1alpha1.StatusCleaning:
		time.Sleep(8 * time.Second)
		logx.LogWithLevel("at CLEANING stage", "INFO")
		// r.updateStatus(clientResource, slaoperatorv1alpha1.StatusPending)

	default:
	}

	return ctrl.Result{}, nil
}

var (
	jobOwnerKey = ".metadata.controller"
	apiGVStr    = v1alpha1.SchemeGroupVersion
)

// SetupWithManager sets up the controller with the Manager.
func (r *SlamlReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slaoperatorv1alpha1.Slaml{}).
		Complete(r)
}

func (r *SlamlReconciler) deletePod(ctx context.Context, pod *corev1.Pod, age uint) error {
	if age > 50 {
		if err := r.Delete(ctx, pod); err != nil {
			return err
		}
	}
	return nil
}

func (r *SlamlReconciler) updateStatus(client *slaoperatorv1alpha1.Slaml, status string) (ctrl.Result, error) {
	client.Status.ClientStatus = status

	err := r.Status().Update(context.TODO(), client)
	if err != nil {
		logx.LogWithLevel("failed to update client status to "+status, "ERROR")
		return ctrl.Result{}, err
	} else {
		logx.LogWithLevel("updated client status to "+client.Status.ClientStatus, "SUCCESS")
		return ctrl.Result{Requeue: true}, nil
	}
}

func (r *SlamlReconciler) assignResources(client *slaoperatorv1alpha1.Slaml) error {



	return
}
