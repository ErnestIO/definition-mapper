/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
)

// MapFirewalls : Maps input firewalls to an ernest format ones
func MapFirewalls(payload input.Payload) (firewalls []output.Firewall, err error) {
	var fr output.FirewallRule

	for _, r := range payload.Service.Routers {
		if r.Rules == nil && payload.Service.Bootstrapping != "salt" {
			continue
		}

		f := output.Firewall{}
		f.Name = payload.Service.Name + "-" + r.Name
		f.RouterName = r.Name

		if payload.Service.Bootstrapping == "salt" {
			// Allow port 22 & 5985 from salt master to other networks for ssh/winrm
			fr = output.FirewallRule{}
			fr.SourceIP = "10.254.254.0/24"
			fr.SourcePort = "any"
			fr.DestinationIP = "any"
			fr.DestinationPort = "22"
			fr.Protocol = "tcp"
			f.Rules = append(f.Rules, fr)

			fr = output.FirewallRule{}
			fr.SourceIP = "10.254.254.0/24"
			fr.SourcePort = "any"
			fr.DestinationIP = "any"
			fr.DestinationPort = "5985"
			fr.Protocol = "tcp"
			f.Rules = append(f.Rules, fr)

			// Allow services/salt range to talk to DNS, minions to external Salt packages
			fr = output.FirewallRule{}
			fr.SourceIP = "internal"
			fr.SourcePort = "any"
			fr.DestinationIP = "external"
			fr.DestinationPort = "any"
			fr.Protocol = "any"
			f.Rules = append(f.Rules, fr)

			// Allow port 8000 to current ernest instance
			for _, ip := range payload.Service.ErnestIP {
				sw := false
				for _, rule := range f.Rules {
					if rule.SourceIP == ip {
						sw = true
					}
				}
				if sw == false {
					f.Rules = append(f.Rules, output.FirewallRule{
						SourceIP:        ip,
						SourcePort:      "any",
						DestinationIP:   "",
						DestinationPort: "8000",
						Protocol:        "tcp",
					})
				}
			}

			for _, network := range r.Networks {
				fnr := output.FirewallRule{}
				fnr.Name = network.Name + "-salt-firewall-4505-rule"
				fnr.SourceIP = network.Subnet
				fnr.SourcePort = "any"
				fnr.DestinationIP = "10.254.254.100"
				fnr.DestinationPort = "4505"
				fnr.Protocol = "tcp"
				f.Rules = append(f.Rules, fnr)
			}

			for _, network := range r.Networks {
				fnr := output.FirewallRule{}
				fnr.Name = network.Name + "-salt-firewall-4506-rule"
				fnr.SourceIP = network.Subnet
				fnr.SourcePort = "any"
				fnr.DestinationIP = "10.254.254.100"
				fnr.DestinationPort = "4506"
				fnr.Protocol = "tcp"
				f.Rules = append(f.Rules, fnr)
			}
		}

		// Validate Firewall Rules
		if r.Rules != nil {
			for _, rule := range r.Rules {
				// Check if firewall is valid
				if valid, err := rule.IsValid(r.Networks, payload.Service.Loadbalancers); valid == false {
					return firewalls, err
				}

				fr := output.FirewallRule{}
				fr.Name = rule.Name
				fr.SourceIP = rule.Source
				fr.SourcePort = rule.FromPort
				fr.DestinationIP = rule.Destination
				fr.DestinationPort = rule.ToPort
				fr.Protocol = rule.Protocol

				f.Rules = append(f.Rules, fr)
			}
		}

		firewalls = append(firewalls, f)
	}
	return firewalls, err
}

// HasChangedFirewalls : Checks if a firewall has changed
func HasChangedFirewalls(oldFirewalls []output.Firewall, newFirewalls []output.Firewall) bool {
	for _, o := range oldFirewalls {
		for _, n := range newFirewalls {
			if n.Name == o.Name {
				if HasChangedFirewallRules(o.Rules, n.Rules) {
					return true
				}
			}
		}
	}
	return false
}

// HasChangedFirewallRules : Checks if a firewall has any changed rules
func HasChangedFirewallRules(oldRules []output.FirewallRule, newRules []output.FirewallRule) bool {
	if len(newRules) != len(oldRules) {
		return true
	}
	for i, n := range newRules {
		if hasChangedDestinationIP(n.DestinationIP, oldRules[i].DestinationIP) ||
			n.DestinationPort != oldRules[i].DestinationPort ||
			n.Protocol != oldRules[i].Protocol ||
			n.SourceIP != oldRules[i].SourceIP ||
			n.SourcePort != oldRules[i].SourcePort {
			return true
		}
	}
	return false
}

func hasChangedDestinationIP(n, o string) bool {
	// In case the destination ip is empty it won't be empty on the previous
	// build as it's internally replaced by the endpoint
	if n == "" {
		return false
	}
	if n == o {
		return false
	}
	return true
}

// MapResultingFirewalls maps all input firewalls to an firewall output struct
func MapResultingFirewalls(result []output.Firewall) []output.Firewall {
	out := make([]output.Firewall, len(result))
	for i, firewall := range result {
		out[i].Rules = make([]output.FirewallRule, len(firewall.Rules))

		out[i].Name = firewall.Name
		out[i].RouterName = firewall.RouterName
		for x, rule := range firewall.Rules {
			out[i].Rules[x].DestinationIP = rule.DestinationIP
			out[i].Rules[x].DestinationPort = rule.DestinationPort
			out[i].Rules[x].SourceIP = rule.SourceIP
			out[i].Rules[x].SourcePort = rule.SourcePort
			out[i].Rules[x].Protocol = rule.Protocol

		}
	}
	return out
}
