/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapSQLServers ...
func MapSQLServers(d *definition.Definition) (ips []*components.SQLServer) {
	for _, rg := range d.ResourceGroups {
		for _, ss := range rg.SQLServers {
			n := &components.SQLServer{}
			n.Name = ss.Name
			n.Version = ss.Version
			n.AdministratorLogin = ss.AdministratorLogin
			n.AdministratorLoginPassword = ss.AdministratorLoginPassword
			n.ResourceGroupName = rg.Name
			n.Location = rg.Location
			n.Tags = mapTags(ss.Name, d.Name)
			for k, v := range ss.Tags {
				n.Tags[k] = v
			}

			if n.ID != "" {
				n.SetAction("none")
			}

			n.SetDefaultVariables()

			ips = append(ips, n)
		}
	}

	return
}

// MapDefinitionSQLServers : ...
func MapDefinitionSQLServers(g *graph.Graph, rg *definition.ResourceGroup) (ss []definition.SQLServer) {
	for _, c := range g.GetComponents().ByType("sql_server") {
		sqls := c.(*components.SQLServer)

		if sqls.ResourceGroupName != rg.Name {
			continue
		}

		dsqls := definition.SQLServer{
			ID:                         sqls.GetProviderID(),
			Name:                       sqls.Name,
			Version:                    sqls.Version,
			AdministratorLogin:         sqls.AdministratorLogin,
			AdministratorLoginPassword: sqls.AdministratorLoginPassword,
		}

		for _, cd := range g.GetComponents().ByType("sql_database") {
			sqld := cd.(*components.SQLDatabase)

			if sqld.ResourceGroupName != rg.Name && sqld.ServerName != dsqls.Name {
				continue
			}

			dsqld := definition.SQLDatabase{
				ID:                            sqld.GetProviderID(),
				Name:                          sqld.Name,
				CreateMode:                    sqld.CreateMode,
				SourceDatabaseID:              sqld.SourceDatabaseID,
				RestorePointInTime:            sqld.RestorePointInTime,
				Edition:                       sqld.Edition,
				Collation:                     sqld.Collation,
				MaxSizeBytes:                  sqld.MaxSizeBytes,
				RequestedServiceObjectiveID:   sqld.RequestedServiceObjectiveID,
				RequestedServiceObjectiveName: sqld.RequestedServiceObjectiveName,
				SourceDatabaseDeletionData:    sqld.SourceDatabaseDeletionData,
			}

			dsqls.Databases = append(dsqls.Databases, dsqld)
		}

		ss = append(ss, dsqls)
	}

	return
}
