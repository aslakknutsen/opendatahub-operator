package v1

import operatorv1 "github.com/openshift/api/operator/v1"

// ServiceMeshSpec configures Service Mesh.
type ServiceMeshSpec struct {
	// +kubebuilder:validation:Enum=Managed;Unmanaged;Removed
	// +kubebuilder:default=Removed
	ManagementState operatorv1.ManagementState `json:"managementState,omitempty"`
	// ControlPlane holds configuration of Service Mesh used by Opendatahub.
	ControlPlane ControlPlaneSpec `json:"controlPlane,omitempty"`
	// Auth holds configuration of authentication and authorization services
	// used by Service Mesh in Opendatahub.
	Auth AuthSpec `json:"auth,omitempty"`
}

type ControlPlaneSpec struct {
	// Name is a name Service Mesh Control Plane. Defaults to "data-science-smcp".
	// +kubebuilder:default=data-science-smcp
	Name string `json:"name,omitempty"`
	// Namespace is a namespace where Service Mesh is deployed. Defaults to "istio-system".
	// +kubebuilder:default=istio-system
	Namespace string `json:"namespace,omitempty"`
	// MetricsCollection specifies if metrics from components on the Mesh namespace
	// should be collected. Setting the value to "Istio" will collect metrics from the
	// control plane and any proxies on the Mesh namespace (like gateway pods). Setting
	// to "None" will disable metrics collection.
	// +kubebuilder:validation:Enum=Istio;None
	// +kubebuilder:default=Istio
	MetricsCollection string `json:"metricsCollection,omitempty"`
	//IngressGatewaySpec specifices where the Service Meshs Ingress gateway should behave
	IngressGateway RealIngressGatewaySpec `json:"ingressGateway,omitempty"`
}

type RealIngressGatewaySpec struct {
	// Namespace is a namespace where the Service Mesh Ingress should be deployed is deployed. Defaults to "opendatahub-ingress".
	// +kubebuilder:default=opendatahub-ingress
	Namespace string `json:"namespace,omitempty"`
	// Name is the name of the Ingress Gateway Service
	// +kubebuilder:default=istio-ingressgateway
	Name string `json:"name,omitempty"` // TODO: Not API
	// LabelSelectorKey is a key:value defining the label to use for the ingress gateway objects
	// +kubebuilder:default="opendatahub"
	LabelSelectorKey string `json:"labelSelectorKey,omitempty"` // TODO: Not API
	// LabelSelectorValue is a key:value defining the label to use for the ingress gateway objects
	// +kubebuilder:default="ingressgateway"
	LabelSelectorValue string `json:"labelSelectorValue,omitempty"` // TODO: Not API
	//Gateway holds configuration for the default Gateway
	//TODO: Should Ceritifacte be extracted at a higher level, a more common place between Istio, KNative ++?
	Gateway IngressGatewaySpec `json:"gateway,omitempty"`
}

// TODO: Isn't this really GatewaySpec?
// IngressGatewaySpec represents the configuration of the Ingress Gateways.
type IngressGatewaySpec struct {
	// Domain specifies the DNS name for intercepting ingress requests coming from
	// outside the cluster. Most likely, you will want to use a wildcard name,
	// like *.example.com. If not set, the domain of the OpenShift Ingress is used.
	// If you choose to generate a certificate, this is the domain used for the certificate request.
	Domain string `json:"domain,omitempty"`
	// Certificate specifies configuration of the TLS certificate securing communications of
	// the for Ingress Gateway.
	Certificate CertificateSpec `json:"certificate,omitempty"`
}

type AuthSpec struct {
	// Namespace where it is deployed. If not provided, the default is to
	// use '-auth-provider' suffix on the ApplicationsNamespace of the DSCI.
	Namespace string `json:"namespace,omitempty"`
	// Audiences is a list of the identifiers that the resource server presented
	// with the token identifies as. Audience-aware token authenticators will verify
	// that the token was intended for at least one of the audiences in this list.
	// If no audiences are provided, the audience will default to the audience of the
	// Kubernetes apiserver (kubernetes.default.svc).
	// +kubebuilder:default={"https://kubernetes.default.svc"}
	Audiences *[]string `json:"audiences,omitempty"`
}
