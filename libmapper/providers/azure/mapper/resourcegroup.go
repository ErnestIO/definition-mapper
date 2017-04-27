/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapResourceGroups ...
func MapResourceGroups(d *definition.Definition) (groups []*components.ResourceGroup) {
	for _, group := range d.ResourceGroups {
		cv := components.ResourceGroup{}
		cv.Name = group.Name
		cv.Location = group.Location
		cv.Tags = mapTags(group.Name, d.Name)

		if group.ID != "" {
			cv.SetAction("none")
		}

		cv.SetDefaultVariables()

		groups = append(groups, &cv)
	}

	return
}

// MapDefinitionResourceGroups : ...
func MapDefinitionResourceGroups(g *graph.Graph) (rgs []definition.ResourceGroup) {
	for _, c := range g.GetComponents().ByType("resource_group") {
		rg := c.(*components.ResourceGroup)

		rgs = append(rgs, definition.ResourceGroup{
			ID:       rg.GetProviderID(),
			Name:     rg.Name,
			Location: rg.Location,
		})
	}

	return
}
