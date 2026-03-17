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

package vault

import (
	"context"
	"net/http"
	"testing"

	vaultv1alpha1 "github.com/bank-vaults/vault-operator/pkg/apis/vault/v1alpha1"
	"github.com/bank-vaults/vault-operator/pkg/utils"
	"github.com/siliconbrain/go-seqs/seqs"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var envs = []corev1.EnvVar{
	{
		Name:  "VAULT_NAMESPACE",
		Value: "default",
	},
	{
		Name:  "VAULT_IGNORE_MISSING_SECRETS",
		Value: "true",
	},
}

func TestFluentDConfFile(t *testing.T) {
	testFilename := "test.conf"

	v := &vaultv1alpha1.Vault{
		Spec: vaultv1alpha1.VaultSpec{
			FluentDConfFile: testFilename,
		},
	}

	configMap := configMapForFluentD(v)
	if configMap == nil || configMap.Data == nil {
		t.Errorf("configmap is nil or configmap data is nil")
	}

	if _, ok := configMap.Data[testFilename]; !ok {
		t.Errorf("configmap did not contain a key matching %q", testFilename)
		t.Logf("configmap: %+v", configMap)
	}
}

func TestFluentDConfFileDefault(t *testing.T) {
	defaultFilename := "fluent.conf"

	v := &vaultv1alpha1.Vault{
		Spec: vaultv1alpha1.VaultSpec{},
	}

	configMap := configMapForFluentD(v)

	if configMap == nil || configMap.Data == nil {
		t.Errorf("configmap is nil or configmap data is nil")
	}

	if _, ok := configMap.Data[defaultFilename]; !ok {
		t.Errorf("configmap did not contain a key matching %q", defaultFilename)
		t.Logf("configmap: %+v", configMap)
	}
}

func TestHandleStorageConfiguration_MissingStorage(t *testing.T) {
	// Vault object with missing storage configuration
	vault := &vaultv1alpha1.Vault{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-vault",
			Namespace: "default",
		},
		Spec: vaultv1alpha1.VaultSpec{
			Config: extv1beta1.JSON{
				Raw: []byte(`{"listener": {"tcp": {"address": "127.0.0.1:8200", "tls_disable": 1}}, "storage": {}}`),
			},
		},
	}

	// ReconcileVault instance with a fake client and scheme
	scheme := runtime.NewScheme()
	err := vaultv1alpha1.AddToScheme(scheme)
	client := fake.NewClientBuilder().WithScheme(scheme).Build()
	assert.NoError(t, err, "Failed to add Vault custom resource to scheme")

	reconciler := &ReconcileVault{
		client:              client,
		nonNamespacedClient: client,
		scheme:              client.Scheme(),
		httpClient:          &http.Client{},
	}

	err = reconciler.handleStorageConfiguration(context.Background(), vault)
	assert.Error(t, err, "Expected an error")
}

