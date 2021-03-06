/*
Copyright (C) 2019 Synopsys, Inc.

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements. See the NOTICE file
distributed with this work for additional information
regarding copyright ownership. The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License. You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied. See the License for the
specific language governing permissions and limitations
under the License.
*/

package synopsysctl

import (
	"fmt"

	horizonapi "github.com/blackducksoftware/horizon/pkg/api"
	horizoncomponents "github.com/blackducksoftware/horizon/pkg/components"
	"github.com/blackducksoftware/horizon/pkg/deployer"
	soperator "github.com/blackducksoftware/synopsys-operator/pkg/soperator"
	"github.com/blackducksoftware/synopsys-operator/pkg/util"
	operatorutil "github.com/blackducksoftware/synopsys-operator/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

//  Deploy Command Defaults
var exposeUI = false
var deployNamespace = "synopsys-operator"
var deploySynopsysOperatorImage = "docker.io/blackducksoftware/synopsys-operator:2019.2.0-RC"
var deployPrometheusImage = "docker.io/prom/prometheus:v2.1.0"
var deployTerminationGracePeriodSeconds int64 = 180
var deployDockerConfigPath = ""
var deploySecretType = "Opaque"
var deploySecretAdminPassword = "blackduck"
var deploySecretPostgresPassword = "blackduck"
var deploySecretUserPassword = "blackduck"
var deploySecretBlackduckPassword = "blackduck"

// Deploy Global Variables
var secretType horizonapi.SecretType

// deployCmd represents the deploy command
var deployCmd = &cobra.Command{
	Use:   "deploy [NAMESPACE]",
	Short: "Deploys the synopsys operator onto your cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		// Check number of arguments
		if len(args) > 1 {
			return fmt.Errorf("this command only accepts up to 1 argument")
		}
		// Check the Secret Type
		var err error
		secretType, err = operatorutil.SecretTypeNameToHorizon(deploySecretType)
		if err != nil {
			log.Errorf("Failed to Convert Secret Type: %s", err)
			return nil
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Debugf("Deploying the Synopsys-Operator: %s\n", deployNamespace)
		// Read Commandline Parameters
		if len(args) == 1 {
			deployNamespace = args[0]
		}
		// check if operator is already installed
		out, err := util.RunKubeCmd(restconfig, kube, openshift, "get", "clusterrolebindings", "synopsys-operator-admin", "-o", "go-template='{{range .subjects}}{{.namespace}}{{end}}'")
		if err == nil {
			log.Errorf("Synopsys-Operator is already installed in namespace %s.", out)
			return nil
		}

		// Create a Horizon Deployer to set up the environment for the Synopsys Operator
		environmentDeployer, err := deployer.NewDeployer(restconfig)

		// create a new namespace
		ns := horizoncomponents.NewNamespace(horizonapi.NamespaceConfig{
			// APIVersion:  "string",
			// ClusterName: "string",
			Name:      deployNamespace,
			Namespace: deployNamespace,
		})
		environmentDeployer.AddNamespace(ns)

		// Deploy Resources for the Synopsys Operator
		err = environmentDeployer.Run()
		if err != nil {
			log.Errorf("Error deploying Environment with Horizon : %s", err)
			return nil
		}

		sealKey, err := operatorutil.GetRandomString(32)
		if err != nil {
			log.Panicf("unable to generate the random string for SEAL_KEY due to %+v", err)
		}

		// Deploy synopsys-operator
		soperatorSpec := soperator.SpecConfig{
			Namespace:                     deployNamespace,
			SynopsysOperatorImage:         deploySynopsysOperatorImage,
			SecretType:                    secretType,
			SecretAdminPassword:           deploySecretAdminPassword,
			SecretPostgresPassword:        deploySecretPostgresPassword,
			SecretUserPassword:            deploySecretUserPassword,
			SecretBlackduckPassword:       deploySecretBlackduckPassword,
			TerminationGracePeriodSeconds: deployTerminationGracePeriodSeconds,
			SealKey:                       sealKey,
		}
		err = soperator.UpdateSOperatorComponents(restconfig, kubeClient, deployNamespace, &soperatorSpec)
		if err != nil {
			log.Errorf("Error deploying Synopsys Operator: %s", err)
			return nil
		}

		// Deploy prometheus
		promtheusSpec := soperator.PrometheusSpecConfig{
			Namespace:       deployNamespace,
			PrometheusImage: deployPrometheusImage,
		}
		err = soperator.UpdatePrometheus(restconfig, kubeClient, deployNamespace, &promtheusSpec)
		if err != nil {
			log.Errorf("Error deploying Prometheus: %s", err)
			return nil
		}

		// create secrets (TDDO I think this only works on OpenShift)
		util.RunKubeCmd(restconfig, kube, openshift, "create", "secret", "generic", "custom-registry-pull-secret", fmt.Sprintf("--from-file=.dockerconfigjson=%s", deployDockerConfigPath), "--type=kubernetes.io/dockerconfigjson")
		util.RunKubeCmd(restconfig, kube, openshift, "secrets", "link", "default", "custom-registry-pull-secret", "--for=pull")
		util.RunKubeCmd(restconfig, kube, openshift, "secrets", "link", "synopsys-operator", "custom-registry-pull-secret", "--for=pull")
		util.RunKubeCmd(restconfig, kube, openshift, "scale", "replicationcontroller", "synopsys-operator", "--replicas=0")
		util.RunKubeCmd(restconfig, kube, openshift, "scale", "replicationcontroller", "synopsys-operator", "--replicas=1")

		// expose the routes
		if exposeUI {
			out, err = util.RunKubeCmd(restconfig, kube, openshift, "expose", "replicationcontroller", "synopsys-operator", "--port=80", "--target-port=3000", "--name=synopsys-operator-tcp", "--type=LoadBalancer", fmt.Sprintf("--namespace=%s", deployNamespace))
			if err != nil {
				log.Warnf("Error exposing the Synopsys-Operator's Replication Controller: %s", out)
			}
			out, err = util.RunKubeCmd(restconfig, kube, openshift, "create", "route", "edge", "--service=synopsys-operator-tcp", "-n", deployNamespace)
			if err != nil {
				log.Warnf("Could not create route (Possible Reason: Kubernetes doesn't support Routes): %s", out)
			}
		} else {
			log.Warnf("Synopsys-Operator UI is not exposed ( --expose-ui=true to expose )")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().BoolVar(&exposeUI, "expose-ui", exposeUI, "Expose the Synopsys-Operator's User Interface")
	deployCmd.Flags().StringVarP(&deploySynopsysOperatorImage, "synopsys-operator-image", "i", deploySynopsysOperatorImage, "synopsys operator image URL")
	deployCmd.Flags().StringVarP(&deployPrometheusImage, "prometheus-image", "p", deployPrometheusImage, "prometheus image URL")
	deployCmd.Flags().StringVarP(&deployDockerConfigPath, "docker-config", "d", deployDockerConfigPath, "path to docker config (image pull secrets etc)")
	deployCmd.Flags().StringVar(&deploySecretType, "secret-type", deploySecretType, "type of kubernetes secret for postgres and blackduck")
	deployCmd.Flags().StringVar(&deploySecretAdminPassword, "admin-password", deploySecretAdminPassword, "postgres admin password")
	deployCmd.Flags().StringVar(&deploySecretPostgresPassword, "postgres-password", deploySecretPostgresPassword, "postgres password")
	deployCmd.Flags().StringVar(&deploySecretUserPassword, "user-password", deploySecretUserPassword, "postgres user password")
	deployCmd.Flags().StringVar(&deploySecretBlackduckPassword, "blackduck-password", deploySecretBlackduckPassword, "blackduck password for 'sysadmin' account")
	deployCmd.Flags().Int64VarP(&deployTerminationGracePeriodSeconds, "postgres-termination-grace-period", "t", deployTerminationGracePeriodSeconds, "termination grace period in seconds for shutting down postgres")

	// Set Log Level
	log.SetLevel(log.DebugLevel)
}
