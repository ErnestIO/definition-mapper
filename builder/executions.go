/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

// MapExecutions reassigns any previous executions to the current state
func MapExecutions(payload input.Payload, prev *output.FSMMessage, m output.FSMMessage) (execs []output.Execution, err error) {
	if prev == nil {
		return execs, nil
	}

	for _, exec := range prev.Executions.Items {
		newExec := output.Execution{
			Name:    exec.Name,
			Type:    exec.Type,
			Target:  exec.Target,
			Payload: exec.Payload,
		}
		execs = append(execs, newExec)
	}

	return execs, nil
}

// GenerateExecutions : Generates the given executions to a valid internal executions
func GenerateExecutions(payload input.Payload, prev *output.FSMMessage, m output.FSMMessage) (execs []output.Execution, err error) {
	for _, instance := range payload.Service.Instances {
		// Instance Name prefix
		prefix := fmt.Sprintf("%s-%s-%s", payload.Datacenter.Name, payload.Service.Name, instance.Name)

		// Generate all instance names
		instanceNames := make([]string, instance.Count)
		for q := 0; q < instance.Count; q++ {
			instanceNames[q] = fmt.Sprintf("%s-%s", prefix, strconv.Itoa(q+1))
		}

		/// Generate Execution target
		nodes := strings.Join(instanceNames, ",")
		target := fmt.Sprintf("list:%s", nodes)

		// Itterate over each execution
		for q, e := range instance.Provisioner {
			name := fmt.Sprintf("Execution %s %s", instance.Name, strconv.Itoa(q+1))
			execPayload := strings.Join(e.Commands, "; ")

			// Construct the execution and its payload
			execution := output.Execution{
				Name:    name,
				Type:    "salt",
				Target:  target,
				Payload: execPayload,
			}

			// Check the old Execution payload for changes
			// Otherwise, append the execution unchanged
			if prev != nil {
				// Determine if there is an old execution
				if oldExecution := prev.GetExecution(execution.Name); oldExecution != nil {
					// Get any new instances that match current instance group's name
					newInstances := m.FilterNewInstances(prefix)

					// If the execution has changed, run the execution against all instances in the Type
					// If there are no changes, but new instances, modify the target and run only against them
					if execution.Payload != oldExecution.Payload {
						execs = append(execs, execution)
					} else if len(newInstances) > 0 {
						execution.Target = constructTarget(newInstances)
						execs = append(execs, execution)
					}

				} else {
					// If it's a new execution, append it!
					execs = append(execs, execution)
				}

			} else {
				execs = append(execs, execution)
			}
		}
	}
	return execs, err
}

func constructTarget(instances []output.Instance) string {
	var targets []string
	for _, instance := range instances {
		targets = append(targets, instance.Name)
	}
	nodes := strings.Join(targets, ",")
	return fmt.Sprintf("list:%s", nodes)
}
