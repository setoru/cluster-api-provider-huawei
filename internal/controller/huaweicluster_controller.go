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
	"context"

	"github.com/pkg/errors"
	hwclient "github.com/setoru/cluster-api-provider-huawei/internal/cloud/client"
	"github.com/setoru/cluster-api-provider-huawei/internal/cloud/scope"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/setoru/cluster-api-provider-huawei/api/v1alpha1"
)

// HuaweiClusterReconciler reconciles a HuaweiCluster object
type HuaweiClusterReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweiclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweiclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweiclusters/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=clusters;clusters/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HuaweiCluster object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *HuaweiClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	logger := log.FromContext(ctx)

	huaweiCluster := &infrav1.HuaweiCluster{}
	if err := r.Get(ctx, req.NamespacedName, huaweiCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	cluster, err := util.GetOwnerCluster(ctx, r.Client, huaweiCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster != nil {
		logger.Info("Waiting for Cluster Controller to set OwnerRef on HuaweiCluster")
		return ctrl.Result{}, err
	}

	logger = logger.WithValues("cluster", klog.KObj(cluster))
	ctx = ctrl.LoggerInto(ctx, logger)

	if annotations.IsPaused(cluster, huaweiCluster) {
		logger.Info("HuaweiCluster or owning Cluster is marked as paused, not reconciling")

		return ctrl.Result{}, nil
	}

	secretName := huaweiCluster.Spec.CredentialsSecret.Name
	region := huaweiCluster.Spec.Region
	hwClient, err := hwclient.NewClient(r.Client, secretName, huaweiCluster.Namespace, region)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Create the scope.
	clusterScope := &scope.ClusterScope{
		Client:        r.Client,
		Logger:        logger,
		Cluster:       cluster,
		HuaweiCluster: huaweiCluster,
		HuaweiClient:  hwClient,
	}
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(huaweiCluster, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Always attempt to Patch the HuaweiCluster object and status after each reconciliation.
	defer func() {
		if err := patchHelper.Patch(ctx, huaweiCluster); err != nil {
			logger.Error(err, "failed to patch DockerCluster")
			if rerr == nil {
				rerr = err
			}
		}
	}()

	// Handle deleted clusters
	if !huaweiCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, clusterScope)
	}
	return r.reconcileNormal(ctx, clusterScope)
}

func (r *HuaweiClusterReconciler) reconcileNormal(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *HuaweiClusterReconciler) reconcileDelete(ctx context.Context, clusterScope *scope.ClusterScope) (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HuaweiClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HuaweiCluster{}).
		Named("huaweicluster").
		Complete(r)
}
