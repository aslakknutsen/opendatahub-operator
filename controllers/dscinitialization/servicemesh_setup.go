package dscinitialization

import (
	"fmt"
	"path"

	operatorv1 "github.com/openshift/api/operator/v1"
	conditionsv1 "github.com/openshift/custom-resource-status/conditions/v1"
	corev1 "k8s.io/api/core/v1"

	dsciv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/dscinitialization/v1"
	"github.com/opendatahub-io/opendatahub-operator/v2/controllers/status"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature/manifest"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature/servicemesh"
)

func (r *DSCInitializationReconciler) configureServiceMesh(instance *dsciv1.DSCInitialization) error {
	serviceMeshManagementState := operatorv1.Removed
	if instance.Spec.ServiceMesh != nil {
		serviceMeshManagementState = instance.Spec.ServiceMesh.ManagementState
	} else {
		r.Log.Info("ServiceMesh is not configured in DSCI, same as default to 'Removed'")
	}

	switch serviceMeshManagementState {
	case operatorv1.Managed:

		capabilities := []*feature.HandlerWithReporter[*dsciv1.DSCInitialization]{
			r.serviceMeshCapability(instance, serviceMeshCondition(status.ConfiguredReason, "Service Mesh configured")),
		}

		authzCapability, err := r.authorizationCapability(instance, authorizationCondition(status.ConfiguredReason, "Service Mesh Authorization configured"))
		if err != nil {
			return err
		}
		capabilities = append(capabilities, authzCapability)

		for _, capability := range capabilities {
			capabilityErr := capability.Apply()
			if capabilityErr != nil {
				r.Log.Error(capabilityErr, "failed applying service mesh resources")
				r.Recorder.Eventf(instance, corev1.EventTypeWarning, "DSCInitializationReconcileError", "failed applying service mesh resources")
				return capabilityErr
			}
		}

	case operatorv1.Unmanaged:
		r.Log.Info("ServiceMesh CR is not configured by the operator, we won't do anything")
	case operatorv1.Removed:
		r.Log.Info("existing ServiceMesh CR (owned by operator) will be removed")
		if err := r.removeServiceMesh(instance); err != nil {
			return err
		}
	}
	return nil
}

func (r *DSCInitializationReconciler) removeServiceMesh(instance *dsciv1.DSCInitialization) error {
	// on condition of Managed, do not handle Removed when set to Removed it trigger DSCI reconcile to clean up
	if instance.Spec.ServiceMesh == nil {
		return nil
	}
	if instance.Spec.ServiceMesh.ManagementState == operatorv1.Managed {
		capabilities := []*feature.HandlerWithReporter[*dsciv1.DSCInitialization]{
			r.serviceMeshCapability(instance, serviceMeshCondition(status.RemovedReason, "Service Mesh removed")),
		}

		authzCapability, err := r.authorizationCapability(instance, authorizationCondition(status.RemovedReason, "Service Mesh Authorization removed"))
		if err != nil {
			return err
		}

		capabilities = append(capabilities, authzCapability)

		for _, capability := range capabilities {
			capabilityErr := capability.Delete()
			if capabilityErr != nil {
				r.Log.Error(capabilityErr, "failed deleting service mesh resources")
				r.Recorder.Eventf(instance, corev1.EventTypeWarning, "DSCInitializationReconcileError", "failed deleting service mesh resources")

				return capabilityErr
			}
		}
	}
	return nil
}

func (r *DSCInitializationReconciler) serviceMeshCapability(dsci *dsciv1.DSCInitialization, initialCondition *conditionsv1.Condition) *feature.HandlerWithReporter[*dsciv1.DSCInitialization] { //nolint:lll // Reason: generics are long
	return feature.NewHandlerWithReporter(
		feature.ClusterFeaturesHandler(dsci, r.serviceMeshCapabilityFeatures(dsci)),
		createCapabilityReporter(r.Client, dsci, initialCondition),
	)
}

func (r *DSCInitializationReconciler) authorizationCapability(instance *dsciv1.DSCInitialization, condition *conditionsv1.Condition) (*feature.HandlerWithReporter[*dsciv1.DSCInitialization], error) { //nolint:lll // Reason: generics are long
	authorinoInstalled, err := cluster.SubscriptionExists(r.Client, "authorino-operator")
	if err != nil {
		return nil, fmt.Errorf("failed to list subscriptions %w", err)
	}

	if !authorinoInstalled {
		authzMissingOperatorCondition := &conditionsv1.Condition{
			Type:    status.CapabilityServiceMeshAuthorization,
			Status:  corev1.ConditionFalse,
			Reason:  status.MissingOperatorReason,
			Message: "Authorino operator is not installed on the cluster, skipping authorization capability",
		}

		return feature.NewHandlerWithReporter(
			// EmptyFeaturesHandler acts as all the authorization features are disabled (calling Apply/Delete has no actual effect on the cluster)
			// but it's going to be reported as CapabilityServiceMeshAuthorization/MissingOperator condition/reason
			feature.EmptyFeaturesHandler,
			createCapabilityReporter(r.Client, instance, authzMissingOperatorCondition),
		), nil
	}

	return feature.NewHandlerWithReporter(
		feature.ClusterFeaturesHandler(instance, r.authorizationFeatures(instance)),
		createCapabilityReporter(r.Client, instance, condition),
	), nil
}

