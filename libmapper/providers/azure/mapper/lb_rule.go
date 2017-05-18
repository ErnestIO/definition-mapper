/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
)

// MapLBRules ...
func MapLBRules(d *definition.Definition) (rules []*components.LBRule) {
	for _, rg := range d.ResourceGroups {
		for _, loadbalancer := range rg.LBs {
			for _, config := range loadbalancer.FrontendIPConfigurations {
				for _, rule := range config.Rules {
					n := &components.LBRule{}
					n.Name = rule.Name
					n.Probe = rule.Probe
					n.Protocol = rule.Protocol
					n.FrontendPort = rule.FrontendPort
					n.BackendPort = rule.BackendPort
					n.EnableFloatingIP = rule.FloatingIP
					n.IdleTimeoutInMinutes = rule.IdleTimeout
					n.LoadDistribution = rule.LoadDistribution
					n.Loadbalancer = loadbalancer.Name
					n.ResourceGroupName = rg.Name
					n.FrontendIPConfigurationName = config.Name

					n.SetDefaultVariables()

					rules = append(rules, n)
				}
			}
		}
	}

	return
}
