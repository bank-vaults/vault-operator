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

package mutation

import (
	"encoding/json"

	"github.com/wI2L/jsondiff"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	vaultv1alpha1 "github.com/bank-vaults/vault-operator/pkg/apis/vault/v1alpha1"
)

// MutateVaultPatch returns a json patch containing all the mutations needed for
// a given Vault
func MutateVaultPatch(original *vaultv1alpha1.Vault) ([]byte, error) {
	modified := original.DeepCopy()

	// Skip if Affinity defined
	if modified.Spec.Affinity != nil {
		return []byte{}, nil
	}

	// Apply mutations
	if modified.Spec.PodAntiAffinity != "" {
		var podAntiAffinity *corev1.PodAntiAffinity
		err := json.Unmarshal([]byte(modified.Spec.PodAntiAffinity), &podAntiAffinity)
		if err != nil {
			return nil, err
		}
		modified.Spec.Affinity.PodAntiAffinity = &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: modified.LabelsForVault(),
					},
					TopologyKey: podAntiAffinity.String(),
				},
			},
		}
	}
	if modified.Spec.NodeAffinity.Size() != 0 {
		modified.Spec.Affinity.NodeAffinity = modified.Spec.NodeAffinity.DeepCopy()
	}

	// Generate patch
	patch, err := jsondiff.Compare(original, modified)
	if err != nil {
		return nil, err
	}
	return json.Marshal(patch)
}
