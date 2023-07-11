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

//go:build kubeall || helm
// +build kubeall helm

// Fire up a local Kubernetes cluster (`kind create cluster --config test/kind.yaml`)
// and run the acceptance tests against it (`go test -v -tags kubeall ./test`)

package test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
)

var (
	vaultVersion      = "latest"
	bankVaultsVersion = "latest"
	operatorVersion   = "latest"
)

// Installing the operator helm chart before testing
func TestMain(m *testing.M) {
	t := &testing.T{}

	// Setup Vault operator as a dependency for each test
	releaseName := "vault-operator"
	defaultKubectlOptions := k8s.NewKubectlOptions("", "", "default")

	// Set Vault version
	if v := os.Getenv("VAULT_VERSION"); v != "" {
		vaultVersion = v
	}

	// Set Bank vaults version
	if v := os.Getenv("BANK_VAULTS_VERSION"); v != "" {
		bankVaultsVersion = v
	}

	// Set Operator version
	if v := os.Getenv("OPERATOR_VERSION"); v != "" {
		operatorVersion = v
	}

	// Set Helm chart
	chart := "../deploy/charts/vault-operator"
	if v := os.Getenv("HELM_CHART"); v != "" {
		chart = v
	}

	// Setup args for helm.
	helmOptions := &helm.Options{
		KubectlOptions: defaultKubectlOptions,
		SetValues: map[string]string{
			"image.tag":           operatorVersion,
			"image.bankVaultsTag": bankVaultsVersion,
			"image.pullPolicy":    "Never",
		},
	}

	// Deploy the chart using `helm install` and wait until the pod comes up
	helm.Install(t, helmOptions, chart, releaseName)
	defer helm.Delete(t, helmOptions, releaseName, true)

	operatorPods := waitUntilPodsCreated(t, defaultKubectlOptions, releaseName, 10, 5*time.Second)
	k8s.WaitUntilPodAvailable(t, defaultKubectlOptions, operatorPods[0].GetName(), 5, 10*time.Second)

	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, defaultKubectlOptions)
	_, err = clientset.RbacV1().ClusterRoleBindings().Create(
		context.Background(),
		&v1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "vault-auth-delegator",
			},
			Subjects: []v1.Subject{},
			RoleRef: v1.RoleRef{
				Kind: "ClusterRole",
				Name: "system:auth-delegator",
			},
		},
		metav1.CreateOptions{},
	)
	require.NoError(t, err)
	defer clientset.RbacV1().ClusterRoleBindings().Delete(context.Background(), "vault-auth-delegator", metav1.DeleteOptions{})

	// Run tests
	exitCode := m.Run()

	// Exit based on the test results
	os.Exit(exitCode)
}

func TestKvv2(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"kvv2",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-kvv2.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)
}

func TestStatsd(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"statsd",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-statsd.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)
}

func TestExternalSecretsWatcherDeployment(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"external-secrets-watcher-deployment",
		[]string{
			"rbac.yaml",
			"deploy/test-external-secrets-watch-deployment.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)

	// Check pod annotation
	require.Equal(t, "", k8s.GetPod(t, kubectlOptions, "vault-0").GetAnnotations()["vault.banzaicloud.io/watched-secrets-sum"])
}

func TestExternalSecretsWatcherSecrets(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"external-secrets-watcher-secrets",
		[]string{
			"rbac.yaml",
			"deploy/test-external-secrets-watch-secrets.yaml",
			"deploy/test-external-secrets-watch-deployment.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)

	// Check pod annotation
	require.Equal(
		t,
		"bac8dfa8bdf03009f89303c8eb4a6c8f2fd80eb03fa658f53d6d65eec14666d4",
		k8s.GetPod(t, kubectlOptions, "vault-0").GetAnnotations()["vault.banzaicloud.io/watched-secrets-sum"],
	)
}

func TestRaft(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"raft",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-raft.yaml",
		},
	)

	// Wait until all vault pods come up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-1", 60, 10*time.Second)
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-2", 60, 10*time.Second)
}

func TestSoftHSM(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"softhsm",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-hsm-softhsm.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)
}

