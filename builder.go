/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"

	"github.com/r3labs/definition-mapper/builder"
	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
	"github.com/r3labs/definition-mapper/workflow"
)

// BuildFSMMessage : maps input message on a valid internal ernest message
func BuildFSMMessage(payload input.Payload, prev *output.FSMMessage) (m output.FSMMessage, err error) {
	m, err = mapGenericMessage(payload, prev)
	if err != nil {
		return m, err
	}

	if m.Bootstrapping == "none" && (len(m.Executions.Items) > 0 || len(m.Bootstraps.Items) > 0) {
		return m, errors.New("No bootstrapping mode specified but executions found, \ntry adding a 'bootstrapping: salt' to your definition")
	}

	if prev == nil {
		return m, nil
	}

	m.Endpoint = prev.Endpoint

	if m.Bootstrapping != prev.Bootstrapping {
		return m, errors.New("Can't modify the bootstrapping field on an already existing environment")
	}

	w := workflow.New("workflow/workflows/service_create.json")

	if err = mapUpdatedNetworks(payload, &m, prev, w); err != nil {
		return m, err
	}

	if err = mapUpdatedInstances(payload, &m, prev, w); err != nil {
		return m, err
	}

	if err = mapUpdatedFirewalls(payload, &m, prev, w); err != nil {
		return m, err
	}

	if err = mapUpdatedNats(payload, &m, prev, w); err != nil {
		return m, err
	}

	if err = mapUpdatedExecutions(payload, &m, prev, w); err != nil {
		return m, err
	}

	err = w.Finish("workflow/workflows/service_create_done.json")
	if err != nil {
		return m, err
	}

	m.Workflow = *w
	m.Parent = string(prev.ToJSON())

	return m, nil
}

// BuildDeleteMessage : maps input message to a valid delete message
func BuildDeleteMessage(payload output.FSMMessage) (output.FSMMessage, error) {
	payload.RoutersToDelete = payload.Routers
	payload.NetworksToDelete = payload.Networks
	payload.InstancesToDelete = payload.Instances
	payload.FirewallsToDelete = payload.Firewalls
	payload.NatsToDelete = payload.Nats

	w, err := buildDeleteWorkflow(payload)
	if err != nil {
		return payload, err
	}

	payload.Workflow = *w

	return payload, nil
}

// Maps the generic parts of an input message (datacenter, router, networks, instances,
// nats, firewalls, bootstraps, executions)
func mapGenericMessage(payload input.Payload, prev *output.FSMMessage) (m output.FSMMessage, err error) {
	if valid, err := payload.Service.IsNameValid(); valid == false {
		return m, err
	}

	m = output.FSMMessage{
		ID:            payload.ServiceID,
		Service:       payload.ServiceID,
		ServiceName:   payload.Service.Name,
		Client:        payload.Client.ID,
		ClientName:    payload.Client.Name,
		Type:          payload.Datacenter.Type,
		Bootstrapping: builder.Bootstrapping(payload),
	}

	m.ServiceIP, err = builder.MapServiceIP(payload)
	if err != nil {
		return m, err
	}

	m.ErnestIP, err = builder.MapErnestIP(payload)
	if err != nil {
		return m, err
	}

	m.Datacenters.Items, err = builder.MapDatacenters(payload)
	if err != nil {
		return m, err
	}

	m.Routers.Items, err = builder.MapRouters(payload, prev, m.ServiceIP)
	if err != nil {
		return m, err
	}
	if m.ServiceIP != "" {
		m.Endpoint = m.ServiceIP
		m.Routers.Status = "completed"
	}

	m.Networks.Items, err = builder.MapNetworks(payload, m)
	if err != nil {
		return m, err
	}

	m.Instances.Items, err = builder.MapInstances(payload, m)
	if err != nil {
		return m, err
	}
	m.InstancesToUpdate.Items = m.Instances.Items

	m.Nats.Items, err = builder.MapNATS(payload)
	if err != nil {
		return m, err
	}

	m.Firewalls.Items, err = builder.MapFirewalls(payload)
	if err != nil {
		return m, err
	}

	// If bootstrapping is disabled, dont build executions & bootstraps
	if m.Bootstrapping != "none" {
		m.Bootstraps.Items = builder.GenerateBootstraps(payload, prev, m)
		if err != nil {
			return m, err
		}

		m.Executions.Items, err = builder.MapExecutions(payload, prev, m)
		if err != nil {
			return m, err
		}

		m.ExecutionsToCreate.Items, err = builder.GenerateExecutions(payload, prev, m)
		if err != nil {
			return m, err
		}
	}

	if prev == nil {
		w, err := buildGenericWorkflow(m)
		if err != nil {
			return m, err
		}

		m.Workflow = *w
	}

	return m, nil
}

