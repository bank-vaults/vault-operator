//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AWSUnsealConfig) DeepCopyInto(out *AWSUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AWSUnsealConfig.
func (in *AWSUnsealConfig) DeepCopy() *AWSUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(AWSUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AlibabaUnsealConfig) DeepCopyInto(out *AlibabaUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AlibabaUnsealConfig.
func (in *AlibabaUnsealConfig) DeepCopy() *AlibabaUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(AlibabaUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AzureUnsealConfig) DeepCopyInto(out *AzureUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AzureUnsealConfig.
func (in *AzureUnsealConfig) DeepCopy() *AzureUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(AzureUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialsConfig) DeepCopyInto(out *CredentialsConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialsConfig.
func (in *CredentialsConfig) DeepCopy() *CredentialsConfig {
	if in == nil {
		return nil
	}
	out := new(CredentialsConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EmbeddedObjectMetadata) DeepCopyInto(out *EmbeddedObjectMetadata) {
	*out = *in
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EmbeddedObjectMetadata.
func (in *EmbeddedObjectMetadata) DeepCopy() *EmbeddedObjectMetadata {
	if in == nil {
		return nil
	}
	out := new(EmbeddedObjectMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EmbeddedPersistentVolumeClaim) DeepCopyInto(out *EmbeddedPersistentVolumeClaim) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.EmbeddedObjectMetadata.DeepCopyInto(&out.EmbeddedObjectMetadata)
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EmbeddedPersistentVolumeClaim.
func (in *EmbeddedPersistentVolumeClaim) DeepCopy() *EmbeddedPersistentVolumeClaim {
	if in == nil {
		return nil
	}
	out := new(EmbeddedPersistentVolumeClaim)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EmbeddedPodSpec) DeepCopyInto(out *EmbeddedPodSpec) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.InitContainers != nil {
		in, out := &in.InitContainers, &out.InitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Containers != nil {
		in, out := &in.Containers, &out.Containers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.EphemeralContainers != nil {
		in, out := &in.EphemeralContainers, &out.EphemeralContainers
		*out = make([]v1.EphemeralContainer, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TerminationGracePeriodSeconds != nil {
		in, out := &in.TerminationGracePeriodSeconds, &out.TerminationGracePeriodSeconds
		*out = new(int64)
		**out = **in
	}
	if in.ActiveDeadlineSeconds != nil {
		in, out := &in.ActiveDeadlineSeconds, &out.ActiveDeadlineSeconds
		*out = new(int64)
		**out = **in
	}
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.AutomountServiceAccountToken != nil {
		in, out := &in.AutomountServiceAccountToken, &out.AutomountServiceAccountToken
		*out = new(bool)
		**out = **in
	}
	if in.ShareProcessNamespace != nil {
		in, out := &in.ShareProcessNamespace, &out.ShareProcessNamespace
		*out = new(bool)
		**out = **in
	}
	if in.SecurityContext != nil {
		in, out := &in.SecurityContext, &out.SecurityContext
		*out = new(v1.PodSecurityContext)
		(*in).DeepCopyInto(*out)
	}
	if in.ImagePullSecrets != nil {
		in, out := &in.ImagePullSecrets, &out.ImagePullSecrets
		*out = make([]v1.LocalObjectReference, len(*in))
		copy(*out, *in)
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.HostAliases != nil {
		in, out := &in.HostAliases, &out.HostAliases
		*out = make([]v1.HostAlias, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Priority != nil {
		in, out := &in.Priority, &out.Priority
		*out = new(int32)
		**out = **in
	}
	if in.DNSConfig != nil {
		in, out := &in.DNSConfig, &out.DNSConfig
		*out = new(v1.PodDNSConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.ReadinessGates != nil {
		in, out := &in.ReadinessGates, &out.ReadinessGates
		*out = make([]v1.PodReadinessGate, len(*in))
		copy(*out, *in)
	}
	if in.RuntimeClassName != nil {
		in, out := &in.RuntimeClassName, &out.RuntimeClassName
		*out = new(string)
		**out = **in
	}
	if in.EnableServiceLinks != nil {
		in, out := &in.EnableServiceLinks, &out.EnableServiceLinks
		*out = new(bool)
		**out = **in
	}
	if in.PreemptionPolicy != nil {
		in, out := &in.PreemptionPolicy, &out.PreemptionPolicy
		*out = new(v1.PreemptionPolicy)
		**out = **in
	}
	if in.Overhead != nil {
		in, out := &in.Overhead, &out.Overhead
		*out = make(v1.ResourceList, len(*in))
		for key, val := range *in {
			(*out)[key] = val.DeepCopy()
		}
	}
	if in.TopologySpreadConstraints != nil {
		in, out := &in.TopologySpreadConstraints, &out.TopologySpreadConstraints
		*out = make([]v1.TopologySpreadConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SetHostnameAsFQDN != nil {
		in, out := &in.SetHostnameAsFQDN, &out.SetHostnameAsFQDN
		*out = new(bool)
		**out = **in
	}
	if in.OS != nil {
		in, out := &in.OS, &out.OS
		*out = new(v1.PodOS)
		**out = **in
	}
	if in.HostUsers != nil {
		in, out := &in.HostUsers, &out.HostUsers
		*out = new(bool)
		**out = **in
	}
	if in.SchedulingGates != nil {
		in, out := &in.SchedulingGates, &out.SchedulingGates
		*out = make([]v1.PodSchedulingGate, len(*in))
		copy(*out, *in)
	}
	if in.ResourceClaims != nil {
		in, out := &in.ResourceClaims, &out.ResourceClaims
		*out = make([]v1.PodResourceClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EmbeddedPodSpec.
func (in *EmbeddedPodSpec) DeepCopy() *EmbeddedPodSpec {
	if in == nil {
		return nil
	}
	out := new(EmbeddedPodSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GoogleUnsealConfig) DeepCopyInto(out *GoogleUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GoogleUnsealConfig.
func (in *GoogleUnsealConfig) DeepCopy() *GoogleUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(GoogleUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HSMUnsealConfig) DeepCopyInto(out *HSMUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HSMUnsealConfig.
func (in *HSMUnsealConfig) DeepCopy() *HSMUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(HSMUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubernetesUnsealConfig) DeepCopyInto(out *KubernetesUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubernetesUnsealConfig.
func (in *KubernetesUnsealConfig) DeepCopy() *KubernetesUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(KubernetesUnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resources) DeepCopyInto(out *Resources) {
	*out = *in
	if in.Vault != nil {
		in, out := &in.Vault, &out.Vault
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.BankVaults != nil {
		in, out := &in.BankVaults, &out.BankVaults
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.HSMDaemon != nil {
		in, out := &in.HSMDaemon, &out.HSMDaemon
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.PrometheusExporter != nil {
		in, out := &in.PrometheusExporter, &out.PrometheusExporter
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.FluentD != nil {
		in, out := &in.FluentD, &out.FluentD
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UnsealConfig) DeepCopyInto(out *UnsealConfig) {
	*out = *in
	in.Options.DeepCopyInto(&out.Options)
	out.Kubernetes = in.Kubernetes
	if in.Google != nil {
		in, out := &in.Google, &out.Google
		*out = new(GoogleUnsealConfig)
		**out = **in
	}
	if in.Alibaba != nil {
		in, out := &in.Alibaba, &out.Alibaba
		*out = new(AlibabaUnsealConfig)
		**out = **in
	}
	if in.Azure != nil {
		in, out := &in.Azure, &out.Azure
		*out = new(AzureUnsealConfig)
		**out = **in
	}
	if in.AWS != nil {
		in, out := &in.AWS, &out.AWS
		*out = new(AWSUnsealConfig)
		**out = **in
	}
	if in.Vault != nil {
		in, out := &in.Vault, &out.Vault
		*out = new(VaultUnsealConfig)
		**out = **in
	}
	if in.HSM != nil {
		in, out := &in.HSM, &out.HSM
		*out = new(HSMUnsealConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnsealConfig.
func (in *UnsealConfig) DeepCopy() *UnsealConfig {
	if in == nil {
		return nil
	}
	out := new(UnsealConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UnsealOptions) DeepCopyInto(out *UnsealOptions) {
	*out = *in
	if in.PreFlightChecks != nil {
		in, out := &in.PreFlightChecks, &out.PreFlightChecks
		*out = new(bool)
		**out = **in
	}
	if in.StoreRootToken != nil {
		in, out := &in.StoreRootToken, &out.StoreRootToken
		*out = new(bool)
		**out = **in
	}
	if in.SecretThreshold != nil {
		in, out := &in.SecretThreshold, &out.SecretThreshold
		*out = new(uint)
		**out = **in
	}
	if in.SecretShares != nil {
		in, out := &in.SecretShares, &out.SecretShares
		*out = new(uint)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UnsealOptions.
func (in *UnsealOptions) DeepCopy() *UnsealOptions {
	if in == nil {
		return nil
	}
	out := new(UnsealOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Vault) DeepCopyInto(out *Vault) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Vault.
func (in *Vault) DeepCopy() *Vault {
	if in == nil {
		return nil
	}
	out := new(Vault)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Vault) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultList) DeepCopyInto(out *VaultList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Vault, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultList.
func (in *VaultList) DeepCopy() *VaultList {
	if in == nil {
		return nil
	}
	out := new(VaultList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VaultList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultSpec) DeepCopyInto(out *VaultSpec) {
	*out = *in
	if in.BankVaultsVolumeMounts != nil {
		in, out := &in.BankVaultsVolumeMounts, &out.BankVaultsVolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.WatchedSecretsLabels != nil {
		in, out := &in.WatchedSecretsLabels, &out.WatchedSecretsLabels
		*out = make([]map[string]string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
	if in.WatchedSecretsAnnotations != nil {
		in, out := &in.WatchedSecretsAnnotations, &out.WatchedSecretsAnnotations
		*out = make([]map[string]string, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
			}
		}
	}
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VaultAnnotations != nil {
		in, out := &in.VaultAnnotations, &out.VaultAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VaultLabels != nil {
		in, out := &in.VaultLabels, &out.VaultLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VaultPodSpec != nil {
		in, out := &in.VaultPodSpec, &out.VaultPodSpec
		*out = new(EmbeddedPodSpec)
		(*in).DeepCopyInto(*out)
	}
	in.VaultContainerSpec.DeepCopyInto(&out.VaultContainerSpec)
	if in.VaultConfigurerAnnotations != nil {
		in, out := &in.VaultConfigurerAnnotations, &out.VaultConfigurerAnnotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VaultConfigurerLabels != nil {
		in, out := &in.VaultConfigurerLabels, &out.VaultConfigurerLabels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.VaultConfigurerPodSpec != nil {
		in, out := &in.VaultConfigurerPodSpec, &out.VaultConfigurerPodSpec
		*out = new(EmbeddedPodSpec)
		(*in).DeepCopyInto(*out)
	}
	in.Config.DeepCopyInto(&out.Config)
	in.ExternalConfig.DeepCopyInto(&out.ExternalConfig)
	in.UnsealConfig.DeepCopyInto(&out.UnsealConfig)
	out.CredentialsConfig = in.CredentialsConfig
	if in.EnvsConfig != nil {
		in, out := &in.EnvsConfig, &out.EnvsConfig
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.SecurityContext.DeepCopyInto(&out.SecurityContext)
	if in.ServicePorts != nil {
		in, out := &in.ServicePorts, &out.ServicePorts
		*out = make(map[string]int32, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	in.NodeAffinity.DeepCopyInto(&out.NodeAffinity)
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]EmbeddedPersistentVolumeClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VaultEnvsConfig != nil {
		in, out := &in.VaultEnvsConfig, &out.VaultEnvsConfig
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SecretInitsConfig != nil {
		in, out := &in.SecretInitsConfig, &out.SecretInitsConfig
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.SidecarEnvsConfig != nil {
		in, out := &in.SidecarEnvsConfig, &out.SidecarEnvsConfig
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(Resources)
		(*in).DeepCopyInto(*out)
	}
	if in.Ingress != nil {
		in, out := &in.Ingress, &out.Ingress
		*out = new(Ingress)
		(*in).DeepCopyInto(*out)
	}
	if in.TLSAdditionalHosts != nil {
		in, out := &in.TLSAdditionalHosts, &out.TLSAdditionalHosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.CANamespaces != nil {
		in, out := &in.CANamespaces, &out.CANamespaces
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VaultContainers != nil {
		in, out := &in.VaultContainers, &out.VaultContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VaultInitContainers != nil {
		in, out := &in.VaultInitContainers, &out.VaultInitContainers
		*out = make([]v1.Container, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultSpec.
func (in *VaultSpec) DeepCopy() *VaultSpec {
	if in == nil {
		return nil
	}
	out := new(VaultSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultStatus) DeepCopyInto(out *VaultStatus) {
	*out = *in
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.ComponentCondition, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultStatus.
func (in *VaultStatus) DeepCopy() *VaultStatus {
	if in == nil {
		return nil
	}
	out := new(VaultStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VaultUnsealConfig) DeepCopyInto(out *VaultUnsealConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VaultUnsealConfig.
func (in *VaultUnsealConfig) DeepCopy() *VaultUnsealConfig {
	if in == nil {
		return nil
	}
	out := new(VaultUnsealConfig)
	in.DeepCopyInto(out)
	return out
}
