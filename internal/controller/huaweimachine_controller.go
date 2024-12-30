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
	"fmt"

	"github.com/pkg/errors"
	infrav1 "github.com/setoru/cluster-api-provider-huawei/api/v1alpha1"
	hwclient "github.com/setoru/cluster-api-provider-huawei/internal/cloud/client"
	"github.com/setoru/cluster-api-provider-huawei/internal/cloud/scope"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/cluster-api/util/patch"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// HuaweiMachineReconciler reconciles a HuaweiMachine object
type HuaweiMachineReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweimachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweimachines/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=huaweimachines/finalizers,verbs=update
// +kubebuilder:rbac:groups=cluster.x-k8s.io,resources=machines;machines/status,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the HuaweiMachine object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *HuaweiMachineReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, rerr error) {
	logger := log.FromContext(ctx)

	huaweiMachine := &infrav1.HuaweiMachine{}
	if err := r.Get(ctx, req.NamespacedName, huaweiMachine); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	machine, err := util.GetOwnerMachine(ctx, r.Client, huaweiMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine != nil {
		logger.Info("Waiting for Machine Controller to set OwnerRef on HuaweiMachine")
		return ctrl.Result{}, err
	}

	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, huaweiMachine.ObjectMeta)
	if err != nil {
		logger.Error(err, "HuaweiMachine owner Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, err
	}

	if cluster == nil {
		logger.Info(fmt.Sprintf("Please associate this machien with a cluster using the label %s: <name of cluster>", clusterv1.ClusterNameLabel))
		return ctrl.Result{}, nil
	}

	logger = logger.WithValues("cluster", cluster.Name)

	huaweiCluster := &infrav1.HuaweiCluster{}
	huaweiClusterName := client.ObjectKey{
		Namespace: huaweiMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}

	if err := r.Client.Get(ctx, huaweiClusterName, huaweiCluster); err != nil {
		logger.Error(err, "failed to get huawei cluster")
		return ctrl.Result{}, nil
	}

	secretName := huaweiCluster.Spec.CredentialsSecret.Name
	region := huaweiCluster.Spec.Region
	hwClient, err := hwclient.NewClient(r.Client, secretName, huaweiCluster.Namespace, region)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Create the scope.
	machineScope := &scope.MachineScope{
		Client:        r.Client,
		Logger:        logger,
		Cluster:       cluster,
		Machine:       machine,
		HuaweiMachine: huaweiMachine,
		HuaweiClient:  hwClient,
	}
	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}

	// Initialize the patch helper
	patchHelper, err := patch.NewHelper(huaweiMachine, r.Client)
	if err != nil {
		return ctrl.Result{}, err
	}

	defer func() {
		if err := patchHuaweiMachine(ctx, patchHelper, huaweiMachine); err != nil && rerr == nil {
			logger.Error(err, "failed to patch HuaweiMachine")
			rerr = err
		}
	}()

	// Return early if the object or Cluster is paused
	if annotations.IsPaused(cluster, huaweiMachine) {
		logger.Info("huaweiMachine or linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	// Add finalizer first if not exist to avoid the race condition between init and delete
	if !controllerutil.ContainsFinalizer(huaweiMachine, infrav1.MachineFinalizer) {
		controllerutil.AddFinalizer(huaweiMachine, infrav1.MachineFinalizer)
		return ctrl.Result{}, nil
	}

	// Handle deleted Machines
	if !huaweiMachine.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, machineScope)
	}
	return r.reconcileNormal(ctx, machineScope)
}

func (r *HuaweiMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.MachineScope) (_ ctrl.Result, retErr error) {
	return ctrl.Result{}, nil
}

func (r *HuaweiMachineReconciler) reconcileDelete(ctx context.Context, machineScope *scope.MachineScope) (_ ctrl.Result, retErr error) {
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *HuaweiMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.HuaweiMachine{}).
		Named("huaweimachine").
		Complete(r)
}

func patchHuaweiMachine(ctx context.Context, patchHelper *patch.Helper, huaweiMachine *infrav1.HuaweiMachine) error {
	// Always update the readyCondition by summarizing the state of other conditions.
	// A step counter is added to represent progress during the provisioning process (instead we are hiding the step counter during the deletion process).
	conditions.SetSummary(huaweiMachine,
		conditions.WithConditions(
			infrav1.InstanceReadyCondition,
		),
		conditions.WithStepCounterIf(huaweiMachine.ObjectMeta.DeletionTimestamp.IsZero() && huaweiMachine.Spec.ProviderID == nil),
	)
	// Patch the object, ignoring conflicts on the conditions owned by this controller.
	return patchHelper.Patch(
		ctx,
		huaweiMachine,
		patch.WithOwnedConditions{Conditions: []clusterv1.ConditionType{
			clusterv1.ReadyCondition,
			infrav1.InstanceReadyCondition,
		}},
	)
}
