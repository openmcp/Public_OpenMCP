// +build !ignore_autogenerated

// Code generated by operator-sdk. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenMCPLoadbalancing) DeepCopyInto(out *OpenMCPLoadbalancing) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenMCPLoadbalancing.
func (in *OpenMCPLoadbalancing) DeepCopy() *OpenMCPLoadbalancing {
	if in == nil {
		return nil
	}
	out := new(OpenMCPLoadbalancing)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OpenMCPLoadbalancing) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenMCPLoadbalancingList) DeepCopyInto(out *OpenMCPLoadbalancingList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]OpenMCPLoadbalancing, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenMCPLoadbalancingList.
func (in *OpenMCPLoadbalancingList) DeepCopy() *OpenMCPLoadbalancingList {
	if in == nil {
		return nil
	}
	out := new(OpenMCPLoadbalancingList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *OpenMCPLoadbalancingList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenMCPLoadbalancingSpec) DeepCopyInto(out *OpenMCPLoadbalancingSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenMCPLoadbalancingSpec.
func (in *OpenMCPLoadbalancingSpec) DeepCopy() *OpenMCPLoadbalancingSpec {
	if in == nil {
		return nil
	}
	out := new(OpenMCPLoadbalancingSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenMCPLoadbalancingStatus) DeepCopyInto(out *OpenMCPLoadbalancingStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenMCPLoadbalancingStatus.
func (in *OpenMCPLoadbalancingStatus) DeepCopy() *OpenMCPLoadbalancingStatus {
	if in == nil {
		return nil
	}
	out := new(OpenMCPLoadbalancingStatus)
	in.DeepCopyInto(out)
	return out
}