/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"errors"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

// MapNetworks : Maps the networks from a given input payload. Additionally it will create an
// extra salt-master network if the service bootstrapping mode is defined as salt
func MapNetworks(payload input.Payload, m output.FSMMessage) (networks []output.Network, err error) {
	for _, r := range payload.Service.Routers {
		networks = appendSaltNetwork(networks, payload, m)
		networkNames := []string{}
		for _, network := range r.Networks {
			networkNames = append(networkNames, network.Name)
		}

		for i, network := range r.Networks {
			for j, name := range networkNames {
				if name == network.Name && i != j {
					return networks, errors.New("Input includes duplicate network names")
				}
			}
			if valid, err := network.IsValid(); valid == false {
				return networks, err
			}

			n := output.Network{}
			n.Name = payload.Datacenter.Name + "-" + m.ServiceName + "-" + network.Name
			n.RouterName = r.Name
			n.Subnet = network.Subnet
			n.DNS = network.DNS

			networks = append(networks, n)
		}
	}

	return networks, err
}

// MapNetworksToCreate : calculates the network count difference between the existing
// service and the new one, and returns the network diff to be created.
//
// It bases the difference on the network name, so updating network name can cause problems
// on this calculation
func MapNetworksToCreate(old *output.FSMMessage, current output.FSMMessage) (instances []output.Network, err error) {
	networkNames := []string{}
	newNetworks := []output.Network{}
	for _, network := range current.Networks.Items {
		networkNames = append(networkNames, network.Name)
	}

	for i, network := range current.Networks.Items {
		for j, name := range networkNames {
			if name == network.Name && i != j {
				return newNetworks, errors.New("Input includes duplicate network names")
			}
		}
	}

	for i := range old.Networks.Items {
		old.Networks.Items[i].Exists = false
	}

	for j := range current.Networks.Items {
		current.Networks.Items[j].Exists = false
		for i, o := range old.Networks.Items {
			if current.Networks.Items[j].Name == o.Name && o.Exists == false {
				old.Networks.Items[i].Exists = true
				current.Networks.Items[j].Exists = true
				break
			}
		}
	}

	for _, n := range current.Networks.Items {
		if n.Exists == false {
			newNetworks = append(newNetworks, n)
		}
	}

	return newNetworks, err
}

// Appends salt networks if service bootstrapping mode is salt
func appendSaltNetwork(networks []output.Network, p input.Payload, m output.FSMMessage) []output.Network {
	if p.Service.IsSaltBootstrapped() {
		networks = append(networks, output.Network{
			Name:       p.Datacenter.Name + "-" + m.ServiceName + "-salt",
			RouterName: m.Routers.Items[0].Name,
			Subnet:     "10.254.254.0/24",
		})
	}

	return networks
}

// MapResultingNetworks maps all input networks to an network output struct
func MapResultingNetworks(result []output.Network) []output.Network {
	out := make([]output.Network, len(result))
	for i, network := range result {
		out[i].Name = network.Name
		out[i].Subnet = network.Subnet
		out[i].RouterName = network.RouterName
	}
	return out
}