func TestDisabledRootTokenStorage(t *testing.T) {
	// t.Parallel()

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"disabled-root-token-storage",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-disabled-root-token-storage.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)

	// Check that the vault-root secret is not created
	_, err := k8s.GetSecretE(t, kubectlOptions, "vault-root")
	require.Errorf(t, err, `secrets "vault-root" not found`)
}

func TestPriorityClass(t *testing.T) {
	// t.Parallel()

	// TODO: Disable test for now until examples are fixed
	return

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"priority-class",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-priority.yaml",
		},
	)

	// Add ServiceAccount to ClusterRoleBinding
	clientset, err := k8s.GetKubernetesClientFromOptionsE(t, kubectlOptions)
	crb, err := clientset.RbacV1().ClusterRoleBindings().Get(context.Background(), "vault-auth-delegator", metav1.GetOptions{})
	crb.Subjects = append(crb.Subjects, v1.Subject{
		Kind:      "ServiceAccount",
		Name:      "vault",
		Namespace: kubectlOptions.Namespace,
	})
	_, err = clientset.RbacV1().ClusterRoleBindings().Update(context.Background(), crb, metav1.UpdateOptions{})
	require.NoError(t, err)

	// Wait until vault-0 pod comes up healthy and secrets are populated
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)
	time.Sleep(10 * time.Second)

	// Run an internal client in the default namespace which tries to read from Vault with the configured Kubernetes auth backend
	path, err := filepath.Abs("../cmd/examples/main.go")
	require.NoError(t, err)
	command := fmt.Sprintf("kurun run %s --env VAULT_ADDR=https://vault.%s:8200 --namespace %s", path, kubectlOptions.Namespace, kubectlOptions.Namespace)
	stdout, stderr, err := executeShellCommand(command)
	t.Logf("kurun run stdout: %s, stderr: %s, err: %v", stdout, stderr, err)
	require.NoError(t, err)
}

func TestOIDC(t *testing.T) {
	// t.Parallel()

	// TODO: Disable test for now until examples are fixed
	return

	// Prepare and apply resources
	kubectlOptions := prepareEnv(t,
		"default",
		[]string{
			"rbac.yaml",
			"../deploy/examples/cr-oidc.yaml",
		},
	)

	// Wait until vault-0 pod comes up healthy
	k8s.WaitUntilPodAvailable(t, kubectlOptions, "vault-0", 60, 10*time.Second)

	// Create a pod in the default namespace that uses OIDC authentication
	oidcPodFilePath, _ := filepath.Abs("oidc-pod.yaml")
	command := fmt.Sprintf("kurun apply -f %s -v", oidcPodFilePath)
	stdout, stderr, err := executeShellCommand(command)
	t.Logf("kurun apply stdout: %s, stderr: %s, err: %v", stdout, stderr, err)
	require.NoError(t, err)
	waitUntilPodSucceeded(t, kubectlOptions, "oidc", 60, 10*time.Second)

	// Clean up
	k8s.KubectlDelete(t, kubectlOptions, "../deploy/examples/cr-oidc.yaml")
	k8s.RunKubectl(t, kubectlOptions, "delete", "secret", "vault-unseal-keys")
	k8s.KubectlDelete(t, kubectlOptions, oidcPodFilePath)
}