func TestVaultConfigurerPodSpecContainerMerge(t *testing.T) {
	baseVaultConfig := []byte(`{"listener": {"tcp": {"address": "127.0.0.1:8200", "tls_disable": 1}}, "storage": {"file": {"path": "/vault/file"}}}`)

	tests := []struct {
		name                   string
		vaultConfigurerPodSpec *vaultv1alpha1.EmbeddedPodSpec
		expectedContainerCount int
		validate               func(t *testing.T, deployment *appsv1.Deployment)
	}{
		{
			name:                   "no VaultConfigurerPodSpec - default deployment",
			vaultConfigurerPodSpec: nil,
			expectedContainerCount: 1,
			validate: func(t *testing.T, deployment *appsv1.Deployment) {
				container := &deployment.Spec.Template.Spec.Containers[0]
				assert.Equal(t, "bank-vaults", container.Name)
				assert.NotEmpty(t, container.Image, "Image should be set")
				assert.Equal(t, []string{"bank-vaults", "configure"}, container.Command)
				assert.NotEmpty(t, container.Ports, "Ports should be set")
				assert.Equal(t, int32(9091), container.Ports[0].ContainerPort)
				assert.NotEmpty(t, container.Env, "Env should be set")
				assert.Equal(t, "/config", container.WorkingDir)
			},
		},
		{
			name: "override existing bank-vaults container fields",
			vaultConfigurerPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						SecurityContext: &corev1.SecurityContext{
							RunAsUser:  utils.To(int64(1000)),
							Privileged: utils.To(false),
						},
						Env: []corev1.EnvVar{
							{
								Name:  "AZURE_CLIENT_ID",
								Value: "test-azure-client-id",
							},
						},
					},
				},
			},
			expectedContainerCount: 1,
			validate: func(t *testing.T, deployment *appsv1.Deployment) {
				container := &deployment.Spec.Template.Spec.Containers[0]
				assert.Equal(t, "bank-vaults", container.Name)
				assert.NotNil(t, container.SecurityContext)
				assert.Equal(t, int64(1000), *container.SecurityContext.RunAsUser)
				assert.Equal(t, false, *container.SecurityContext.Privileged)
				env, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "AZURE_CLIENT_ID" }))
				assert.True(t, found, "AZURE_CLIENT_ID env var should exist")
				assert.Equal(t, "test-azure-client-id", env.Value)
				assert.NotEmpty(t, container.Image, "Image should still be set")
				assert.Equal(t, []string{"bank-vaults", "configure"}, container.Command)
				assert.Equal(t, "/config", container.WorkingDir)
			},
		},
		{
			name: "add additional sidecar container",
			vaultConfigurerPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name:    "sidecar",
						Image:   "busybox:latest",
						Command: []string{"sleep", "infinity"},
					},
				},
			},
			expectedContainerCount: 2,
			validate: func(t *testing.T, deployment *appsv1.Deployment) {
				containers := deployment.Spec.Template.Spec.Containers
				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should exist")
				assert.NotEmpty(t, bankVaults.Image)
				assert.Equal(t, []string{"bank-vaults", "configure"}, bankVaults.Command)
				sidecar, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "sidecar" }))
				assert.True(t, found, "sidecar container should exist")
				assert.Equal(t, "busybox:latest", sidecar.Image)
				assert.Equal(t, []string{"sleep", "infinity"}, sidecar.Command)
			},
		},
		{
			name: "override bank-vaults and add sidecar",
			vaultConfigurerPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						Env:  []corev1.EnvVar{{Name: "CUSTOM_VAR", Value: "custom-value"}},
					},
					{
						Name:  "logger",
						Image: "fluentd:latest",
					},
				},
			},
			expectedContainerCount: 2,
			validate: func(t *testing.T, deployment *appsv1.Deployment) {
				containers := deployment.Spec.Template.Spec.Containers
				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found)
				env, found := seqs.First(seqs.Filter(seqs.FromSlice(bankVaults.Env), func(e corev1.EnvVar) bool { return e.Name == "CUSTOM_VAR" }))
				assert.True(t, found, "CUSTOM_VAR should exist")
				assert.Equal(t, "custom-value", env.Value)
				assert.NotEmpty(t, bankVaults.Image)
				assert.Equal(t, []string{"bank-vaults", "configure"}, bankVaults.Command)

				logger, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "logger" }))
				assert.True(t, found)
				assert.Equal(t, "fluentd:latest", logger.Image)
			},
		},
		{
			name: "empty containers slice - no changes",
			vaultConfigurerPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{},
			},
			expectedContainerCount: 1,
			validate: func(t *testing.T, deployment *appsv1.Deployment) {
				container := &deployment.Spec.Template.Spec.Containers[0]
				assert.Equal(t, "bank-vaults", container.Name)
				assert.NotEmpty(t, container.Image)
				assert.Equal(t, []string{"bank-vaults", "configure"}, container.Command)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &vaultv1alpha1.Vault{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-vault",
					Namespace: "default",
				},
				Spec: vaultv1alpha1.VaultSpec{
					Config:                 extv1beta1.JSON{Raw: baseVaultConfig},
					VaultConfigurerPodSpec: tt.vaultConfigurerPodSpec,
				},
			}

			deployment, err := deploymentForConfigurer(v, corev1.ConfigMapList{}, corev1.SecretList{}, map[string]string{})
			assert.NoError(t, err)
			assert.NotNil(t, deployment)
			assert.Len(t, deployment.Spec.Template.Spec.Containers, tt.expectedContainerCount)
			tt.validate(t, deployment)
		})
	}
}

