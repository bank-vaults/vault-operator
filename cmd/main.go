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

package main

import (
	"flag"
	"net"
	"os"
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	"github.com/bank-vaults/vault-operator/pkg/apis"
	vaultv1alpha1 "github.com/bank-vaults/vault-operator/pkg/apis/vault/v1alpha1"
	"github.com/bank-vaults/vault-operator/pkg/controller"
)

const (
	envOperatorNamespace   = "OPERATOR_NAMESPACE"
	envWatchNamespace      = "WATCH_NAMESPACE"
	envKubeServiceHost     = "KUBERNETES_SERVICE_HOST"
	envKubeServicePort     = "KUBERNETES_SERVICE_PORT"
	envBankVaultsImage     = "BANK_VAULTS_IMAGE"
	healthProbeBindAddress = ":8080"
	metricsBindAddress     = ":8383"
	defaultSyncPeriod      = 30 * time.Second
)

var log = ctrl.Log.WithName("cmd")

func main() {
	// Register CLI flags
	syncPeriod := flag.Duration("sync_period", defaultSyncPeriod,
		"Determines the minimum frequency at which watched resources are reconciled")
	verbose := flag.Bool("verbose", false, "Enables verbose logging")
	flag.Parse()

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	ctrl.SetLogger(zap.New(zap.UseDevMode(*verbose)))

	// Update default bank vaults image if needed
	defaultImage := os.Getenv(envBankVaultsImage)
	if defaultImage != "" {
		vaultv1alpha1.DefaultBankVaultsImage = defaultImage
	}

	// Get namespace config
	namespace := os.Getenv(envOperatorNamespace)
	if namespace == "" {
		namespace = os.Getenv(envWatchNamespace)
	}

	namespaces := make(map[string]cache.Config)
	if namespace == "" {
		log.Info("no watched namespace found, watching the entire cluster")
	} else {
		namespaces[namespace] = cache.Config{}
		log.Info("watched namespace: " + namespace)
	}

	// Load kube client config
	k8sConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "unable to get k8s config")
		os.Exit(1)
	}

	// Configure leader election
	host := os.Getenv(envKubeServiceHost)
	port := os.Getenv(envKubeServicePort)
	leaderElectionNamespace := ""
	if k8sConfig.Host != "https://"+net.JoinHostPort(host, port) {
		leaderElectionNamespace = "default"
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(k8sConfig, manager.Options{
		Cache: cache.Options{
			DefaultNamespaces: namespaces,
			SyncPeriod:        syncPeriod,
		},
		LeaderElection:          true,
		LeaderElectionNamespace: leaderElectionNamespace,
		LeaderElectionID:        "vault-operator-lock",
		HealthProbeBindAddress:  healthProbeBindAddress,
		Metrics:                 metricsserver.Options{BindAddress: metricsBindAddress},
		LivenessEndpointName:    "/",      // For Chart backwards compatibility
		ReadinessEndpointName:   "/ready", // For Chart backwards compatibility
	})
	if err != nil {
		log.Error(err, "unable to create manager as defined")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	// Register checks
	log.Info("registering manager checks")

	err = mgr.AddReadyzCheck("ping", healthz.Ping)
	if err != nil {
		log.Error(err, "unable to add readyz check")
		os.Exit(1)
	}

	err = mgr.AddHealthzCheck("ping", healthz.Ping)
	if err != nil {
		log.Error(err, "unable to add heatlhz check")
		os.Exit(1)
	}

	// Setup scheme and controller
	log.Info("bootstrapping manager")

	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "unable to add scheme to manager")
		os.Exit(1)
	}

	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "unable to add manager to controller")
		os.Exit(1)
	}

	// Start manager
	log.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "problem running manager")
		os.Exit(1)
	}
}
