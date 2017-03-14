/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/aws/definition"
)

// MapInternetGateways ...
func MapInternetGateways(d *definition.Definition) []*components.InternetGateway {
	var igs []*components.InternetGateway
	var vpcs []string

	for _, network := range d.Networks {
		if network.Public {
			vpcs = appendUnique(vpcs, network.VPC)
		}
	}

	for _, vpc := range vpcs {
		ig := &components.InternetGateway{
			Name: vpc,
			Vpc:  vpc,
			Tags: mapTags(vpc, d.Name),
		}

		ig.SetDefaultVariables()

		igs = append(igs, ig)
	}

	return igs
}

func appendUnique(s []string, v string) []string {
	for _, x := range s {
		if x == v {
			return s
		}
	}
	return append(s, v)
}
