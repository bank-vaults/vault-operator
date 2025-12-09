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
	"github.com/stretchr/testify/assert"
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
	v := &vaultv1alpha1.Vault{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-vault",
			Namespace: "default",
		},
		Spec: vaultv1alpha1.VaultSpec{
			Config: extv1beta1.JSON{
				Raw: []byte(`{"listener": {"tcp": {"address": "127.0.0.1:8200", "tls_disable": 1}}, "storage": {"file": {"path": "/vault/file"}}}`),
			},
			VaultConfigurerPodSpec: &vaultv1alpha1.EmbeddedPodSpec{
				Containers: []corev1.Container{
					{
						Name: "bank-vaults",
						SecurityContext: &corev1.SecurityContext{
							RunAsUser:  utils.To(int64(1000)),
							Privileged: utils.To(false),
						},
					},
				},
			},
		},
	}

	deployment, err := deploymentForConfigurer(v, corev1.ConfigMapList{}, corev1.SecretList{}, map[string]string{})
	assert.NoError(t, err, "Failed to create deployment for configurer")
	assert.NotNil(t, deployment, "Deployment should not be nil")

	containers := deployment.Spec.Template.Spec.Containers
	assert.Greater(t, len(containers), 0, "Should have at least one container")

	var foundContainer bool
	for i := range containers {
		c := &containers[i]
		if c.Name == "bank-vaults" &&
			c.SecurityContext != nil &&
			c.SecurityContext.RunAsUser != nil &&
			*c.SecurityContext.RunAsUser == *utils.To(int64(1000)) &&
			c.SecurityContext.Privileged != nil &&
			*c.SecurityContext.Privileged == *utils.To(false) {
			foundContainer = true
			break
		}
	}
	assert.True(t, foundContainer, "Should find bank-vaults container with specified SecurityContext (RunAsUser=1000, Privileged=false)")
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
