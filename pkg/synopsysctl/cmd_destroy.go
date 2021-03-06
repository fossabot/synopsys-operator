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
	"strings"

	"github.com/blackducksoftware/synopsys-operator/pkg/soperator"
	util "github.com/blackducksoftware/synopsys-operator/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Destroy Command Defaults
var destroyNamespace = "synopsys-operator"

// destroyCmd removes the Synopsys-Operator from the cluster
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Removes the Synopsys-Operator and CRDs from Cluster",
	Args: func(cmd *cobra.Command, args []string) error {
		// Check number of arguments
		if len(args) != 0 {
			return fmt.Errorf("this command accepts 0 arguments")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the namespace of the Synopsys-Operator
		destroyNamespace, err := soperator.GetOperatorNamespace(restconfig)
		if err != nil {
			log.Warnf("Error finding Synopsys-Operator: %s", err)
		}
		log.Debugf("Destroying the Synopsys-Operator: %s\n", destroyNamespace)
		// Delete the namespace
		out, err := util.RunKubeCmd(restconfig, kube, openshift, "delete", "ns", destroyNamespace)
		if err != nil {
			log.Warnf("Could not delete %s - %s", destroyNamespace, out)
		}
		cleanCommands := [...]string{
			"delete crd alerts.synopsys.com",
			"delete crd blackducks.synopsys.com",
			"delete crd hubs.synopsys.com",
			"delete crd opssights.synopsys.com",
			"delete clusterrolebinding synopsys-operator-admin",
			"delete clusterrole synopsys-operator-admin",
		}

		for cmd := range cleanCommands {
			out, _ = util.RunKubeCmd(restconfig, kube, openshift, strings.Split(cleanCommands[cmd], " ")...)
			log.Debugf(" > %s", out)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
}