// Installs k8s and kustomize components from configuration
func prepareEnv(t *testing.T, testName string, k8sRes []string) *k8s.KubectlOptions {
	// Setup a unique namespace for the resources for this test.
	namespaceName := fmt.Sprintf("test-%s-vault-%s", testName, strings.ReplaceAll(vaultVersion, ".", "-"))

	// Setup the kubectl config (default HOME/.kube/config) and the default current context.
	kubectlOptions := k8s.NewKubectlOptions("", "", namespaceName)

	// Prepare namespace
	prepareNamespace(t, namespaceName, kubectlOptions)

	// Reading files into byte slices
	var files [][]byte
	for _, crd := range k8sRes {
		data, err := os.ReadFile(crd)
		if err != nil {
			t.Fatal(err)
		}
		files = append(files, data)
	}

	// Decode byte slices into individual yaml documents
	var documents []interface{}
	for _, file := range files {
		dec := yaml.NewDecoder(bytes.NewReader(file))
		// Slice of the individual resources found in the file
		for {
			var v interface{}
			err := dec.Decode(&v)
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatal(err)
			}
			documents = append(documents, v)
		}
	}

	// Iterate on yaml documents and change namespace name where necessary
	for _, v := range documents {
		if v.(map[string]interface{})["kind"] == "Vault" {
			if i, ok := v.(map[string]interface{})["spec"].(map[string]interface{}); ok {
				if i["image"] != "" {
					i["image"] = "vault:" + vaultVersion
				}
			}

			if s, ok := v.(map[string]interface{})["spec"].(map[string]interface{})["unsealConfig"].(map[string]interface{})["kubernetes"].(map[string]interface{}); ok {
				if s["secretNamespace"] != "" {
					s["secretNamespace"] = namespaceName
				}
			}

			apiAddress := fmt.Sprintf("http://vault.%s:8200", namespaceName)
			if a, ok := v.(map[string]interface{})["spec"].(map[string]interface{})["config"].(map[string]interface{}); ok {
				if a["api_addr"] != "" {
					a["api_addr"] = apiAddress
				}
			}

			if b, ok := v.(map[string]interface{})["spec"].(map[string]interface{})["externalConfig"].(map[string]interface{})["auth"].([]interface{})[0].(map[string]interface{})["roles"].([]interface{})[0].(map[string]interface{})["bound_service_account_namespaces"].([]interface{}); ok {
				if b[0] != "" {
					b[0] = namespaceName
				}
			}
		}

		vaultAddress := fmt.Sprintf("https://vault.%s.svc.cluster.local:8200", namespaceName)
		if v.(map[string]interface{})["kind"] == "Secret" || v.(map[string]interface{})["kind"] == "ConfigMap" {
			if a, ok := v.(map[string]interface{})["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}); ok {
				a["vault.security.banzaicloud.io/vault-addr"] = vaultAddress
			}
		}

		resource, err := yaml.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}

		k8s.KubectlApplyFromString(t, kubectlOptions, string(resource))
	}

	return kubectlOptions
}

func prepareNamespace(t *testing.T, namespaceName string, kubectlOptions *k8s.KubectlOptions) {
	k8s.CreateNamespace(t, kubectlOptions, namespaceName)
	t.Cleanup(func() {
		k8s.DeleteNamespace(t, kubectlOptions, kubectlOptions.Namespace)
	})
}

func executeShellCommand(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func waitUntilPodsCreated(t *testing.T, options *k8s.KubectlOptions, deploymentName string, retries int, sleepBetweenRetries time.Duration) []corev1.Pod {
	statusMsg := fmt.Sprintf("Wait for Pod(s) %s to be created.", deploymentName)
	podsInterface, err := retry.DoWithRetryInterfaceE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (interface{}, error) {
			pods := k8s.ListPods(t, options, metav1.ListOptions{LabelSelector: labels.Set(map[string]string{"app.kubernetes.io/name": deploymentName}).String()})
			if len(pods) == 0 {
				return nil, errors.New("Pod(s) not created yet")
			}
			return pods, nil
		},
	)
	if err != nil {
		logger.Logf(t, "Timedout waiting for Pod(s) to be created: %s", err)
		require.NoError(t, err)
	}
	logger.Logf(t, "Pod(s) created")

	var createdPods []corev1.Pod
	if pods, ok := podsInterface.([]corev1.Pod); ok {
		createdPods = pods
	}

	return createdPods
}

func waitUntilPodSucceeded(t *testing.T, options *k8s.KubectlOptions, podName string, retries int, sleepBetweenRetries time.Duration) {
	statusMsg := fmt.Sprintf("Wait for Pod %s to succeed.", podName)
	message, err := retry.DoWithRetryE(
		t,
		statusMsg,
		retries,
		sleepBetweenRetries,
		func() (string, error) {
			pod, err := k8s.GetPodE(t, options, podName)
			if err != nil {
				return "", err
			}
			if string(pod.Status.Phase) != "Succeeded" {
				return "", errors.New("Pod is not succeeded yet")
			}
			return "Pod is now succeeded", nil
		},
	)
	if err != nil {
		logger.Logf(t, "Timedout waiting for Pod to succeed: %s", err)
		require.NoError(t, err)
	}
	logger.Logf(t, message)
}
