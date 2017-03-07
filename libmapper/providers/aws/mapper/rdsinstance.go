/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapRDSInstances : Maps the rds instances for the input payload on a ernest internal format
func MapRDSInstances(d *definition.Definition) []*components.RDSInstance {
	var instances []*components.RDSInstance

	for _, instance := range d.RDSInstances {

		i := &components.RDSInstance{
			Name:              instance.Name,
			Size:              instance.Size,
			Engine:            instance.Engine,
			EngineVersion:     instance.EngineVersion,
			Port:              instance.Port,
			Cluster:           instance.Cluster,
			Public:            instance.Public,
			MultiAZ:           instance.MultiAZ,
			PromotionTier:     instance.PromotionTier,
			StorageType:       instance.Storage.Type,
			StorageSize:       instance.Storage.Size,
			StorageIops:       instance.Storage.Iops,
			AvailabilityZone:  instance.AvailabilityZone,
			SecurityGroups:    instance.SecurityGroups,
			Networks:          instance.Networks,
			DatabaseName:      instance.DatabaseName,
			DatabaseUsername:  instance.DatabaseUsername,
			DatabasePassword:  instance.DatabasePassword,
			AutoUpgrade:       instance.AutoUpgrade,
			BackupRetention:   instance.Backups.Retention,
			BackupWindow:      instance.Backups.Window,
			MaintenanceWindow: instance.MaintenanceWindow,
			ReplicationSource: instance.ReplicationSource,
			FinalSnapshot:     instance.FinalSnapshot,
			License:           instance.License,
			Timezone:          instance.Timezone,
			Tags:              mapTagsServiceOnly(d.Name),
		}

		for _, cluster := range d.RDSClusters {
			if cluster.Name == i.Cluster {
				i.Engine = cluster.Engine
			}
		}

		i.SetDefaultVariables()

		instances = append(instances, i)
	}
	return instances
}

// MapDefinitionRDSInstances : Maps the rds instances from the internal format to the input definition format
func MapDefinitionRDSInstances(g *graph.Graph) []definition.RDSInstance {
	var instances []definition.RDSInstance

	for _, gi := range g.GetComponents().ByType("rds_instance") {
		instance := gi.(*components.RDSInstance)

		i := definition.RDSInstance{
			Name:              instance.Name,
			Size:              instance.Size,
			Engine:            instance.Engine,
			EngineVersion:     instance.EngineVersion,
			Port:              instance.Port,
			Cluster:           instance.Cluster,
			Public:            instance.Public,
			MultiAZ:           instance.MultiAZ,
			PromotionTier:     instance.PromotionTier,
			AvailabilityZone:  instance.AvailabilityZone,
			SecurityGroups:    instance.SecurityGroups,
			Networks:          instance.Networks,
			DatabaseName:      instance.DatabaseName,
			DatabaseUsername:  instance.DatabaseUsername,
			DatabasePassword:  instance.DatabasePassword,
			AutoUpgrade:       instance.AutoUpgrade,
			MaintenanceWindow: instance.MaintenanceWindow,
			ReplicationSource: instance.ReplicationSource,
			FinalSnapshot:     instance.FinalSnapshot,
			License:           instance.License,
			Timezone:          instance.Timezone,
		}

		i.Storage.Type = instance.StorageType
		i.Storage.Size = instance.StorageSize
		i.Storage.Iops = instance.StorageIops
		i.Backups.Retention = instance.BackupRetention
		i.Backups.Window = instance.BackupWindow

		if i.Storage.Type != "io1" {
			i.Storage.Iops = nil
		}

		instances = append(instances, i)
	}

	return instances
}
