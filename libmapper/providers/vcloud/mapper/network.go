/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package mapper

import (
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/components"
	"github.com/ernestio/definition-mapper/libmapper/providers/vcloud/definition"
)

// MapNetworks : Maps the networks from a given input payload. Additionally it will create an
// extra salt-master network if the service bootstrapping mode is defined as salt
func MapNetworks(d *definition.Definition) []*components.Network {
	var networks []*components.Network

	for _, r := range d.Gateways {
		for _, network := range r.Networks {
			octets := getIPOctets(network.Subnet)

			n := &components.Network{
				Name:         network.Name,
				Subnet:       network.Subnet,
				StartAddress: octets + ".5",
				EndAddress:   octets + ".250",
				Gateway:      octets + ".1",
				Netmask:      parseNetmask(network.Subnet),
				DNS:          network.DNS,
				EdgeGateway:  r.Name,
			}

			n.SetDefaultVariables()

			networks = append(networks, n)
		}
	}

	return networks
}

func getIPOctets(rng string) string {
	// Splits the network range and returns the first three octets
	ip, _, err := net.ParseCIDR(rng)
	if err != nil {
		log.Println(err)
	}
	octets := strings.Split(ip.String(), ".")
	octets = append(octets[:3], octets[3+1:]...)
	octetString := strings.Join(octets, ".")
	return octetString
}

func parseNetmask(rng string) string {
	// Convert netmask hex to string, generated from network range CIDR
	_, nw, _ := net.ParseCIDR(rng)
	hx, _ := hex.DecodeString(nw.Mask.String())
	netmask := fmt.Sprintf("%v.%v.%v.%v", hx[0], hx[1], hx[2], hx[3])
	return netmask
}
