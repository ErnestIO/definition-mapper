/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/r3labs/binary-prefix"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

// MapInstances : Maps the instances for the input payload on a ernest internal format
func MapInstances(payload input.Payload, m output.FSMMessage) (instances []output.Instance, err error) {
	instances = appendSaltInstance(instances, payload, m.ServiceName)

	for _, instance := range payload.Service.Instances {
		if valid, err := instance.IsValid(); valid == false {
			return instances, err
		}

		imageParts := strings.Split(instance.Image, "/")
		imageCatalog := imageParts[0]
		imageImage := imageParts[1]

		instanceNetwork := payload.Datacenter.Name + "-" + m.ServiceName + "-" + instance.Networks.Name

		// Convert net.IP to net.IPv4 so that we can increment it
		startIP := net.ParseIP(instance.Networks.StartIP.String()).To4()
		ip := instance.Networks.StartIP.To4()

		for index := 0; index < instance.Count; index++ {
			if m.Networks.Items != nil {
				network, err := getInstanceNetwork(m, instanceNetwork, instance.Networks.StartIP)
				if err != nil {
					return instances, err
				}

				err = isValidInstanceIP(network, ip, startIP, m.Instances.Items)
				if err != nil {
					return instances, err
				}
			} else {
				// Leave the network name as the user specified if there are no networks to build
				instanceNetwork = instance.Networks.Name
			}

			name := payload.Datacenter.Name + "-" + m.ServiceName + "-" + instance.Name + "-" + strconv.Itoa(index+1)
			i, err := buildInstance(instance, name, instanceNetwork, ip, imageCatalog, imageImage)
			if err != nil {
				return instances, err
			}

			instances = append(instances, i)

			// Increment IP address
			ip[3]++
		}
	}
	return instances, err
}

// MapInstancesToUpdate : Returns a list of instances being updated on its
// CPU, Memory and Disks resources
func MapInstancesToUpdate(old *output.FSMMessage, current output.FSMMessage) (instances []output.Instance, err error) {
	updatedInstances := []output.Instance{}

	for _, o := range old.Instances.Items {
		for _, n := range current.Instances.Items {
			if n.Name == o.Name {
				if hasChangedInstanceResource(o, n) {
					in := output.Instance{
						Name:        o.Name,
						Catalog:     o.Catalog,
						Image:       o.Image,
						Cpus:        n.Cpus,
						Memory:      n.Memory,
						NetworkName: o.NetworkName,
						IP:          o.IP,
					}
					in.Disks = n.Disks
					updatedInstances = append(updatedInstances, in)
				}
			}
		}
	}

	return updatedInstances, nil
}

// MapInstancesToDelete : Calculates a diff between existing instances and given ones
func MapInstancesToDelete(old *output.FSMMessage, current output.FSMMessage) (instances []output.Instance, err error) {
	for i := range current.Instances.Items {
		current.Instances.Items[i].Exists = false
	}

	for j, o := range old.Instances.Items {
		old.Instances.Items[j].Exists = false
		for i, n := range current.Instances.Items {
			if n.Name == o.Name && n.Exists == false {
				current.Instances.Items[i].Exists = true
				old.Instances.Items[j].Exists = true
				break
			}
		}
	}

	deletedInstances := []output.Instance{}
	for _, o := range old.Instances.Items {
		if o.Exists == false {
			// TODO remove network_id from everywhere
			network := o.NetworkName
			if network == "" {
				network = o.NetworkName
			}
			in := output.Instance{
				Name:        o.Name,
				Catalog:     o.Catalog,
				Image:       o.Image,
				Cpus:        o.Cpus,
				Memory:      o.Memory,
				NetworkName: network,
				IP:          o.IP,
				Disks:       o.Disks,
			}
			deletedInstances = append(deletedInstances, in)
		}
	}

	return deletedInstances, nil
}

// MapInstancesToCreate : Calculates a diff between given instances and existing ones
func MapInstancesToCreate(old *output.FSMMessage, current output.FSMMessage) (instances []output.Instance, err error) {
	for i := range old.Instances.Items {
		old.Instances.Items[i].Exists = false
	}

	for j := range current.Instances.Items {
		current.Instances.Items[j].Exists = false
		for i, o := range old.Instances.Items {
			if current.Instances.Items[j].Name == o.Name && o.Exists == false {
				old.Instances.Items[i].Exists = true
				current.Instances.Items[j].Exists = true
				break
			}
		}
	}

	newInstances := []output.Instance{}
	for _, n := range current.Instances.Items {
		if n.Exists == false {
			newInstances = append(newInstances, n)
		}
	}

	return newInstances, nil
}

