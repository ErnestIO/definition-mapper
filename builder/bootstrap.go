/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"fmt"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

// GenerateBootstraps : generates necessary bootstraps for instances
func GenerateBootstraps(payload input.Payload, prev *output.FSMMessage, m output.FSMMessage) (bootstraps []output.Execution) {
	// Check there are any executions
	// Validate and generate Executions and Bootstraps
	// nodes with executions
	var instances []output.Instance
	if prev == nil {
		instances = m.Instances.Items
	} else {
		instances = m.InstancesToCreate.Items
	}

	// Add instance to bootstrap if not salt-master
	for _, instance := range instances {
		if instance.IP.String() != "10.254.254.100" {
			command := fmt.Sprintf("/usr/bin/bootstrap -master 10.254.254.100 -host %s -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name %s", instance.IP, instance.Name)
			name := fmt.Sprintf("Bootstrap %s", instance.Name)
			bootstrap := output.Execution{
				Name:    name,
				Type:    "salt",
				Target:  "list:salt-master.localdomain",
				Payload: command,
			}
			bootstraps = append(bootstraps, bootstrap)
		}
	}

	return bootstraps
}

// BootstrapCleanup : When a service updates its bootstraps may need a salt cleanup
func BootstrapCleanup(m output.FSMMessage) []output.Execution {
	cleanupExecs := []output.Execution{}

	// Add instance to bootstrap if not salt-master
	for _, instance := range m.InstancesToDelete.Items {
		if instance.IP.String() != "10.254.254.100" && !instance.Exists {
			command := fmt.Sprintf("salt-key -y -d %s", instance.Name)
			name := fmt.Sprintf("Bootstrap %s", instance.Name)
			cmd := output.Execution{
				Name:    name,
				Type:    "salt",
				Target:  "list:salt-master.localdomain",
				Payload: command,
			}
			cleanupExecs = append(cleanupExecs, cmd)
		}
	}

	return cleanupExecs
}

// Bootstrapping : gets the specific defined bootstraping method
func Bootstrapping(payload input.Payload) string {
	if payload.Service.Bootstrapping == "salt" {
		return "salt"
	}

	return "none"
}
