/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"net"
	"strconv"

	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
	binaryprefix "github.com/r3labs/binary-prefix"
	"github.com/r3labs/graph"
)

// MapInstances : Maps the instances for the input payload on a ernest internal format
func MapInstances(d *definition.Definition) []*components.Instance {
	var instances []*components.Instance

	for _, instance := range d.Instances {
		var commands []string

		ip := net.ParseIP(instance.StartIP).To4()
		memory, _ := binaryprefix.GetMB(instance.Memory)

		for _, prov := range instance.Provisioner {
			if len(prov.Shell) > 0 {
				commands = prov.Shell
			}
		}

		for i := 0; i < instance.Count; i++ {
			var disks []components.Disk

			if instance.RootDisk != "" {
				size, _ := binaryprefix.GetMB(instance.RootDisk)
				disks = append(disks, components.Disk{
					ID:   0,
					Size: size,
				})
			}

			disks = append(disks, MapInstanceDisks(instance.Disks)...)

			newInstance := &components.Instance{
				Name:          instance.Name + "-" + strconv.Itoa(i+1),
				Hostname:      instance.Name + "-" + strconv.Itoa(i+1),
				Catalog:       instance.Catalog(),
				Image:         instance.Template(),
				Cpus:          instance.Cpus,
				Memory:        memory,
				Disks:         disks,
				Network:       instance.Network,
				IP:            ip.String(),
				ShellCommands: commands,
				Tags:          mapInstanceTags(d.Name, instance.Name),
			}

			if len(d.Gateways) < 1 {
				newInstance.InstanceOnly = true
			}

			newInstance.SetDefaultVariables()

			instances = append(instances, newInstance)

			// Increment IP address
			ip[3]++
		}
	}
	return instances
}

// MapInstanceDisks : Maps the instances disks
func MapInstanceDisks(d []string) []components.Disk {
	var disks []components.Disk

	for x, disk := range d {
		size, _ := binaryprefix.GetMB(disk)
		disks = append(disks, components.Disk{
			ID:   (x + 1),
			Size: size,
		})
	}

	return disks
}

func mapInstanceTags(service, instanceGroup string) map[string]string {
	tags := make(map[string]string)

	tags["ernest.service"] = service
	tags["ernest.instance_group"] = instanceGroup

	return tags
}

/*
Count       int      `json:"count"`
Cpus        int      `json:"cpus"`
Image       string   `json:"image"`
Memory      string   `json:"memory"`
RootDisk    string   `json:"root_disk"`
Disks       []string `json:"disks"`
Name        string   `json:"name"`
Network     string   `json:"network"`
StartIP     string   `json:"start_ip"`
Provisioner []Exec   `json:"provisioner"`
*/

// MapDefinitionInstances :
func MapDefinitionInstances(g *graph.Graph) []definition.Instance {
	var instances []definition.Instance

	ci := g.GetComponents().ByType("instance")

	for _, ig := range ci.TagValues("ernest.instance_group") {
		is := ci.ByGroup("ernest.instance_group", ig)

		if len(is) < 1 {
			continue
		}

		firstInstance := is[0].(*components.Instance)

		instance := definition.Instance{
			Name:    ig,
			Cpus:    firstInstance.Cpus,
			Memory:  strconv.Itoa(firstInstance.Memory) + "MB",
			Image:   firstInstance.Image,
			Network: firstInstance.Network,
			StartIP: firstInstance.IP,
			Count:   len(is),
		}

		for _, disk := range firstInstance.Disks {
			size := strconv.Itoa(disk.Size) + "MB"

			if disk.ID == 0 {
				instance.RootDisk = size
				continue
			}

			instance.Disks = append(instance.Disks, size)
		}

		if len(firstInstance.ShellCommands) > 0 {
			instance.Provisioner = append(instance.Provisioner, definition.Exec{Shell: firstInstance.ShellCommands})
		}

		instances = append(instances, instance)
	}

	return instances
}
