// Copyright © 2023 Bank-Vaults Maintainers
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

	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	t.Run("Good", func(t *testing.T) {
		tests := []string{
			"bank-vaults/my-vault:1.2.3",
			"bank-vaults/my-vault:1.2",
			"my.local.proxy/bank-vaults/my-vault:1.2.3",
			"my.local.proxy:5000/bank-vaults/my-vault:1.2.3",
			"bank-vaults/my-vault:v1.2.3",
			"bank-vaults/my-vault:v1.2",
			"my.local.proxy/bank-vaults/my-vault:v1.2.3",
			"my.local.proxy:5000/bank-vaults/my-vault:v1.2.3",
		}

		for _, tt := range tests {
			tt := tt

			t.Run("", func(t *testing.T) {
				vault := &VaultSpec{
					Image: tt,
				}

				_, err := vault.GetVersion()
				require.NoError(t, err)
			})
		}
	})

	t.Run("Bad", func(t *testing.T) {
		tests := []string{
			"bank-vaults/my-vault",
			"bank-vaults/my-vault:latest",
			"bank-vaults/my-vault:my-custom-build",
		}

		for _, tt := range tests {
			tt := tt

			t.Run("", func(t *testing.T) {
				vault := &VaultSpec{
					Image: tt,
				}

				_, err := vault.GetVersion()
				require.Error(t, err)
			})
		}
	})
}