func buildGenericWorkflow(m output.FSMMessage) (*workflow.Workflow, error) {
	w := workflow.New("workflow/workflows/service_create.json")

	if len(m.Routers.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/routers_create.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.Networks.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/networks_create.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.Instances.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/instance_create.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.Firewalls.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/firewalls_create.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.Nats.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/nats_create.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.Bootstraps.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/bootstraps.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.ExecutionsToCreate.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/executions.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	err := w.Finish("workflow/workflows/service_create_done.json")
	if err != nil {
		return w, err
	}

	return w, nil
}

func buildDeleteWorkflow(m output.FSMMessage) (*workflow.Workflow, error) {
	w := workflow.New("workflow/workflows/service_delete.json")

	if len(m.InstancesToDelete.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/instance_delete.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.NetworksToDelete.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/networks_delete.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	if len(m.RoutersToDelete.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/routers_delete.json")
		if err != nil {
			return w, err
		}
		w.Add(arc)
	}

	err := w.Finish("workflow/workflows/service_delete_done.json")
	if err != nil {
		return w, err
	}

	return w, nil
}

// In case bootstraps or executions has changed this will map it on a valid input message
func mapUpdatedExecutions(payload input.Payload, m *output.FSMMessage, prev *output.FSMMessage, w *workflow.Workflow) error {
	if m.Bootstrapping == "none" {
		return nil
	}

	m.Bootstraps.Items = builder.GenerateBootstraps(payload, prev, *m)
	if len(m.Bootstraps.Items) > 0 {
		if !w.ContainsSubject("bootstrapping") {
			arc, err := workflow.LoadArcs("workflow/workflows/bootstraps.json")
			if err != nil {
				return err
			}
			w.Add(arc)
		}
	}

	exec, _ := builder.GenerateExecutions(payload, prev, *m)
	m.ExecutionsToCreate.Items = exec
	if len(exec) > 0 {
		if !w.ContainsSubject("running_executions") {
			arc, err := workflow.LoadArcs("workflow/workflows/executions.json")
			if err != nil {
				return err
			}
			w.Add(arc)
		}
	}

	cleanupExecs := builder.BootstrapCleanup(*m)
	if len(cleanupExecs) > 0 {
		if !w.ContainsSubject("running_executions") {
			arc, err := workflow.LoadArcs("workflow/workflows/executions.json")
			if err != nil {
				return err
			}
			w.Add(arc)
		}
		for _, exec := range cleanupExecs {
			m.ExecutionsToCreate.Items = append(m.ExecutionsToCreate.Items, exec)
		}
	}

	return nil
}

// Maps updated rules for firewalls in case they are changed
func mapUpdatedFirewalls(payload input.Payload, m *output.FSMMessage, prev *output.FSMMessage, w *workflow.Workflow) error {
	oldFirewalls := prev.Firewalls.Items
	newFirewalls := m.Firewalls
	if builder.HasChangedFirewalls(oldFirewalls, newFirewalls.Items) {
		arc, err := workflow.LoadArcs("workflow/workflows/firewalls_update.json")
		if err != nil {
			return err
		}
		w.Add(arc)
	}
	return nil
}

// Maps updated nats in case they are changed
func mapUpdatedNats(payload input.Payload, m *output.FSMMessage, prev *output.FSMMessage, w *workflow.Workflow) error {
	oldNats := prev.Nats.Items
	newNats := m.Nats.Items
	if builder.HasChangedNats(oldNats, newNats) {
		arc, err := workflow.LoadArcs("workflow/workflows/nats_update.json")
		if err != nil {
			return err
		}
		w.Add(arc)
	}
	return nil
}

// Maps updated instances in case they are changed, new ones to be created or deleted
func mapUpdatedInstances(payload input.Payload, m *output.FSMMessage, prev *output.FSMMessage, w *workflow.Workflow) (err error) {
	m.InstancesToUpdate.Items, err = builder.MapInstancesToUpdate(prev, *m)
	if err != nil {
		return err
	}

	m.InstancesToDelete.Items, err = builder.MapInstancesToDelete(prev, *m)
	if err != nil {
		return err
	}

	m.InstancesToCreate.Items, err = builder.MapInstancesToCreate(prev, *m)
	if err != nil {
		return err
	}

	if len(m.InstancesToCreate.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/instance_create.json")
		if err != nil {
			return err
		}
		w.Add(arc)

		for _, i := range m.InstancesToCreate.Items {
			sw := 0
			for _, j := range m.InstancesToUpdate.Items {
				if i.Name == j.Name {
					sw = 1
					break
				}
			}
			if sw == 0 {
				m.InstancesToUpdate.Items = append(m.InstancesToUpdate.Items, i)
			}
		}
	}

	if len(m.InstancesToUpdate.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/instances_update.json")
		if err != nil {
			return err
		}
		w.Add(arc)
	}

	if len(m.InstancesToDelete.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/instance_delete.json")
		if err != nil {
			return err
		}
		w.Add(arc)
	}

	return nil
}

// Maps updated networks in case they are new ones
func mapUpdatedNetworks(payload input.Payload, m *output.FSMMessage, prev *output.FSMMessage, w *workflow.Workflow) (err error) {
	m.NetworksToCreate.Items, err = builder.MapNetworksToCreate(prev, *m)
	if err != nil {
		return err
	}
	if len(m.NetworksToCreate.Items) > 0 {
		arc, err := workflow.LoadArcs("workflow/workflows/networks_create.json")
		if err != nil {
			return err
		}
		w.Add(arc)
	}

	return nil
}
