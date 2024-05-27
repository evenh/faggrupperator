/*
Copyright 2024.

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
	"bytes"
	"context"
	"fmt"
	bekknov1alpha1 "github.com/evenh/faggrupperator/api/v1alpha1"
	"github.com/evenh/faggrupperator/assets"
	"github.com/gosimple/slug"
	"html/template"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// PracticeGroupReconciler reconciles a PracticeGroup object
type PracticeGroupReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

var t = template.New("master")

func init() {
	// Read templates from binary
	_, err := t.ParseFS(assets.BundledAssets, "*.html")
	if err != nil {
		panic(err)
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PracticeGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&bekknov1alpha1.PracticeGroup{}).
		Owns(&corev1.ConfigMap{}).
		Watches(
			&corev1.ConfigMap{},
			&handler.EnqueueRequestForObject{},
			builder.WithPredicates(predicate.Funcs{
				UpdateFunc: func(e event.UpdateEvent) bool {
					return e.ObjectNew.GetName() == PracticeGroupConfigMapName
				},
				CreateFunc: func(e event.CreateEvent) bool {
					return e.Object.GetName() == PracticeGroupConfigMapName
				},
				DeleteFunc: func(e event.DeleteEvent) bool {
					return e.Object.GetName() == PracticeGroupConfigMapName
				},
			}),
		).
		Complete(r)
}

// Define what RBAC permissions we need
//
//+kubebuilder:rbac:groups=bekk.no,resources=practicegroups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=bekk.no,resources=practicegroups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=bekk.no,resources=practicegroups/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *PracticeGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Check if the request is for the PracticeGroup or the ConfigMap
	if req.Name == PracticeGroupConfigMapName {
		return RequeueWithError(r.reconcileConfigMap(ctx, req.Namespace))
	}

	// Look up existing object
	pg := &bekknov1alpha1.PracticeGroup{}
	err := r.Client.Get(ctx, req.NamespacedName, pg)

	if client.IgnoreNotFound(err) != nil {
		// Something other than 'not found'
		logger.Error(err, "unable to fetch PracticeGroup")
		r.EmitWarningEvent(pg, "ReconcileStartFail", "error while fetching the PracticeGroup, it might have been deleted")
		return RequeueWithError(err)
	}

	err = r.reconcileConfigMap(ctx, pg.Namespace)
	if client.IgnoreNotFound(err) != nil {
		logger.Error(err, "unable to update ConfigMap with group "+pg.Name)
		r.EmitWarningEvent(pg, "ReconcileConfigMapFailed", "could not update ConfigMap with PracticeGroup")
		// Already retried
		return DoNotRequeue()
	}

	r.EmitNormalEvent(pg, "UpdatedConfigMap", "successfully updated ConfigMap with new HTML")

	// All good
	return DoNotRequeue()
}

func (r *PracticeGroupReconciler) reconcileConfigMap(ctx context.Context, namespace string) error { //nolint:lll
	// Fetch all PracticeGroup instances
	practiceGroupList := &bekknov1alpha1.PracticeGroupList{}
	opts := []client.ListOption{client.InNamespace(namespace)}
	if err := r.List(ctx, practiceGroupList, opts...); err != nil {
		return err
	}

	// Render into a map of HTML files
	renderedData, err := mapPracticeGroupsToHTML(practiceGroupList.Items)
	if err != nil {
		return fmt.Errorf("could not render HTML files or index.html: %w", err)
	}

	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{Name: PracticeGroupConfigMapName, Namespace: namespace}

	// Get the existing ConfigMap if it exists
	err = r.Client.Get(ctx, cmKey, cm)
	if client.IgnoreNotFound(err) != nil {
		return fmt.Errorf("failed to get ConfigMap: %v", err)
	}

	if errors.IsNotFound(err) {
		// ConfigMap doesn't exist, create a new one
		cm = &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      PracticeGroupConfigMapName,
				Namespace: namespace,
			},
			Data: renderedData,
		}
		if err := r.Client.Create(ctx, cm); err != nil {
			return fmt.Errorf("failed to create ConfigMap: %v", err)
		}
		return nil
	}

	// Update the existing ConfigMap
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Get the latest version of the ConfigMap
		if err := r.Client.Get(ctx, cmKey, cm); err != nil {
			return fmt.Errorf("failed to get ConfigMap: %v", err)
		}

		// Update the entry in the data map
		cm.Data = renderedData

		// Update the ConfigMap
		if err := r.Client.Update(ctx, cm); err != nil {
			return fmt.Errorf("failed to update ConfigMap: %v", err)
		}
		return nil
	})
}

func mapPracticeGroupsToHTML(groups []bekknov1alpha1.PracticeGroup) (map[string]string, error) {
	result := make(map[string]string)
	indexList := make(map[string]string)
	if len(groups) == 0 {
		return result, nil
	}

	for _, group := range groups {
		fileName := fmt.Sprintf("%s.html", slug.Make(group.Name))
		content, err := renderGroup(&group)
		if err != nil {
			return nil, err
		}

		result[fileName] = content
		indexList[fileName] = group.Spec.Name
	}

	// Create index.html
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "index.html", indexList)
	if err != nil {
		return result, err
	}

	result[IndexFile] = buf.String()

	return result, nil
}

func renderGroup(group *bekknov1alpha1.PracticeGroup) (string, error) {
	buf := new(bytes.Buffer)
	err := t.ExecuteTemplate(buf, "pg.html", group)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