func TestWithVaultEnv(t *testing.T) {
	tests := []struct {
		name     string
		vault    *vaultv1alpha1.Vault
		envs     []corev1.EnvVar
		expected []corev1.EnvVar
	}{
		{
			name: "vaultEnvsConfig specified",
			vault: &vaultv1alpha1.Vault{
				Spec: vaultv1alpha1.VaultSpec{
					VaultEnvsConfig: []corev1.EnvVar{
						{
							Name:  "VAULT_TOKEN",
							Value: "vault:login",
						},
						{
							Name:  "VAULT_ADDR",
							Value: "http://vault:8200",
						},
					},
				},
			},
			envs: envs,
			expected: []corev1.EnvVar{
				{
					Name:  "VAULT_NAMESPACE",
					Value: "default",
				},
				{
					Name:  "VAULT_IGNORE_MISSING_SECRETS",
					Value: "true",
				},
				{
					Name:  "VAULT_TOKEN",
					Value: "vault:login",
				},
				{
					Name:  "VAULT_ADDR",
					Value: "http://vault:8200",
				},
			},
		},
		{
			name: "secretInitEnvsConfig specified",
			vault: &vaultv1alpha1.Vault{
				Spec: vaultv1alpha1.VaultSpec{
					SecretInitsConfig: []corev1.EnvVar{

						{
							Name:  "SECRET_INIT_LOG_LEVEL",
							Value: "info",
						},
						{
							Name:  "SECRET_INIT_JSON_LOG",
							Value: "true",
						},
						{
							Name:  "SECRET_INIT_DAEMON",
							Value: "true",
						},
					},
				},
			},
			envs: envs,
			expected: []corev1.EnvVar{
				{
					Name:  "VAULT_NAMESPACE",
					Value: "default",
				},
				{
					Name:  "VAULT_IGNORE_MISSING_SECRETS",
					Value: "true",
				},
				{
					Name:  "SECRET_INIT_LOG_LEVEL",
					Value: "info",
				},
				{
					Name:  "SECRET_INIT_JSON_LOG",
					Value: "true",
				},
				{
					Name:  "SECRET_INIT_DAEMON",
					Value: "true",
				},
			},
		},
		{
			name: "VaultEnvsConfig specified with deprecated envs",
			vault: &vaultv1alpha1.Vault{
				Spec: vaultv1alpha1.VaultSpec{
					VaultEnvsConfig: []corev1.EnvVar{
						{
							Name:  "VAULT_TOKEN",
							Value: "vault:login",
						},
						{
							Name:  "VAULT_ADDR",
							Value: "http://vault:8200",
						},
						{
							Name:  "VAULT_JSON_LOG",
							Value: "true",
						},
						{
							Name:  "VAULT_ENV_LOG_SERVER",
							Value: "https://logserver:8200",
						},
						{
							Name:  "VAULT_ENV_DAEMON",
							Value: "true",
						},
						{
							Name:  "VAULT_ENV_DELAY",
							Value: "10",
						},
						{
							Name:  "VAULT_ENV_FROM_PATH",
							Value: "vault:secret/data/test",
						},
						{
							Name:  "VAULT_ENV_PASSTHROUGH",
							Value: "VAULT_TOKEN",
						},
					},
				},
			},
			envs: envs,
			expected: []corev1.EnvVar{
				{
					Name:  "VAULT_NAMESPACE",
					Value: "default",
				},
				{
					Name:  "VAULT_IGNORE_MISSING_SECRETS",
					Value: "true",
				},
				{
					Name:  "VAULT_TOKEN",
					Value: "vault:login",
				},
				{
					Name:  "VAULT_ADDR",
					Value: "http://vault:8200",
				},
				{
					Name:  "SECRET_INIT_JSON_LOG",
					Value: "true",
				},
				{
					Name:  "SECRET_INIT_LOG_SERVER",
					Value: "https://logserver:8200",
				},
				{
					Name:  "SECRET_INIT_DAEMON",
					Value: "true",
				},
				{
					Name:  "SECRET_INIT_DELAY",
					Value: "10",
				},
				{
					Name:  "VAULT_FROM_PATH",
					Value: "vault:secret/data/test",
				},
				{
					Name:  "VAULT_PASSTHROUGH",
					Value: "VAULT_TOKEN",
				},
			},
		},
	}

	for _, tt := range tests {
		ttp := tt
		t.Run(tt.name, func(t *testing.T) {
			envs := withVaultEnv(ttp.vault, ttp.envs)
			assert.Equal(t, ttp.expected, envs, "envs did not match")
		})
	}
}

