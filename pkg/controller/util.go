package controller

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const PracticeGroupConfigMapName = "practice-groups"
const IndexFile = "index.html"

func DoNotRequeue() (reconcile.Result, error) {
	return reconcile.Result{}, nil
}

func RequeueWithError(err error) (reconcile.Result, error) {
	return reconcile.Result{}, err
}

// Bad practice, but nice for demo
func (r *PracticeGroupReconciler) EmitWarningEvent(object runtime.Object, reason string, message string) {
	r.Recorder.Event(
		object,
		corev1.EventTypeWarning, reason,
		message,
	)
}

func (r *PracticeGroupReconciler) EmitNormalEvent(object runtime.Object, reason string, message string) {
	r.Recorder.Event(
		object,
		corev1.EventTypeNormal, reason,
		message,
	)
}

func configMapFor() {

}
