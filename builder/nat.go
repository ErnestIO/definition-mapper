/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
)

// MapNATS : Generates necessary nats rules for input networks + salt if
// required
func MapNATS(payload input.Payload) (nats []output.Nat, err error) {
	for _, r := range payload.Service.Routers {
		if r.Networks == nil {
			continue
		}

		// Generate Nats
		nt := output.Nat{}

		nt.Name = payload.Service.Name
		nt.RouterName = r.Name
		nt.Service = payload.Service.Name
		nt.Rules = saltNATSRules(payload)

		// All Outbound Nat rules for networks
		for _, network := range r.Networks {
			rule := output.NatRule{
				Type:            "snat",
				OriginIP:        network.Subnet,
				OriginPort:      "any",
				TranslationIP:   "",
				TranslationPort: "any",
				Protocol:        "any",
				Network:         payload.Datacenter.ExternalNetwork,
			}
			nt.Rules = append(nt.Rules, rule)
		}

		for _, rule := range r.PortForwarding {
			ntr := output.NatRule{}
			ntr.Type = "dnat"
			ntr.OriginIP = rule.Source
			ntr.OriginPort = rule.FromPort
			ntr.TranslationIP = rule.Destination
			ntr.TranslationPort = rule.ToPort
			ntr.Protocol = "tcp"
			ntr.Network = payload.Datacenter.ExternalNetwork

			nt.Rules = append(nt.Rules, ntr)
		}

		nats = append(nats, nt)
	}
	return nats, err
}

// HasChangedNats : Verifies if  nats has been changed on the new input
func HasChangedNats(oldNats []output.Nat, newNats []output.Nat) bool {
	for _, o := range oldNats {
		for _, n := range newNats {
			if n.Name == o.Name {
				if hasChangedNatRules(o.Rules, n.Rules) {
					return true
				}
			}
		}
	}
	return false
}

// Check if a nat rules for the same network has changed
func hasChangedNatRules(oldRules []output.NatRule, newRules []output.NatRule) bool {
	if len(oldRules) != len(newRules) {
		return true
	}
	for i, n := range newRules {
		if n.Network != oldRules[i].Network ||
			hasChangedIP(n.OriginIP, oldRules[i].OriginIP) ||
			n.OriginPort != oldRules[i].OriginPort ||
			hasChangedIP(n.TranslationIP, oldRules[i].TranslationIP) ||
			n.TranslationPort != oldRules[i].TranslationPort ||
			n.Protocol != oldRules[i].Protocol ||
			n.Type != oldRules[i].Type {
			return true
		}
	}
	return false
}

func hasChangedIP(n, o string) bool {
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

// Generatest salt specific nats rules
func saltNATSRules(payload input.Payload) []output.NatRule {
	var ntr output.NatRule
	var rules []output.NatRule

	if payload.Service.IsSaltBootstrapped() == false {
		return rules
	}

	ntr = output.NatRule{}
	ntr.Type = "dnat"
	ntr.OriginIP = "" // TODO: Should be 'any' ?
	ntr.OriginPort = "8000"
	ntr.TranslationIP = "10.254.254.100"
	ntr.TranslationPort = "8000"
	ntr.Protocol = "tcp"
	ntr.Network = payload.Datacenter.ExternalNetwork

	rules = append(rules, ntr)

	ntr = output.NatRule{}
	ntr.Type = "dnat"
	ntr.OriginIP = "" // TODO: hould be 'any' ?
	ntr.OriginPort = "22"
	ntr.TranslationIP = "10.254.254.100"
	ntr.TranslationPort = "22"
	ntr.Protocol = "tcp"
	ntr.Network = payload.Datacenter.ExternalNetwork

	rules = append(rules, ntr)

	ntr = output.NatRule{}
	ntr.Type = "snat"
	ntr.OriginIP = "10.254.254.0/24"
	ntr.OriginPort = "any"
	ntr.TranslationIP = ""
	ntr.TranslationPort = "any"
	ntr.Protocol = "any"
	ntr.Network = payload.Datacenter.ExternalNetwork

	rules = append(rules, ntr)

	return rules
}

// MapResultingNats maps all input nats to an nat output struct
func MapResultingNats(result []output.Nat) []output.Nat {
	out := make([]output.Nat, len(result))
	for i, nat := range result {
		out[i].Rules = make([]output.NatRule, len(nat.Rules))

		out[i].Name = nat.Name
		out[i].RouterName = nat.RouterName
		for x, rule := range nat.Rules {
			out[i].Rules[x].Type = rule.Type
			out[i].Rules[x].Protocol = rule.Protocol
			out[i].Rules[x].Network = rule.Network
			out[i].Rules[x].OriginIP = rule.OriginIP
			out[i].Rules[x].OriginPort = rule.OriginPort
			out[i].Rules[x].TranslationIP = rule.TranslationIP
			out[i].Rules[x].TranslationPort = rule.TranslationPort
		}
	}
	return out
}
