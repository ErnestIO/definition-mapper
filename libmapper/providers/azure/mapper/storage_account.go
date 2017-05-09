/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapStorageAccounts ...
func MapStorageAccounts(d *definition.Definition) (ips []*components.StorageAccount) {
	for _, rg := range d.ResourceGroups {
		for _, sa := range rg.StorageAccounts {
			n := &components.StorageAccount{}
			n.Name = sa.Name
			n.ResourceGroupName = rg.Name
			n.Location = rg.Location
			n.AccountKind = sa.AccountKind
			n.AccountType = sa.AccountType
			n.EnableBlobEncryption = sa.EnableBlobEncryption
			n.Tags = mapTags(sa.Name, d.Name)

			if n.ID != "" {
				n.SetAction("none")
			}

			n.SetDefaultVariables()

			ips = append(ips, n)
		}
	}

	return
}

// MapDefinitionStorageAccounts : ...
func MapDefinitionStorageAccounts(g *graph.Graph, rg *definition.ResourceGroup) (sa []definition.StorageAccount) {
	for _, c := range g.GetComponents().ByType("storage_account") {
		ca := c.(*components.StorageAccount)

		if ca.ResourceGroupName != rg.Name {
			continue
		}

		daccount := definition.StorageAccount{
			ID:                   ca.GetProviderID(),
			Name:                 ca.Name,
			AccountKind:          ca.AccountKind,
			AccountType:          ca.AccountType,
			EnableBlobEncryption: ca.EnableBlobEncryption,
		}

		for _, cx := range g.GetComponents().ByType("storage_container") {
			cc := cx.(*components.StorageContainer)

			if cc.ResourceGroupName != rg.Name && cc.StorageAccountName != daccount.Name {
				continue
			}

			dcontainer := definition.StorageContainer{
				ID:         cc.GetProviderID(),
				Name:       cc.Name,
				AccessType: cc.ContainerAccessType,
			}

			daccount.Containers = append(daccount.Containers, dcontainer)
		}

		sa = append(sa, daccount)
	}

	return
}