func TestMergeContainersByName(t *testing.T) {
	tests := []struct {
		name     string
		src      []corev1.Container
		dst      []corev1.Container
		validate func(t *testing.T, result []corev1.Container, err error)
	}{
		{
			name: "env vars are appended not replaced",
			src: []corev1.Container{
				{
					Name: "bank-vaults",
					Env: []corev1.EnvVar{
						{Name: "AZURE_CLIENT_ID", Value: "azure-client-id"},
						{Name: "AZURE_TENANT_ID", Value: "azure-tenant-id"},
					},
				},
			},
			dst: []corev1.Container{
				{
					Name:    "bank-vaults",
					Image:   "ghcr.io/bank-vaults/bank-vaults:latest",
					Command: []string{"bank-vaults", "configure"},
					Env: []corev1.EnvVar{
						{Name: "VAULT_ADDR", Value: "https://vault:8200"},
						{Name: "VAULT_CACERT", Value: "/vault/tls/ca.crt"},
						{Name: "NAMESPACE", Value: "vault"},
					},
				},
			},
			validate: func(t *testing.T, result []corev1.Container, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				container := result[0]

				vaultAddr, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "VAULT_ADDR" }))
				assert.True(t, found, "VAULT_ADDR should be preserved")
				assert.Equal(t, "https://vault:8200", vaultAddr.Value)

				vaultCacert, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "VAULT_CACERT" }))
				assert.True(t, found, "VAULT_CACERT should be preserved")
				assert.Equal(t, "/vault/tls/ca.crt", vaultCacert.Value)

				namespace, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "NAMESPACE" }))
				assert.True(t, found, "NAMESPACE should be preserved")
				assert.Equal(t, "vault", namespace.Value)

				azureClientID, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "AZURE_CLIENT_ID" }))
				assert.True(t, found, "AZURE_CLIENT_ID should be appended")
				assert.Equal(t, "azure-client-id", azureClientID.Value)

				azureTenantID, found := seqs.First(seqs.Filter(seqs.FromSlice(container.Env), func(e corev1.EnvVar) bool { return e.Name == "AZURE_TENANT_ID" }))
				assert.True(t, found, "AZURE_TENANT_ID should be appended")
				assert.Equal(t, "azure-tenant-id", azureTenantID.Value)

				assert.Len(t, container.Env, 5)
				assert.Equal(t, "ghcr.io/bank-vaults/bank-vaults:latest", container.Image)
				assert.Equal(t, []string{"bank-vaults", "configure"}, container.Command)
			},
		},
		{
			name: "volume mounts are appended not replaced",
			src: []corev1.Container{
				{
					Name: "bank-vaults",
					VolumeMounts: []corev1.VolumeMount{
						{Name: "azure-token", MountPath: "/var/run/secrets/azure"},
					},
				},
			},
			dst: []corev1.Container{
				{
					Name:  "bank-vaults",
					Image: "ghcr.io/bank-vaults/bank-vaults:latest",
					VolumeMounts: []corev1.VolumeMount{
						{Name: "vault-tls", MountPath: "/vault/tls"},
						{Name: "vault-config", MountPath: "/config"},
					},
				},
			},
			validate: func(t *testing.T, result []corev1.Container, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				container := result[0]
				assert.Len(t, container.VolumeMounts, 3)

				vaultTLS, found := seqs.First(seqs.Filter(seqs.FromSlice(container.VolumeMounts), func(v corev1.VolumeMount) bool { return v.Name == "vault-tls" }))
				assert.True(t, found, "vault-tls mount should be preserved")
				assert.Equal(t, "/vault/tls", vaultTLS.MountPath)

				vaultConfig, found := seqs.First(seqs.Filter(seqs.FromSlice(container.VolumeMounts), func(v corev1.VolumeMount) bool { return v.Name == "vault-config" }))
				assert.True(t, found, "vault-config mount should be preserved")
				assert.Equal(t, "/config", vaultConfig.MountPath)

				azureToken, found := seqs.First(seqs.Filter(seqs.FromSlice(container.VolumeMounts), func(v corev1.VolumeMount) bool { return v.Name == "azure-token" }))
				assert.True(t, found, "azure-token mount should be appended")
				assert.Equal(t, "/var/run/secrets/azure", azureToken.MountPath)
			},
		},
		{
			name: "non-slice fields are overridden",
			src: []corev1.Container{
				{
					Name:       "bank-vaults",
					Image:      "custom-image:v1.0.0",
					WorkingDir: "/custom/workdir",
				},
			},
			dst: []corev1.Container{
				{
					Name:       "bank-vaults",
					Image:      "ghcr.io/bank-vaults/bank-vaults:latest",
					WorkingDir: "/config",
					Command:    []string{"bank-vaults", "configure"},
				},
			},
			validate: func(t *testing.T, result []corev1.Container, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				container := result[0]

				assert.Equal(t, "custom-image:v1.0.0", container.Image)
				assert.Equal(t, "/custom/workdir", container.WorkingDir)
				assert.Equal(t, []string{"bank-vaults", "configure"}, container.Command)
			},
		},
		{
			name: "new container from src is appended",
			src: []corev1.Container{
				{
					Name:  "sidecar",
					Image: "busybox:latest",
				},
			},
			dst: []corev1.Container{
				{
					Name:  "bank-vaults",
					Image: "ghcr.io/bank-vaults/bank-vaults:latest",
				},
			},
			validate: func(t *testing.T, result []corev1.Container, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 2)

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(result), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found)
				assert.Equal(t, "ghcr.io/bank-vaults/bank-vaults:latest", bankVaults.Image)

				sidecar, found := seqs.First(seqs.Filter(seqs.FromSlice(result), func(c corev1.Container) bool { return c.Name == "sidecar" }))
				assert.True(t, found)
				assert.Equal(t, "busybox:latest", sidecar.Image)
			},
		},
		{
			name: "empty src returns dst unchanged",
			src:  []corev1.Container{},
			dst: []corev1.Container{
				{
					Name:  "bank-vaults",
					Image: "ghcr.io/bank-vaults/bank-vaults:latest",
					Env:   []corev1.EnvVar{{Name: "VAULT_ADDR", Value: "https://vault:8200"}},
				},
			},
			validate: func(t *testing.T, result []corev1.Container, err error) {
				assert.NoError(t, err)
				assert.Len(t, result, 1)
				assert.Equal(t, "bank-vaults", result[0].Name)
				assert.Len(t, result[0].Env, 1)
				assert.Equal(t, "VAULT_ADDR", result[0].Env[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mergeContainersByName(tt.src, tt.dst)
			tt.validate(t, result, err)
		})
	}
}

