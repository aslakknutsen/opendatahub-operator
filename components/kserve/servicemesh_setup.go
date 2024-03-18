package kserve

import (
	"fmt"
	"path"

	operatorv1 "github.com/openshift/api/operator/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dsciv1 "github.com/opendatahub-io/opendatahub-operator/v2/apis/dscinitialization/v1"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature"
	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/feature/servicemesh"
)

func (k *Kserve) configureServiceMesh(cli client.Client, dscispec *dsciv1.DSCInitializationSpec) error {
	if dscispec.ServiceMesh.ManagementState == operatorv1.Managed && k.GetManagementState() == operatorv1.Managed {
		serviceMeshInitializer := feature.ComponentFeaturesHandler(k.GetComponentName(), dscispec.ApplicationsNamespace, k.defineServiceMeshFeatures(cli, dscispec))
		return serviceMeshInitializer.Apply()
	}
	if dscispec.ServiceMesh.ManagementState == operatorv1.Unmanaged && k.GetManagementState() == operatorv1.Managed {
		return nil
	}

	return k.removeServiceMeshConfigurations(cli, dscispec)
}

func (k *Kserve) removeServiceMeshConfigurations(cli client.Client, dscispec *dsciv1.DSCInitializationSpec) error {
	serviceMeshInitializer := feature.ComponentFeaturesHandler(k.GetComponentName(), dscispec.ApplicationsNamespace, k.defineServiceMeshFeatures(cli, dscispec))
	return serviceMeshInitializer.Delete()
}

func (k *Kserve) defineServiceMeshFeatures(cli client.Client, dscispec *dsciv1.DSCInitializationSpec) feature.FeaturesProvider {
	return func(registry feature.FeaturesRegistry) error {
		authorinoInstalled, err := cluster.SubscriptionExists(cli, "authorino-operator")
		if err != nil {
			return fmt.Errorf("failed to list subscriptions %w", err)
		}

		if authorinoInstalled {
			kserveExtAuthzErr := registry.Add(feature.Define("kserve-external-authz").
				ManifestsLocation(Resources.Location).
				Manifests(
					path.Join(Resources.ServiceMeshDir, "activator-envoyfilter.tmpl.yaml"),
					path.Join(Resources.ServiceMeshDir, "envoy-oauth-temp-fix.tmpl.yaml"),
					path.Join(Resources.ServiceMeshDir, "kserve-predictor-authorizationpolicy.tmpl.yaml"),
					path.Join(Resources.ServiceMeshDir, "z-migrations"),
				).
				WithData(
					feature.Entry("Domain", cluster.GetDomain),
					servicemesh.FeatureData.ControlPlane.Create(dscispec).AsAction(),
				).
				WithData(
					servicemesh.FeatureData.Authorization.All(dscispec)...,
				),
			)

			if kserveExtAuthzErr != nil {
				return kserveExtAuthzErr
			}
		} else {
			fmt.Println("WARN: Authorino operator is not installed on the cluster, skipping authorization capability")
		}

		return registry.Add(feature.Define("kserve-temporary-fixes").
			ManifestsLocation(Resources.Location).
			Manifests(
				path.Join(Resources.ServiceMeshDir, "grpc-envoyfilter-temp-fix.tmpl.yaml"),
			).
			WithData(servicemesh.FeatureData.ControlPlane.Create(dscispec).AsAction()),
		)
	}
}
