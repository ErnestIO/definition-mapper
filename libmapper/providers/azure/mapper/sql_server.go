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
			n := components.SQLServer{}
			n.Name = ss.Name
			n.Version = ss.Version
			n.AdministratorLogin = ss.AdministratorLogin
			n.AdministratorLoginPassword = ss.AdministratorLoginPassword
			n.ResourceGroupName = rg.Name
			n.Location = rg.Location
			n.Tags = mapTags(ss.Name, d.Name)

			if n.ID != "" {
				n.SetAction("none")
			}

			n.SetDefaultVariables()

			ips = append(ips, &n)
		}
	}

	return
}

// MapDefinitionSQLServers : ...
func MapDefinitionSQLServers(g *graph.Graph, rg *definition.ResourceGroup) (ss []definition.SQLServer) {
	for _, c := range g.GetComponents().ByType("sql_server") {
		sql := c.(*components.SQLServer)

		if sql.ResourceGroupName != rg.Name {
			continue
		}

		n := definition.SQLServer{
			ID:                         sql.GetProviderID(),
			Name:                       sql.Name,
			Version:                    sql.Version,
			AdministratorLogin:         sql.AdministratorLogin,
			AdministratorLoginPassword: sql.AdministratorLoginPassword,
		}

		ss = append(ss, n)
	}

	return
}
