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
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"time"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/bank-vaults/vault-operator/pkg/apis"
	"github.com/bank-vaults/vault-operator/pkg/controller"
)

var log = ctrl.Log.WithName("cmd")

const (
	operatorNamespace      = "OPERATOR_NAMESPACE"
	watchNamespaceEnvVar   = "WATCH_NAMESPACE"
	healthProbeBindAddress = ":8080"
	metricsBindAddress     = ":8383"
)

func main() {
	syncPeriod := flag.Duration("sync_period", 30*time.Second, "SyncPeriod determines the minimum frequency at which watched resources are reconciled")
	verbose := flag.Bool("verbose", false, "enable verbose logging")

	flag.Parse()

	// The logger instantiated here can be changed to any logger
	// implementing the logr.Logger interface. This logger will
	// be propagated through the whole operator, generating
	// uniform and structured logs.
	ctrl.SetLogger(zap.New(zap.UseDevMode(*verbose)))

	// Fetch namespace data
	namespace, isSet := os.LookupEnv(operatorNamespace)
	if !isSet {
		namespace, isSet = os.LookupEnv(watchNamespaceEnvVar)
	}

	var namespaces []string
	if !isSet {
		log.Info("No watched namespace found, watching the entire cluster")
	} else {
		log.Info("Watched namespace: " + namespace)
		namespaces = []string{namespace}
	}

	// Get a config to talk to the apiserver
	k8sConfig, err := config.GetConfig()
	if err != nil {
		log.Error(err, "Unable to get k8s config")
		os.Exit(1)
	}

	leaderElectionNamespace := ""
	if !isInClusterConfig(k8sConfig) {
		leaderElectionNamespace = "default"
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := manager.New(k8sConfig, manager.Options{
		Cache: cache.Options{
			Namespaces: namespaces,
			SyncPeriod: syncPeriod,
		},
		LeaderElection:          true,
		LeaderElectionNamespace: leaderElectionNamespace,
		LeaderElectionID:        "vault-operator-lock",
		HealthProbeBindAddress:  healthProbeBindAddress,
		LivenessEndpointName:    "/",      // For Chart backwards compatibility
		ReadinessEndpointName:   "/ready", // For Chart backwards compatibility
		MetricsBindAddress:      metricsBindAddress,
	})
	if err != nil {
		log.Error(err, "Unable to create manager as defined")
		os.Exit(1)
	}

	err = mgr.AddReadyzCheck("ping", healthz.Ping)
	if err != nil {
		log.Error(err, "Add Readyz Check failed")
		os.Exit(1)
	}
	err = mgr.AddHealthzCheck("ping", healthz.Ping)
	if err != nil {
		log.Error(err, "Unable to add heatlh check")
		os.Exit(1)
	}

	log.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "Failed to use api to add scheme")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		log.Error(err, "Unable to add manager to controller")
		os.Exit(1)
	}

	log.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		log.Error(err, "manager exited non-zero")
		os.Exit(1)
	}
}

func isInClusterConfig(k8sConfig *rest.Config) bool {
	host, port := os.Getenv("KUBERNETES_SERVICE_HOST"), os.Getenv("KUBERNETES_SERVICE_PORT")
	return k8sConfig.Host == "https://"+net.JoinHostPort(host, port)
}