// Calculates if there is any change on the resources of given instances
func hasChangedInstanceResource(o output.Instance, n output.Instance) bool {
	if o.Cpus != n.Cpus {
		return true
	}
	if o.Memory != n.Memory {
		return true
	}
	if len(o.Disks) != len(n.Disks) {
		return true
	}
	for _, do := range o.Disks {
		sw := false
		for _, dn := range n.Disks {
			if do == dn {
				sw = true
			}
		}
		if sw == false {
			return true
		}
	}

	return false
}

// Appends a salt instance if bootstrapping mode is salt
func appendSaltInstance(instances []output.Instance, payload input.Payload, serviceName string) []output.Instance {
	if payload.Service.IsSaltBootstrapped() {
		instances = append(instances, output.Instance{
			Name:        payload.Datacenter.Name + "-" + serviceName + "-salt-master",
			Catalog:     "r3",
			Image:       "r3-salt-master",
			Cpus:        1,
			Memory:      2048,
			Disks:       []output.InstanceDisk{},
			NetworkName: payload.Datacenter.Name + "-" + serviceName + "-salt",
			IP:          net.ParseIP("10.254.254.100"),
		})
	}

	return instances
}

// Get the instance related network
func getInstanceNetwork(m output.FSMMessage, instanceNetwork string, startIP net.IP) (*net.IPNet, error) {
	var subnet string
	isNetwork := false

	for _, network := range m.Networks.Items {
		if network.Name == instanceNetwork {
			subnet = network.Subnet
			isNetwork = true
		}
	}

	if !isNetwork {
		err := errors.New("Instance must specify a valid network defined on the spec!")
		return nil, err
	}

	_, network, err := net.ParseCIDR(subnet)
	if err != nil {
		err := errors.New("Could not parse Network Subnet")
		return nil, err
	}

	inRange := network.Contains(startIP)
	if !inRange {
		err := errors.New("Start IP invalid. IP must be a valid IP in the same range as it's network")
		return nil, err
	}

	return network, nil
}

// Check if an instance has a valid ip for its related network
func isValidInstanceIP(network *net.IPNet, ip net.IP, startIP net.IP, instances []output.Instance) error {
	// Check IP is in Range
	inRange := network.Contains(ip)
	if !inRange {
		err := errors.New("Instance IP invalid. IP must be a valid IP in the same range as it's network")
		return err
	}

	// Check IP is greater than Start IP (Bounds checking)
	if ip[3] < startIP[3] {
		err := errors.New("Instance IP invalid. Allocated IP is lower than Start IP")
		return err
	}
	// Check IP Is available
	for _, otherInstance := range instances {
		if otherInstance.IP.String() == ip.String() {
			err := errors.New("Instance IP is already assigned on this network.")
			return err
		}
	}

	return nil
}

// Builds a single output instance from a validated input instance
func buildInstance(instance input.Instance, name string, instanceNetwork string, ip net.IP, imageCatalog string, imageImage string) (output.Instance, error) {
	var err error

	i := output.Instance{}
	i.Name = name
	i.Catalog = imageCatalog
	i.Image = imageImage
	i.Cpus = instance.Cpus

	i.Memory, err = binaryprefix.GetMB(instance.Memory)
	if err != nil {
		err := errors.New("Invalid memory format")
		return i, err
	}

	disks := make([]output.InstanceDisk, len(instance.Disks))
	for j, d := range instance.Disks {
		size, err := binaryprefix.GetMB(d)
		if err != nil {
			err := errors.New("Invalid disk format")
			return i, err
		}

		disks[j] = output.InstanceDisk{
			ID:   (j + 1),
			Size: size,
		}
	}
	i.Disks = disks
	i.NetworkName = instanceNetwork
	i.IP = net.ParseIP(ip.String())

	return i, nil
}

// MapResultingInstances maps all resulting instances to output instance struct
func MapResultingInstances(result []output.Instance) []output.Instance {
	out := make([]output.Instance, len(result))
	for i, instance := range result {
		out[i].Name = instance.Name
		out[i].Cpus = instance.Cpus
		out[i].NetworkName = instance.NetworkName
		out[i].Memory = instance.Memory
	}
	return out
}
