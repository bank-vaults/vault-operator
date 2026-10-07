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
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	bankvaultsdevv1alpha1 "github.com/bank-vaults/vault-operator/pkg/apis/vault/v1alpha1"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	// +kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

type updateCountingClient struct {
	client.Client
	updateCalls int
}

func (c *updateCountingClient) Update(ctx context.Context, object client.Object, options ...client.UpdateOption) error {
	c.updateCalls++
	return c.Client.Update(ctx, object, options...)
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	gatewayModuleDirectory, err := exec.Command("go", "list", "-f", "{{.Module.Dir}}", "sigs.k8s.io/gateway-api/apis/v1").Output()
	Expect(err).NotTo(HaveOccurred())
	gatewayCRDDirectory := filepath.Join(strings.TrimSpace(string(gatewayModuleDirectory)), "config", "crd", "standard")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "..", "deploy", "crd", "bases"),
			gatewayCRDDirectory,
		},
		ErrorIfCRDPathMissing: true,
	}

	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = bankvaultsdevv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	Expect(gatewayv1.Install(scheme.Scheme)).To(Succeed())

	// +kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

})

var _ = It("does not update an HTTPRoute after server defaults backend references", func() {
	ctx := context.Background()
	parentGroup := gatewayv1.Group("gateway.networking.k8s.io")
	parentKind := gatewayv1.Kind("Gateway")
	pathType := gatewayv1.PathMatchPathPrefix
	rootPath := "/"
	vault := &bankvaultsdevv1alpha1.Vault{
		Name:      "http-route-defaults-vault",
		Namespace: "default",
		UID:       "http-route-defaults-vault-uid",
		Spec: bankvaultsdevv1alpha1.VaultSpec{
			HTTPRoute: &bankvaultsdevv1alpha1.HTTPRoute{
				Spec: gatewayv1.HTTPRouteSpec{
					CommonRouteSpec: gatewayv1.CommonRouteSpec{
						ParentRefs: []gatewayv1.ParentReference{{
							Group: &parentGroup,
							Kind:  &parentKind,
							Name:  "test-gateway",
						}},
					},
					Rules: []gatewayv1.HTTPRouteRule{{
						Matches: []gatewayv1.HTTPRouteMatch{{
							Path: &gatewayv1.HTTPPathMatch{Type: &pathType, Value: &rootPath},
						}},
					}},
				},
			},
		},
	}

	countingClient := &updateCountingClient{Client: k8sClient}
	reconciler := &ReconcileVault{client: countingClient, apiReader: k8sClient, scheme: scheme.Scheme}
	Expect(reconciler.reconcileHTTPRoute(ctx, vault)).To(Succeed())
	Expect(countingClient.updateCalls).To(Equal(0))

	serverRoute := &gatewayv1.HTTPRoute{}
	Expect(k8sClient.Get(ctx, client.ObjectKeyFromObject(vault), serverRoute)).To(Succeed())
	serverBackendRef := serverRoute.Spec.Rules[0].BackendRefs[0]
	Expect(serverBackendRef.Group).NotTo(BeNil())
	Expect(*serverBackendRef.Group).To(Equal(gatewayv1.Group("")))
	Expect(serverBackendRef.Kind).NotTo(BeNil())
	Expect(*serverBackendRef.Kind).To(Equal(gatewayv1.Kind("Service")))
	Expect(serverBackendRef.Weight).NotTo(BeNil())
	Expect(*serverBackendRef.Weight).To(Equal(int32(1)))
	Expect(serverRoute.Spec).To(Equal(httpRouteForVault(vault).Spec))

	Expect(reconciler.reconcileHTTPRoute(ctx, vault)).To(Succeed())
	Expect(countingClient.updateCalls).To(Equal(0))
	DeferCleanup(func() {
		Expect(k8sClient.Delete(ctx, serverRoute)).To(Succeed())
	})
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