func TestVaultPodSpecContainerMerge(t *testing.T) {
	baseVaultConfig := []byte(`{"listener": {"tcp": {"address": "127.0.0.1:8200", "tls_disable": 1}}, "storage": {"file": {"path": "/vault/file"}}}`)
	service := &corev1.Service{
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	tests := []struct {
		name                       string
		vaultPodSpec               *vaultv1alpha1.EmbeddedPodSpec
		expectedContainerCount     int
		expectedInitContainerCount int
		validate                   func(t *testing.T, sts *appsv1.StatefulSet)
	}{
		{
			// vault + bank-vaults + prometheus-exporter (statsd enabled by default)
			name:                       "no VaultPodSpec - default statefulset",
			vaultPodSpec:               nil,
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers

				vault, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "vault" }))
				assert.True(t, found, "vault container should exist")
				assert.NotEmpty(t, vault.Image, "Image should be set")
				assert.Equal(t, []string{"server"}, vault.Args)
				assert.NotEmpty(t, vault.Env, "Env should be set")

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should exist")
				assert.NotEmpty(t, bankVaults.Image, "Image should be set")
				assert.NotEmpty(t, bankVaults.Command, "Command should be set")
				assert.NotEmpty(t, bankVaults.Env, "Env should be set")

				initContainers := sts.Spec.Template.Spec.InitContainers
				configTemplating, found := seqs.First(seqs.Filter(seqs.FromSlice(initContainers), func(c corev1.Container) bool { return c.Name == "config-templating" }))
				assert.True(t, found, "config-templating init container should exist")
				assert.NotEmpty(t, configTemplating.Image, "Image should be set")
				assert.NotEmpty(t, configTemplating.Command, "Command should be set")
			},
		},
		{
			name: "override bank-vaults sidecar security context",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: utils.To(false),
							ReadOnlyRootFilesystem:   utils.To(true),
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{"ALL"},
							},
						},
					},
				},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should exist")
				assert.NotNil(t, bankVaults.SecurityContext, "security context should be set")
				assert.Equal(t, false, *bankVaults.SecurityContext.AllowPrivilegeEscalation)
				assert.Equal(t, true, *bankVaults.SecurityContext.ReadOnlyRootFilesystem)
				assert.Equal(t, []corev1.Capability{"ALL"}, bankVaults.SecurityContext.Capabilities.Drop)

				// vault container should still exist and be unmodified
				vault, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "vault" }))
				assert.True(t, found, "vault container should exist")
				assert.NotEmpty(t, vault.Image, "vault image should be set")

				// bank-vaults should preserve operator-set fields
				assert.NotEmpty(t, bankVaults.Image, "bank-vaults image should be preserved")
				assert.NotEmpty(t, bankVaults.Command, "bank-vaults command should be preserved")
			},
		},
		{
			name: "override config-templating init container security context",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				InitContainers: []corev1.Container{
					{
						Name: "config-templating",
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: utils.To(false),
							ReadOnlyRootFilesystem:   utils.To(true),
							Capabilities: &corev1.Capabilities{
								Drop: []corev1.Capability{"ALL"},
							},
						},
					},
				},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				initContainers := sts.Spec.Template.Spec.InitContainers

				configTemplating, found := seqs.First(seqs.Filter(seqs.FromSlice(initContainers), func(c corev1.Container) bool { return c.Name == "config-templating" }))
				assert.True(t, found, "config-templating init container should exist")
				assert.NotNil(t, configTemplating.SecurityContext, "security context should be set")
				assert.Equal(t, false, *configTemplating.SecurityContext.AllowPrivilegeEscalation)
				assert.Equal(t, true, *configTemplating.SecurityContext.ReadOnlyRootFilesystem)
				assert.Equal(t, []corev1.Capability{"ALL"}, configTemplating.SecurityContext.Capabilities.Drop)

				// config-templating should preserve operator-set fields
				assert.NotEmpty(t, configTemplating.Image, "image should be preserved")
				assert.NotEmpty(t, configTemplating.Command, "command should be preserved")
			},
		},
		{
			name: "override multiple containers and init containers simultaneously",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: utils.To(false),
						},
					},
				},
				InitContainers: []corev1.Container{
					{
						Name: "config-templating",
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: utils.To(false),
						},
					},
				},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers
				initContainers := sts.Spec.Template.Spec.InitContainers

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found)
				assert.NotNil(t, bankVaults.SecurityContext)
				assert.Equal(t, false, *bankVaults.SecurityContext.AllowPrivilegeEscalation)

				configTemplating, found := seqs.First(seqs.Filter(seqs.FromSlice(initContainers), func(c corev1.Container) bool { return c.Name == "config-templating" }))
				assert.True(t, found)
				assert.NotNil(t, configTemplating.SecurityContext)
				assert.Equal(t, false, *configTemplating.SecurityContext.AllowPrivilegeEscalation)

				// All original containers should still be present
				vault, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "vault" }))
				assert.True(t, found, "vault container should still exist")
				assert.NotEmpty(t, vault.Image, "vault image should be preserved")
			},
		},
		{
			name: "env vars for bank-vaults sidecar are appended",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						Env: []corev1.EnvVar{
							{Name: "CUSTOM_VAR", Value: "custom-value"},
						},
					},
				},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found)

				// Custom var should be appended
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(bankVaults.Env), func(e corev1.EnvVar) bool { return e.Name == "CUSTOM_VAR" }))
				assert.True(t, found, "CUSTOM_VAR should be appended")

				// Operator-set env vars should be preserved
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(bankVaults.Env), func(e corev1.EnvVar) bool { return e.Name == "POD_NAME" }))
				assert.True(t, found, "operator POD_NAME env var should be preserved")

				assert.NotEmpty(t, bankVaults.Image, "Image should be preserved")
				assert.NotEmpty(t, bankVaults.Command, "Command should be preserved")
			},
		},
		{
			name: "add additional sidecar container",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name:    "log-shipper",
						Image:   "fluentd:latest",
						Command: []string{"fluentd", "-c", "/etc/fluentd.conf"},
					},
				},
			},
			expectedContainerCount:     4,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers

				// New container should be appended
				logShipper, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "log-shipper" }))
				assert.True(t, found, "log-shipper container should be appended")
				assert.Equal(t, "fluentd:latest", logShipper.Image)
				assert.Equal(t, []string{"fluentd", "-c", "/etc/fluentd.conf"}, logShipper.Command)

				// Existing containers should be preserved
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "vault" }))
				assert.True(t, found, "vault container should be preserved")

				_, found = seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should be preserved")
			},
		},
		{
			name: "override bank-vaults and add sidecar",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						SecurityContext: &corev1.SecurityContext{
							AllowPrivilegeEscalation: utils.To(false),
						},
					},
					{
						Name:  "log-shipper",
						Image: "fluentd:latest",
					},
				},
			},
			expectedContainerCount:     4,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers

				bankVaults, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should exist")
				assert.NotNil(t, bankVaults.SecurityContext, "security context should be set")
				assert.Equal(t, false, *bankVaults.SecurityContext.AllowPrivilegeEscalation)
				assert.NotEmpty(t, bankVaults.Image, "bank-vaults image should be preserved")
				assert.NotEmpty(t, bankVaults.Command, "bank-vaults command should be preserved")

				logShipper, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "log-shipper" }))
				assert.True(t, found, "log-shipper container should be appended")
				assert.Equal(t, "fluentd:latest", logShipper.Image)
			},
		},
		{
			name: "empty containers slice - no changes",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers:     []corev1.Container{},
				InitContainers: []corev1.Container{},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				containers := sts.Spec.Template.Spec.Containers
				initContainers := sts.Spec.Template.Spec.InitContainers

				_, found := seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "vault" }))
				assert.True(t, found, "vault container should exist")

				_, found = seqs.First(seqs.Filter(seqs.FromSlice(containers), func(c corev1.Container) bool { return c.Name == "bank-vaults" }))
				assert.True(t, found, "bank-vaults container should exist")

				_, found = seqs.First(seqs.Filter(seqs.FromSlice(initContainers), func(c corev1.Container) bool { return c.Name == "config-templating" }))
				assert.True(t, found, "config-templating init container should exist")
			},
		},
		{
			name: "extra volumes from vaultPodSpec are appended to operator volumes",
			vaultPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "custom-volume",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					},
				},
			},
			expectedContainerCount:     3,
			expectedInitContainerCount: 1,
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				volumes := sts.Spec.Template.Spec.Volumes

				// Custom volume should be appended
				_, found := seqs.First(seqs.Filter(seqs.FromSlice(volumes), func(v corev1.Volume) bool { return v.Name == "custom-volume" }))
				assert.True(t, found, "custom volume should be appended")

				// Operator-set volumes should be preserved
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(volumes), func(v corev1.Volume) bool { return v.Name == "vault-config" }))
				assert.True(t, found, "operator vault-config volume should be preserved")

				_, found = seqs.First(seqs.Filter(seqs.FromSlice(volumes), func(v corev1.Volume) bool { return v.Name == "vault-raw-config" }))
				assert.True(t, found, "operator vault-raw-config volume should be preserved")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &vaultv1alpha1.Vault{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-vault",
					Namespace: "default",
				},
				Spec: vaultv1alpha1.VaultSpec{
					Size:         1,
					Config:       extv1beta1.JSON{Raw: baseVaultConfig},
					VaultPodSpec: tt.vaultPodSpec,
				},
			}

			sts, err := statefulSetForVault(v, []corev1.Secret{}, map[string]string{}, service)
			assert.NoError(t, err)
			assert.NotNil(t, sts)
			assert.Len(t, sts.Spec.Template.Spec.Containers, tt.expectedContainerCount)
			assert.Len(t, sts.Spec.Template.Spec.InitContainers, tt.expectedInitContainerCount)
			tt.validate(t, sts)
		})
	}
}

