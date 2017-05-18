/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapSQLDatabases ...
func MapSQLDatabases(d *definition.Definition) (ips []*components.SQLDatabase) {
	for _, rg := range d.ResourceGroups {
		for _, ss := range rg.SQLServers {
			for _, sd := range ss.Databases {
				n := &components.SQLDatabase{}
				n.Name = sd.Name
				n.ResourceGroupName = rg.Name
				n.Location = rg.Location
				n.ServerName = ss.Name
				n.CreateMode = sd.CreateMode
				n.SourceDatabaseID = sd.SourceDatabaseID
				n.RestorePointInTime = sd.RestorePointInTime
				n.Edition = sd.Edition
				n.Collation = sd.Collation
				n.MaxSizeBytes = sd.MaxSizeBytes
				n.RequestedServiceObjectiveID = sd.RequestedServiceObjectiveID
				n.RequestedServiceObjectiveName = sd.RequestedServiceObjectiveName
				n.SourceDatabaseDeletionData = sd.SourceDatabaseDeletionData
				n.ElasticPoolName = sd.ElasticPoolName
				n.Tags = mapTags(sd.Name, d.Name)
				for k, v := range sd.Tags {
					n.Tags[k] = v
				}

				if n.ID != "" {
					n.SetAction("none")
				}

				n.SetDefaultVariables()

				ips = append(ips, n)
			}
		}
	}

	return
}
