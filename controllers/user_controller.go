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

package controllers

import (
	"context"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"

	monitoringv1beta1 "github.com/ArthurVardevanyan/openshift-monitoring-cr-controller/api/v1beta1"
)

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=monitoring.arthurvardevanyan.com,resources=users,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=monitoring.arthurvardevanyan.com,resources=users/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=monitoring.arthurvardevanyan.com,resources=users/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the User object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *UserReconciler) Reconcile(reconcilerContext context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(reconcilerContext)
	log.V(1).Info(req.Name)

	// Incept Object
	var monitoring monitoringv1beta1.User
	if err := r.Get(reconcilerContext, req.NamespacedName, &monitoring); err != nil {
		if strings.Contains(err.Error(), "not found") {
			log.V(1).Info("Monitoring Config Not Found or No Longer Exists!")
			return ctrl.Result{}, nil
		} else {
			log.Error(err, "Unable to fetch Monitoring Object")
			return ctrl.Result{}, client.IgnoreNotFound(err)
		}
	}

	const (
		finalizer     string = "user.monitoring.arthurvardevanyan.com/finalizer"
		namespace     string = "openshift-user-workload-monitoring"
		configMapName string = "user-workload-monitoring-config"
	)

	configMapData := make(map[string]string)
	MonitoringYaml, err := yaml.Marshal(&monitoring.Spec)
	if err != nil {
		log.Error(err, "Unable to Marshal ConfigMap Struct to Yaml!")
	}
	configMapData["config.yaml"] = string(MonitoringYaml)

	configMap := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: namespace,
		},
		Data: configMapData,
	}

	// Delete Logic, Create / Update Finalizers, and Garbage Collect LogSink on Object Deletion
	// https://book.kubebuilder.io/reference/using-finalizers.html
	if monitoring.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// registering our finalizer.
		if !controllerutil.ContainsFinalizer(&monitoring, finalizer) {
			// https://sdk.operatorframework.io/docs/building-operators/golang/references/client/#patch
			patch := client.MergeFrom(monitoring.DeepCopy())
			controllerutil.AddFinalizer(&monitoring, finalizer)
			// https://sdk.operatorframework.io/docs/building-operators/golang/references/client/#patch
			if err := r.Patch(reconcilerContext, &monitoring, patch); err != nil {
				return ctrl.Result{}, err
			} else {
				return ctrl.Result{}, nil
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(&monitoring, finalizer) {
			// our finalizer is present, so lets handle any external dependency
			log.V(1).Info("Deleting ConfigMap!")
			err := r.Delete(reconcilerContext, &configMap)
			if err != nil {
				log.Error(err, "Unable to Delete ConfigMap!")
			}

			// remove our finalizer from the list and update it.
			patch := client.MergeFrom(monitoring.DeepCopy())
			controllerutil.RemoveFinalizer(&monitoring, finalizer)
			if err := r.Patch(reconcilerContext, &monitoring, patch); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

	// https://medium.com/@aneeshputtur/kubernetes-operators-with-external-configmap-b972c9c36bbe
	err = r.Get(reconcilerContext, types.NamespacedName{Name: configMapName, Namespace: namespace}, &corev1.ConfigMap{})
	if err != nil {
		log.V(1).Info("Create ConfigMap")
		r.Create(reconcilerContext, &configMap)
	} else {
		log.V(1).Info("Update ConfigMap")
		r.Update(reconcilerContext, &configMap)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1beta1.User{}).
		Complete(r)
}
