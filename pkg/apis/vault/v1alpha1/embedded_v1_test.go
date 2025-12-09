// Copyright © 2025 Bank-Vaults Maintainers
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

package v1alpha1

import (
	"testing"

	"github.com/bank-vaults/vault-operator/pkg/utils"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestToPodSpec(t *testing.T) {
	eps := &EmbeddedPodSpec{
		Volumes:                       []v1.Volume{{Name: "vol1"}},
		InitContainers:                []v1.Container{{Name: "init"}},
		Containers:                    []v1.Container{{Name: "main"}},
		EphemeralContainers:           []v1.EphemeralContainer{{EphemeralContainerCommon: v1.EphemeralContainerCommon{Name: "eph"}}},
		RestartPolicy:                 v1.RestartPolicyOnFailure,
		TerminationGracePeriodSeconds: utils.To(int64(42)),
		ActiveDeadlineSeconds:         utils.To(int64(100)),
		DNSPolicy:                     v1.DNSClusterFirst,
		NodeSelector:                  map[string]string{"k": "v"},
		ServiceAccountName:            "svc",
		DeprecatedServiceAccount:      "svc-deprecated",
		AutomountServiceAccountToken:  utils.To(true),
		NodeName:                      "node",
		HostNetwork:                   true,
		HostPID:                       true,
		HostIPC:                       true,
		ShareProcessNamespace:         utils.To(true),
		SecurityContext:               &v1.PodSecurityContext{RunAsUser: utils.To(int64(1000))},
		ImagePullSecrets:              []v1.LocalObjectReference{{Name: "pullsecret"}},
		Hostname:                      "host",
		Subdomain:                     "sub",
		Affinity:                      &v1.Affinity{},
		SchedulerName:                 "sched",
		Tolerations:                   []v1.Toleration{{Key: "key"}},
		HostAliases:                   []v1.HostAlias{{IP: "127.0.0.1"}},
		PriorityClassName:             "high",
		Priority:                      utils.To(int32(10)),
		DNSConfig:                     &v1.PodDNSConfig{},
		ReadinessGates:                []v1.PodReadinessGate{{ConditionType: "ready"}},
		RuntimeClassName:              utils.To("runtime"),
		EnableServiceLinks:            utils.To(false),
		PreemptionPolicy:              utils.To(v1.PreemptLowerPriority),
		Overhead:                      func() v1.ResourceList { qv, _ := resource.ParseQuantity("1"); return v1.ResourceList{"cpu": qv} }(),
		TopologySpreadConstraints:     []v1.TopologySpreadConstraint{{MaxSkew: 1}},
		SetHostnameAsFQDN:             utils.To(true),
		OS:                            &v1.PodOS{Name: "linux"},
		HostUsers:                     utils.To(false),
		SchedulingGates:               []v1.PodSchedulingGate{{Name: "gate"}},
		ResourceClaims:                []v1.PodResourceClaim{{Name: "claim"}},
		Resources:                     &v1.ResourceRequirements{},
	}

	ps := eps.ToPodSpec()

	t.Run("Volumes copied", func(t *testing.T) {
		require.Equal(t, eps.Volumes, ps.Volumes)
	})
	t.Run("InitContainers copied", func(t *testing.T) {
		require.Equal(t, eps.InitContainers, ps.InitContainers)
	})
	t.Run("Containers copied", func(t *testing.T) {
		require.Equal(t, eps.Containers, ps.Containers)
	})
	t.Run("EphemeralContainers copied", func(t *testing.T) {
		require.Equal(t, eps.EphemeralContainers, ps.EphemeralContainers)
	})
	t.Run("RestartPolicy copied", func(t *testing.T) {
		require.Equal(t, eps.RestartPolicy, ps.RestartPolicy)
	})
	t.Run("TerminationGracePeriodSeconds copied", func(t *testing.T) {
		require.NotNil(t, ps.TerminationGracePeriodSeconds)
		require.Equal(t, *eps.TerminationGracePeriodSeconds, *ps.TerminationGracePeriodSeconds)
	})
	t.Run("ServiceAccountName copied", func(t *testing.T) {
		require.Equal(t, eps.ServiceAccountName, ps.ServiceAccountName)
	})
}
