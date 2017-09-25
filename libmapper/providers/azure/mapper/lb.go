/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/azure/definition"
	lb "github.com/ernestio/ernestprovider/providers/azure/lb"
	graph "gopkg.in/r3labs/graph.v2"
)

// MapLBs ...
func MapLBs(d *definition.Definition) (lbs []*components.LB) {
	for _, rg := range d.ResourceGroups {
		for _, loadbalancer := range rg.LBs {
			n := &components.LB{}
			n.Name = loadbalancer.Name
			n.ResourceGroupName = rg.Name
			n.Location = rg.Location
			for _, d := range loadbalancer.FrontendIPConfigurations {
				lbconfig := lb.FrontendIPConfiguration{
					Name:                       d.Name,
					Subnet:                     d.Subnet,
					PrivateIPAddress:           d.PrivateIPAddress,
					PrivateIPAddressAllocation: d.PrivateIPAddressAllocation,
				}

				if d.PublicIPAddressAllocation != "" {
					lbconfig.PublicIPAddress = n.Name + "-" + d.Name
				}

				n.FrontendIPConfigurations = append(n.FrontendIPConfigurations, lbconfig)
			}

			n.Tags = mapTags(loadbalancer.Name, d.Name)

			for k, v := range loadbalancer.Tags {
				n.Tags[k] = v
			}

			if loadbalancer.ID != "" {
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

		dlb := definition.LB{
			ID:       lb.GetProviderID(),
			Name:     lb.Name,
			Location: lb.Location,
		}

		for _, config := range lb.FrontendIPConfigurations {
			dconfig := definition.FrontendIPConfiguration{
				Name:                       config.Name,
				PrivateIPAddress:           config.PrivateIPAddress,
				PrivateIPAddressAllocation: config.PrivateIPAddressAllocation,
				Subnet: config.Subnet,
			}
			if config.PublicIPAddressID != "" {
				cpip := g.GetComponents().ByProviderID(config.PublicIPAddressID)
				if cpip != nil {
					pip := cpip.(*components.PublicIP)
					dconfig.PublicIPAddressAllocation = pip.PublicIPAddressAllocation
				}
			}

			for _, cn := range g.GetComponents().ByType("lb_rule") {
				lg := cn.(*components.LBRule)
				dconfig.Rules = append(dconfig.Rules, definition.LoadbalancerRule{
					Name:               lg.Name,
					Protocol:           lg.Protocol,
					FrontendPort:       lg.FrontendPort,
					BackendPort:        lg.BackendPort,
					BackendAddressPool: lg.BackendAddressPool,
					Probe:              lg.Probe,
					FloatingIP:         lg.EnableFloatingIP,
					IdleTimeout:        lg.IdleTimeoutInMinutes,
					LoadDistribution:   lg.LoadDistribution,
				})
			}

			dlb.FrontendIPConfigurations = append(dlb.FrontendIPConfigurations, dconfig)
		}

		for _, cn := range g.GetComponents().ByType("lb_backend_address_pool") {
			lg := cn.(*components.LBBackendAddressPool)
			dlb.BackendAddressPools = append(dlb.BackendAddressPools, lg.Name)
		}

		for _, cn := range g.GetComponents().ByType("lb_probe") {
			lg := cn.(*components.LBProbe)
			dlb.Probes = append(dlb.Probes, definition.LoadbalancerProbe{
				Name:            lg.Name,
				Port:            lg.Port,
				Protocol:        lg.Protocol,
				RequestPath:     lg.RequestPath,
				Interval:        lg.IntervalInSeconds,
				MaximumFailures: lg.NumberOfProbes,
			})
		}

		lbs = append(lbs, dlb)
	}

	return
}
