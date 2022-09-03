/*
Copyright 2022.

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

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// Our annotation
	doubleDownAnnotation = "jdocklabs.co.uk/double-down"

	// Already doubled, no change required
	isDoubledAnnotation = "jdocklabs.co.uk/doubled"
)

// DeploymentReconciler reconciles a Deployment object
type DeploymentReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx) // Logs in K-V pairs after the message

	var d appsv1.Deployment

	if err := r.Get(ctx, req.NamespacedName, &d); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		l.Error(err, "unable to fetch Deployment", "Deployment", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if d.Annotations[doubleDownAnnotation] == "true" {
		if d.Annotations[isDoubledAnnotation] == "true" {
			l.Info("nothing to do", "Deployment", req.NamespacedName)
			// Nothing to do, avoids pods ^ 2 constant growth
			return ctrl.Result{}, nil
		}

		var double int32 = 2
		currentReplicas := *d.Spec.Replicas
		doubleReplicas := currentReplicas * double
		d.Spec.Replicas = &doubleReplicas

		m := make(map[string]string)
		m[isDoubledAnnotation] = "true"
		d.SetAnnotations(m)

		if err := r.Update(ctx, &d); err != nil {
			if apierrors.IsConflict(err) {
				// The Deployment has been updated since we read it.
				// Requeue the Deployment to try to reconciliate again.
				return ctrl.Result{Requeue: true}, nil
			}
			if apierrors.IsNotFound(err) {
				// The Deployment has been deleted since we read it.
				// Requeue the Deployment to try to reconciliate again.
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
		l.Info("doubled down", "Deployment", req.NamespacedName)
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
