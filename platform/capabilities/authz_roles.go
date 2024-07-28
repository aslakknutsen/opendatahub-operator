package capabilities

import (
	"context"
	"fmt"

	rbacv1 "k8s.io/api/rbac/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
)

// CreateOrUpdateAuthzBindings.
func CreateOrUpdateAuthzBindings(ctx context.Context, cli client.Client, roleName string, schemas []ResourceSchema, metaOptions ...cluster.MetaOptions) error {
	if _, err := cluster.CreateOrUpdateClusterRole(ctx, cli, roleName, createAuthRules(schemas), metaOptions...); err != nil {
		return fmt.Errorf("failed creating cluster role: %w", err)
	}

	subjects, roleRef := createAuthzRoleBinding(roleName)
	if _, err := cluster.CreateOrUpdateClusterRoleBinding(ctx, cli, roleName, subjects, roleRef, metaOptions...); err != nil {
		return fmt.Errorf("failed creating cluster role binding: %w", err)
	}

	return nil
}

// DeleteAuthzBindings attempts to remove created authz role/bindings but does not fail if these are not existing in the cluster.
func DeleteAuthzBindings(ctx context.Context, cli client.Client, roleName string) error {
	if err := cluster.DeleteClusterRoleBinding(ctx, cli, roleName); !k8serr.IsNotFound(err) {
		return err
	}
	if err := cluster.DeleteClusterRole(ctx, cli, roleName); !k8serr.IsNotFound(err) {
		return err
	}

	return nil
}

func createAuthRules(schemas []ResourceSchema) []rbacv1.PolicyRule {
	apiGroups := make([]string, 0)
	resources := make([]string, 0)
	for _, schema := range schemas {
		apiGroups = append(apiGroups, schema.GroupVersionKind.Group)
		resources = append(resources, schema.Resources)
	}

	return []rbacv1.PolicyRule{
		{
			APIGroups: apiGroups,
			Resources: resources,
			Verbs:     []string{"get", "list", "watch"},
		},
	}
}

func createAuthzRoleBinding(roleName string) ([]rbacv1.Subject, rbacv1.RoleRef) {
	return []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      "odh-platform-manager",
				Namespace: "opendatahub",
			},
		},
		rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     roleName,
		}
}
