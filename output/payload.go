/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package output

import (
	"encoding/json"
	"strings"

	"github.com/ernestio/definition-mapper/workflow"
)

// FSMMessage is the JSON payload that will be sent to the FSM to create a
// service.
type FSMMessage struct {
	ID            string            `json:"id"`
	Body          string            `json:"body"`
	Endpoint      string            `json:"endpoint"`
	Service       string            `json:"service"`
	Bootstrapping string            `json:"bootstrapping"`
	ErnestIP      []string          `json:"ernest_ip"`
	ServiceIP     string            `json:"service_ip"`
	Parent        string            `json:"existing_service"`
	Workflow      workflow.Workflow `json:"workflow"`
	ServiceName   string            `json:"name"`
	Client        string            `json:"client"` // TODO: Use client or client_id not both!
	ClientID      string            `json:"client_id"`
	ClientName    string            `json:"client_name"`
	Started       string            `json:"started"`
	Finished      string            `json:"finished"`
	Status        string            `json:"status"`
	Type          string            `json:"type"`
	Datacenters   struct {
		Started  string       `json:"started"`
		Finished string       `json:"finished"`
		Status   string       `json:"status"`
		Items    []Datacenter `json:"items"`
	} `json:"datacenters"`
	Routers struct {
		Started  string   `json:"started"`
		Finished string   `json:"finished"`
		Status   string   `json:"status"`
		Items    []Router `json:"items"`
	} `json:"routers"`
	RoutersToDelete struct {
		Started  string   `json:"started"`
		Finished string   `json:"finished"`
		Status   string   `json:"status"`
		Items    []Router `json:"items"`
	} `json:"routers_to_delete"`
	Networks struct {
		Started  string    `json:"started"`
		Finished string    `json:"finished"`
		Status   string    `json:"status"`
		Items    []Network `json:"items"`
	} `json:"networks"`
	NetworksToCreate struct {
		Started  string    `json:"started"`
		Finished string    `json:"finished"`
		Status   string    `json:"status"`
		Items    []Network `json:"items"`
	} `json:"networks_to_create"`
	NetworksToDelete struct {
		Started  string    `json:"started"`
		Finished string    `json:"finished"`
		Status   string    `json:"status"`
		Items    []Network `json:"items"`
	} `json:"networks_to_delete"`
	Instances struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Instance `json:"items"`
	} `json:"instances"`
	InstancesToCreate struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Instance `json:"items"`
	} `json:"instances_to_create"`
	InstancesToUpdate struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Instance `json:"items"`
	} `json:"instances_to_update"`
	InstancesToDelete struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Instance `json:"items"`
	} `json:"instances_to_delete"`
	Firewalls struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Firewall `json:"items"`
	} `json:"firewalls"`
	FirewallsToDelete struct {
		Started  string     `json:"started"`
		Finished string     `json:"finished"`
		Status   string     `json:"status"`
		Items    []Firewall `json:"items"`
	} `json:"firewalls_to_delete"`
	Loadbalancers struct {
		Started  string         `json:"started"`
		Finished string         `json:"finished"`
		Status   string         `json:"status"`
		Items    []Loadbalancer `json:"items"`
	} `json:"loadbalancers"`
	Nats struct {
		Started  string `json:"started"`
		Finished string `json:"finished"`
		Status   string `json:"status"`
		Items    []Nat  `json:"items"`
	} `json:"nats"`
	NatsToDelete struct {
		Started  string `json:"started"`
		Finished string `json:"finished"`
		Status   string `json:"status"`
		Items    []Nat  `json:"items"`
	} `json:"nats_to_delete"`
	Bootstraps struct {
		Started  string      `json:"started"`
		Finished string      `json:"finished"`
		Status   string      `json:"status"`
		Items    []Execution `json:"items"`
	} `json:"bootstraps"`
	Executions struct {
		Started  string      `json:"started"`
		Finished string      `json:"finished"`
		Status   string      `json:"status"`
		Items    []Execution `json:"items"`
	} `json:"executions"`
	ExecutionsToCreate struct {
		Started  string      `json:"started"`
		Finished string      `json:"finished"`
		Status   string      `json:"status"`
		Items    []Execution `json:"items"`
	} `json:"executions_to_create"`
}

// FilterNewInstances will return any new instances that match a certain pattern
func (m *FSMMessage) FilterNewInstances(name string) []Instance {
	var instances []Instance
	for _, instance := range m.InstancesToCreate.Items {
		if strings.Contains(instance.Name, name) {
			instances = append(instances, instance)
		}
	}
	return instances
}

// GetExecution from service result and return it
func (m *FSMMessage) GetExecution(name string) *Execution {
	for _, exec := range m.Executions.Items {
		if exec.Name == name {
			return &exec
		}
	}
	return nil
}

// ToJSON : Get this service as a json
func (m *FSMMessage) ToJSON() []byte {
	json, _ := json.Marshal(m)

	return json
}
