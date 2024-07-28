package capabilities

import (
	"context"
	"encoding/json"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opendatahub-io/opendatahub-operator/v2/pkg/cluster"
)

// Consumer

func NewRouting(available bool) RoutingCapability {
	return RoutingCapability{available: available}
}

type Routing interface {
	Availability
	JSONSerializable
	Route(routeResources ...ResourceSchema)
}

// Producer

var _ Routing = (*RoutingCapability)(nil)

type RoutingCapability struct {
	available      bool
	routeResources []ResourceSchema
}

func (a *RoutingCapability) IsAvailable() bool {
	return a.available
}

func (a *RoutingCapability) Route(routeResource ...ResourceSchema) {
	a.routeResources = routeResource
}

func (a *RoutingCapability) AsJSON() ([]byte, error) {
	return json.Marshal(a.routeResources)
}

var _ Handler = (*RoutingCapability)(nil)

func (a *RoutingCapability) IsRequired() bool {
	return len(a.routeResources) > 0
}

// Reconcile ensures Authorization capability and component-specific configuration is wired when needed.
func (a *RoutingCapability) Reconcile(ctx context.Context, cli client.Client, metaOptions ...cluster.MetaOptions) error {
	roleName := "platform-routing-resources-watcher"
	if a.IsRequired() {
		return CreateOrUpdateAuthzBindings(ctx, cli, roleName, a.routeResources, metaOptions...)
	}

	return DeleteAuthzBindings(ctx, cli, roleName)
}
