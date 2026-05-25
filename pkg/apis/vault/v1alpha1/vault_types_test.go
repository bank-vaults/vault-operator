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
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
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

func TestGetConfigPath(t *testing.T) {
	t.Run("No config path specified", func(t *testing.T) {
		vault := &VaultSpec{}
		path := vault.GetConfigPath()
		require.Equal(t, "/vault/config", path)
	})
	t.Run("Config path specified", func(t *testing.T) {
		vault := &VaultSpec{
			ConfigPath: "/openbao/config",
		}
		path := vault.GetConfigPath()
		require.Equal(t, "/openbao/config", path)
	})
}

func TestGetAPIPort(t *testing.T) {
	mkSpec := func(addr string) *VaultSpec {
		cfg := []byte(`{}`)
		if addr != "" {
			cfg = []byte(`{"listener":{"tcp":{"address":"` + addr + `"}}}`)
		}
		return &VaultSpec{Config: extv1beta1.JSON{Raw: cfg}}
	}

	tests := []struct {
		name string
		addr string
		want int
	}{
		{name: "default when no config", addr: "", want: 8200},
		{name: "default IPv4 listener", addr: "0.0.0.0:8200", want: 8200},
		{name: "custom IPv4 port", addr: "0.0.0.0:9200", want: 9200},
		{name: "loopback custom port", addr: "127.0.0.1:8201", want: 8201},
		{name: "IPv6 listener", addr: "[::]:8200", want: 8200},
		{name: "malformed address falls back", addr: "not-a-real-address", want: 8200},
		{name: "non-numeric port falls back", addr: "0.0.0.0:vault", want: 8200},
		{name: "out-of-range port falls back", addr: "0.0.0.0:70000", want: 8200},
		{name: "negative port falls back", addr: "0.0.0.0:-1", want: 8200},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := mkSpec(tc.addr).GetAPIPort()
			require.Equal(t, tc.want, got)
		})
	}
}
