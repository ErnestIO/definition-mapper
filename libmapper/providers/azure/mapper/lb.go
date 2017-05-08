/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	ernestiolb "github.com/ernestio/ernestprovider/providers/azure/lb"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapLBs ...
func MapLBs(d *definition.Definition) (lbs []*components.LB) {
	for _, rg := range d.ResourceGroups {
		for _, lb := range rg.LBs {
			n := &components.LB{}
			n.Name = lb.Name
			n.ResourceGroupName = rg.Name
			n.Location = rg.Location
			for _, d := range lb.FrontendIPConfigurations {
				n.FrontendIPConfigurations = append(n.FrontendIPConfigurations, ernestiolb.FrontendIPConfiguration{
					Name:            d.Name,
					SubnetID:        d.Subnet,
					PublicIPAddress: d.PublicIPAddress,
				})
			}

			n.Tags = mapTags(lb.Name, d.Name)
			if lb.ID != "" {
				n.SetAction("none")
			}

			n.SetDefaultVariables()

			lbs = append(lbs, n)
		}
	}

	return
}

// MapDefinitionLBs : ...
func MapDefinitionLBs(g *graph.Graph, rg *definition.ResourceGroup) (lbs []definition.LB) {
	for _, c := range g.GetComponents().ByType("lb") {
		lb := c.(*components.LB)

		if lb.ResourceGroupName != rg.Name {
			continue
		}

		nLB := definition.LB{
			ID:       lb.GetProviderID(),
			Name:     lb.Name,
			Location: lb.Location,
		}

		lbs = append(lbs, nLB)
	}

	return
}