/*
Copyright 2023 API Testing Authors.

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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1alpha1 "github.com/linuxsuren/api-testing/operator/api/v1alpha1"
)

// ATestReconciler reconciles a ATest object
type ATestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=core.linuxsuren.github.com,resources=atests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core.linuxsuren.github.com,resources=atests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=core.linuxsuren.github.com,resources=atests/finalizers,verbs=update

//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=persistentvolumeclaims,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ATest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ATestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (result ctrl.Result, err error) {
	logger := log.FromContext(ctx)
	atest := &corev1alpha1.ATest{}
	if err = r.Get(ctx, req.NamespacedName, atest); err != nil {
		err = client.IgnoreAlreadyExists(err)
		return
	}

	logger.Info("Reconciling ATest", "atest", atest.Name)
	pvc := newPVC(atest)
	if pvc != nil {
		existingPVC := &corev1.PersistentVolumeClaim{}
		if err = r.Client.Get(ctx, req.NamespacedName, existingPVC); err != nil {
			if !errors.IsNotFound(err) {
				return
			}

			if err = r.Client.Create(ctx, pvc); err != nil {
				return
			}
		} else {
			existingPVC.Spec = pvc.Spec
			err = r.Client.Update(ctx, existingPVC)
		}
	}

	configMap := newConfigMap(atest)
	existingConfigMap := &corev1.ConfigMap{}
	if err = r.Client.Get(ctx, req.NamespacedName, existingConfigMap); err != nil {
		if !errors.IsNotFound(err) {
			return
		}

		if err = r.Client.Create(ctx, configMap); err != nil {
			return
		}
	} else {
		existingConfigMap.Data = configMap.Data
		err = r.Client.Update(ctx, existingConfigMap)
	}

	deploy := newDeployment(atest)
	existingDeploy := &appsv1.Deployment{}
	if err = r.Client.Get(ctx, req.NamespacedName, existingDeploy); err != nil {
		if !errors.IsNotFound(err) {
			return
		}

		if err = r.Client.Create(ctx, deploy); err != nil {
			return
		}
	} else {
		existingDeploy.Spec = deploy.Spec
		err = r.Client.Update(ctx, existingDeploy)
	}

	svc := newService(atest)
	existingSvc := &corev1.Service{}
	if err = r.Client.Get(ctx, req.NamespacedName, existingSvc); err != nil {
		if !errors.IsNotFound(err) {
			return
		}

		if err = r.Client.Create(ctx, svc); err != nil {
			return
		}
	} else {
		existingSvc.Spec.Ports = svc.Spec.Ports
		existingSvc.Spec.Type = svc.Spec.Type
		existingSvc.Spec.Selector = svc.Spec.Selector
		err = r.Client.Update(ctx, existingSvc)
	}

	if err == nil {
		result = ctrl.Result{RequeueAfter: 30 * time.Second}
	}
	return
}

// SetupWithManager sets up the controller with the Manager.
func (r *ATestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1alpha1.ATest{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
