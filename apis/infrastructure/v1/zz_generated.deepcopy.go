//go:build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1

import ()

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AuthSpec) DeepCopyInto(out *AuthSpec) {
	*out = *in
	if in.Audiences != nil {
		in, out := &in.Audiences, &out.Audiences
		*out = new([]string)
		if **in != nil {
			in, out := *in, *out
			*out = make([]string, len(*in))
			copy(*out, *in)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AuthSpec.
func (in *AuthSpec) DeepCopy() *AuthSpec {
	if in == nil {
		return nil
	}
	out := new(AuthSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertificateSpec) DeepCopyInto(out *CertificateSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertificateSpec.
func (in *CertificateSpec) DeepCopy() *CertificateSpec {
	if in == nil {
		return nil
	}
	out := new(CertificateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ControlPlaneSpec) DeepCopyInto(out *ControlPlaneSpec) {
	*out = *in
	out.IngressGateway = in.IngressGateway
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ControlPlaneSpec.
func (in *ControlPlaneSpec) DeepCopy() *ControlPlaneSpec {
	if in == nil {
		return nil
	}
	out := new(ControlPlaneSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IngressGatewaySpec) DeepCopyInto(out *IngressGatewaySpec) {
	*out = *in
	out.Certificate = in.Certificate
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IngressGatewaySpec.
func (in *IngressGatewaySpec) DeepCopy() *IngressGatewaySpec {
	if in == nil {
		return nil
	}
	out := new(IngressGatewaySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RealIngressGatewaySpec) DeepCopyInto(out *RealIngressGatewaySpec) {
	*out = *in
	out.Gateway = in.Gateway
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RealIngressGatewaySpec.
func (in *RealIngressGatewaySpec) DeepCopy() *RealIngressGatewaySpec {
	if in == nil {
		return nil
	}
	out := new(RealIngressGatewaySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceMeshSpec) DeepCopyInto(out *ServiceMeshSpec) {
	*out = *in
	out.ControlPlane = in.ControlPlane
	in.Auth.DeepCopyInto(&out.Auth)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceMeshSpec.
func (in *ServiceMeshSpec) DeepCopy() *ServiceMeshSpec {
	if in == nil {
		return nil
	}
	out := new(ServiceMeshSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServingSpec) DeepCopyInto(out *ServingSpec) {
	*out = *in
	out.IngressGateway = in.IngressGateway
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServingSpec.
func (in *ServingSpec) DeepCopy() *ServingSpec {
	if in == nil {
		return nil
	}
	out := new(ServingSpec)
	in.DeepCopyInto(out)
	return out
}
