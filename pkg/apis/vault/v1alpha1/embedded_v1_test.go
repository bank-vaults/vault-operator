package v1alpha1

import (
	"testing"

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
		TerminationGracePeriodSeconds: int64Ptr(42),
		ActiveDeadlineSeconds:         int64Ptr(100),
		DNSPolicy:                     v1.DNSClusterFirst,
		NodeSelector:                  map[string]string{"k": "v"},
		ServiceAccountName:            "svc",
		DeprecatedServiceAccount:      "svc-deprecated",
		AutomountServiceAccountToken:  boolPtr(true),
		NodeName:                      "node",
		HostNetwork:                   true,
		HostPID:                       true,
		HostIPC:                       true,
		ShareProcessNamespace:         boolPtr(true),
		SecurityContext:               &v1.PodSecurityContext{RunAsUser: int64Ptr(1000)},
		ImagePullSecrets:              []v1.LocalObjectReference{{Name: "pullsecret"}},
		Hostname:                      "host",
		Subdomain:                     "sub",
		Affinity:                      &v1.Affinity{},
		SchedulerName:                 "sched",
		Tolerations:                   []v1.Toleration{{Key: "key"}},
		HostAliases:                   []v1.HostAlias{{IP: "127.0.0.1"}},
		PriorityClassName:             "high",
		Priority:                      int32Ptr(10),
		DNSConfig:                     &v1.PodDNSConfig{},
		ReadinessGates:                []v1.PodReadinessGate{{ConditionType: "ready"}},
		RuntimeClassName:              strPtr("runtime"),
		EnableServiceLinks:            boolPtr(false),
		PreemptionPolicy:              preemptionPtr(v1.PreemptLowerPriority),
		Overhead:                      func() v1.ResourceList { qv, _ := resource.ParseQuantity("1"); return v1.ResourceList{"cpu": qv} }(),
		TopologySpreadConstraints:     []v1.TopologySpreadConstraint{{MaxSkew: 1}},
		SetHostnameAsFQDN:             boolPtr(true),
		OS:                            &v1.PodOS{Name: "linux"},
		HostUsers:                     boolPtr(false),
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

func int64Ptr(i int64) *int64                                  { return &i }
func int32Ptr(i int32) *int32                                  { return &i }
func boolPtr(b bool) *bool                                     { return &b }
func strPtr(s string) *string                                  { return &s }
func preemptionPtr(p v1.PreemptionPolicy) *v1.PreemptionPolicy { return &p }

// Removed resourceQty function
