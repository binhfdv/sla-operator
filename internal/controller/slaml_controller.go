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
	logx "github.com/theritikchoure/logx"
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
			logx.LogWithLevel("updated client status to Running", "SUCCESS")
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
			logx.LogWithLevel("job scheduled", "SUCCESS")

		} else if query.Status.Phase == corev1.PodPending {
			return ctrl.Result{Requeue: true}, nil
		} else {
			return ctrl.Result{Requeue: true}, err
		}
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