func TestVaultContainerSpecEnvAppend(t *testing.T) {
	baseVaultConfig := []byte(`{"listener": {"tcp": {"address": "127.0.0.1:8200", "tls_disable": 1}}, "storage": {"file": {"path": "/vault/file"}}}`)
	service := &corev1.Service{
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	tests := []struct {
		name               string
		vaultContainerSpec corev1.Container
		validate           func(t *testing.T, sts *appsv1.StatefulSet)
	}{
		{
			name: "env vars are appended not replaced",
			vaultContainerSpec: corev1.Container{
				Name: "vault",
				Env: []corev1.EnvVar{
					{Name: "CUSTOM_VAR", Value: "custom-value"},
				},
			},
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				vault := sts.Spec.Template.Spec.Containers[0]
				assert.Equal(t, "vault", vault.Name)

				// Custom var should be appended
				_, found := seqs.First(seqs.Filter(seqs.FromSlice(vault.Env), func(e corev1.EnvVar) bool { return e.Name == "CUSTOM_VAR" }))
				assert.True(t, found, "CUSTOM_VAR should be appended")

				// Operator-set env vars should be preserved
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(vault.Env), func(e corev1.EnvVar) bool { return e.Name == "VAULT_K8S_POD_NAME" }))
				assert.True(t, found, "operator VAULT_K8S_POD_NAME env var should be preserved")
			},
		},
		{
			name: "volume mounts are appended not replaced",
			vaultContainerSpec: corev1.Container{
				Name: "vault",
				VolumeMounts: []corev1.VolumeMount{
					{Name: "custom-mount", MountPath: "/custom"},
				},
			},
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				vault := sts.Spec.Template.Spec.Containers[0]

				// Custom mount should be appended
				_, found := seqs.First(seqs.Filter(seqs.FromSlice(vault.VolumeMounts), func(v corev1.VolumeMount) bool { return v.Name == "custom-mount" }))
				assert.True(t, found, "custom-mount should be appended")

				// Operator-set mounts should be preserved
				_, found = seqs.First(seqs.Filter(seqs.FromSlice(vault.VolumeMounts), func(v corev1.VolumeMount) bool { return v.Name == "vault-config" }))
				assert.True(t, found, "operator vault-config mount should be preserved")
			},
		},
		{
			name: "scalar fields still override",
			vaultContainerSpec: corev1.Container{
				Name:       "vault",
				WorkingDir: "/custom/workdir",
			},
			validate: func(t *testing.T, sts *appsv1.StatefulSet) {
				vault := sts.Spec.Template.Spec.Containers[0]
				assert.Equal(t, "/custom/workdir", vault.WorkingDir)

				// Operator-set fields should be preserved when not overridden
				assert.NotEmpty(t, vault.Image, "Image should be preserved")
				assert.Equal(t, []string{"server"}, vault.Args)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &vaultv1alpha1.Vault{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-vault",
					Namespace: "default",
				},
				Spec: vaultv1alpha1.VaultSpec{
					Size:               1,
					Config:             extv1beta1.JSON{Raw: baseVaultConfig},
					VaultContainerSpec: tt.vaultContainerSpec,
				},
			}

			sts, err := statefulSetForVault(v, []corev1.Secret{}, map[string]string{}, service)
			assert.NoError(t, err)
			assert.NotNil(t, sts)
			tt.validate(t, sts)
		})
	}
}