func (r *DSCInitializationReconciler) serviceMeshCapabilityFeatures(dsci *dsciv1.DSCInitialization) feature.FeaturesProvider {
	return func(registry feature.FeaturesRegistry) error {
		serviceMeshSpec := dsci.Spec.ServiceMesh

		smcp := feature.Define("mesh-control-plane-creation").
			Manifests(
				manifest.Location(Templates.Location).
					Include(
						path.Join(Templates.ServiceMeshDir),
					),
			).
			WithData(servicemesh.FeatureData.ControlPlane.Create(&dsci.Spec).AsAction()).
			PreConditions(
				servicemesh.EnsureServiceMeshOperatorInstalled,
				feature.CreateNamespaceIfNotExists(serviceMeshSpec.ControlPlane.Namespace),
			).
			PostConditions(
				feature.WaitForPodsToBeReady(serviceMeshSpec.ControlPlane.Namespace),
			)

		if serviceMeshSpec.ControlPlane.MetricsCollection == "Istio" {
			metricsCollectionErr := registry.Add(feature.Define("mesh-metrics-collection").
				Manifests(
					manifest.Location(Templates.Location).
						Include(
							path.Join(Templates.MetricsDir),
						),
				).
				WithData(
					servicemesh.FeatureData.ControlPlane.Create(&dsci.Spec).AsAction(),
				).
				PreConditions(
					servicemesh.EnsureServiceMeshInstalled,
				))

			if metricsCollectionErr != nil {
				return metricsCollectionErr
			}
		}

		cfgMap := feature.Define("mesh-shared-configmap").
			WithResources(servicemesh.MeshRefs, servicemesh.AuthRefs).
			WithData(
				servicemesh.FeatureData.ControlPlane.Create(&dsci.Spec).AsAction(),
			).
			WithData(
				servicemesh.FeatureData.Authorization.All(&dsci.Spec)...,
			)

		return registry.Add(smcp, cfgMap)
	}
}

func (r *DSCInitializationReconciler) authorizationFeatures(dsci *dsciv1.DSCInitialization) feature.FeaturesProvider {
	return func(registry feature.FeaturesRegistry) error {
		serviceMeshSpec := dsci.Spec.ServiceMesh

		return registry.Add(
			feature.Define("mesh-control-plane-external-authz").
				Manifests(
					manifest.Location(Templates.Location).
						Include(
							path.Join(Templates.AuthorinoDir, "auth-smm.tmpl.yaml"),
							path.Join(Templates.AuthorinoDir, "base"),
							path.Join(Templates.AuthorinoDir, "mesh-authz-ext-provider.patch.tmpl.yaml"),
						),
				).
				WithData(
					servicemesh.FeatureData.ControlPlane.Create(&dsci.Spec).AsAction(),
				).
				WithData(
					servicemesh.FeatureData.Authorization.All(&dsci.Spec)...,
				).
				PreConditions(
					feature.EnsureOperatorIsInstalled("authorino-operator"),
					servicemesh.EnsureServiceMeshInstalled,
					servicemesh.EnsureAuthNamespaceExists,
				).
				PostConditions(
					feature.WaitForPodsToBeReady(serviceMeshSpec.ControlPlane.Namespace),
				).
				OnDelete(
					servicemesh.RemoveExtensionProvider,
				),
			// We do not have the control over deployment resource creation.
			// It is created by Authorino operator using Authorino CR and labels are not propagated from Authorino CR to spec.template
			//
			// To make it part of Service Mesh we have to patch it with injection
			// enabled instead, otherwise it will not have proxy pod injected.
			feature.Define("enable-proxy-injection-in-authorino-deployment").
				PreConditions(
					servicemesh.EnsureAuthNamespaceExists,
					func(f *feature.Feature) error {
						return feature.WaitForPodsToBeReady(serviceMeshSpec.Auth.Namespace)(f)
					}).
				Manifests(
					manifest.Location(Templates.Location).
						Include(path.Join(Templates.AuthorinoDir, "deployment.injection.patch.tmpl.yaml")),
				),
		)
	}
}
