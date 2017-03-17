/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapRoute53Zones : Maps the zones from a given input payload.
func MapRoute53Zones(d *definition.Definition) []*components.Route53Zone {
	var zones []*components.Route53Zone

	for _, zone := range d.Route53Zones {
		z := &components.Route53Zone{
			Name:    zone.Name,
			Private: zone.Private,
			Vpc:     zone.Vpc,
			Tags:    mapTagsServiceOnly(d.Name),
		}

		for _, record := range zone.Records {
			r := components.Record{
				Entry:         record.Entry,
				Type:          record.Type,
				Instances:     record.Instances,
				Loadbalancers: record.Loadbalancers,
				RDSClusters:   record.RDSClusters,
				RDSInstances:  record.RDSInstances,
				Values:        record.Values,
				TTL:           record.TTL,
			}

			z.Records = append(z.Records, r)
		}

		z.SetDefaultVariables()

		zones = append(zones, z)
	}

	return zones
}

// MapDefinitionRoute53Zones : Maps zones from the internal format to the input definition format
func MapDefinitionRoute53Zones(g *graph.Graph) []definition.Route53Zone {
	var zones []definition.Route53Zone

	for _, gzone := range g.GetComponents().ByType("route53") {
		zone := gzone.(*components.Route53Zone)

		z := definition.Route53Zone{
			Name:    zone.Name,
			Private: zone.Private,
		}

		for _, record := range zone.Records {
			z.Records = append(z.Records, definition.Record{
				Entry:         record.Entry,
				Type:          record.Type,
				TTL:           record.TTL,
				Instances:     record.Instances,
				Loadbalancers: record.Loadbalancers,
				RDSClusters:   record.RDSClusters,
				RDSInstances:  record.RDSInstances,
				Values:        record.Values,
			})
		}

		zones = append(zones, z)
	}

	return zones
}
