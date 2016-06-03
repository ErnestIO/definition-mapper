/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"testing"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNetworksMapping(t *testing.T) {
	Convey("Given I try to map a network", t, func() {
		networks := []input.Network{}
		networks = append(networks, input.Network{
			Name:   "bar",
			Subnet: "10.0.0.10/24",
		})
		iRouters := []input.Router{}
		iRouters = append(iRouters, input.Router{
			Name:     "router",
			Networks: networks,
		})
		p := input.Payload{
			Service: input.ServiceDefinition{
				Name:          "foo",
				Bootstrapping: "salt",
				Routers:       iRouters,
			},
			Datacenter: input.Datacenter{Name: "datacenter"},
		}
		m := output.FSMMessage{
			ServiceName: "foo",
		}
		routers := []output.Router{}
		m.Routers.Items = append(routers, output.Router{Name: "router"})

		Convey("When the input specifies duplicated network names", func() {
			Convey("And the network name is empty", func() {
				p.Service.Routers[0].Networks = append(p.Service.Routers[0].Networks, input.Network{
					Name:   "bar",
					Subnet: "10.0.0.10/24",
				})
				_, err := MapNetworks(p, m)
				Convey("Then we should receive a controlled error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Input includes duplicate network names")
				})
			})
		})

		Convey("When the input specifies invalid networks", func() {
			Convey("And the network name is empty", func() {
				p.Service.Routers[0].Networks[0].Name = ""
				_, err := MapNetworks(p, m)
				Convey("Then we should receive a controlled error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Network name should not be null")
				})
			})

			Convey("And the network name is > 50 chars", func() {
				p.Service.Routers[0].Networks[0].Name = "a_very_big_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
				_, err := MapNetworks(p, m)
				Convey("Then we should receive a controlled error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Network name can't be greater than 50 characters")
				})
			})

			Convey("And the network subnet is not valid CIDR", func() {
				p.Service.Routers[0].Networks[0].Subnet = ""
				_, err := MapNetworks(p, m)
				Convey("Then we should receive a controlled error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Network CIDR is not valid")
				})
			})

		})

		Convey("When the input specifies bootstrap as not salt", func() {
			p.Service.Bootstrapping = "foo"
			n, err := MapNetworks(p, m)
			Convey("Then only input netoworks should be mapped", func() {
				So(err, ShouldBeNil)
				So(len(n), ShouldEqual, len(p.Service.Routers[0].Networks))
				So(n[0].Name, ShouldEqual, "datacenter-foo-bar")
				So(n[0].RouterName, ShouldEqual, m.Routers.Items[0].Name)
				So(n[0].Subnet, ShouldEqual, p.Service.Routers[0].Networks[0].Subnet)
			})
		})

		Convey("When the input specifies bootstrap as salt", func() {
			n, err := MapNetworks(p, m)
			Convey("Then an extra network should be created", func() {
				So(err, ShouldBeNil)
				So(len(n), ShouldEqual, len(p.Service.Routers[0].Networks)+1)
				So(n[0].Name, ShouldEqual, "datacenter-foo-salt")
				So(n[0].RouterName, ShouldEqual, m.Routers.Items[0].Name)
				So(n[0].Subnet, ShouldEqual, "10.254.254.0/24")
			})
			Convey("Then input network should be mapped as usual", func() {
				So(err, ShouldBeNil)
				So(len(n), ShouldEqual, len(p.Service.Routers[0].Networks)+1)
				So(n[1].Name, ShouldEqual, "datacenter-foo-bar")
				So(n[1].RouterName, ShouldEqual, m.Routers.Items[0].Name)
				So(n[1].Subnet, ShouldEqual, p.Service.Routers[0].Networks[0].Subnet)
			})
		})
	})
}

func TestNetworksToCreateMapping(t *testing.T) {
	Convey("Given I try to map the networks to be created", t, func() {
		oNetworks := []output.Network{}
		oNetworks = append(oNetworks, output.Network{Name: "foo"})
		o := output.FSMMessage{}
		o.Networks.Items = oNetworks
		nNetworks := []output.Network{}
		nNetworks = append(nNetworks, output.Network{Name: "foo"})
		n := output.FSMMessage{}
		n.Networks.Items = nNetworks

		Convey("When we're keeping the same number of networks", func() {
			i, err := MapNetworksToCreate(&o, n)
			Convey("Then we should receive an empty list", func() {
				So(err, ShouldBeNil)
				So(len(i), ShouldEqual, 0)
			})
		})

		Convey("When we're increasing the number of networks", func() {
			nNetworks = append(nNetworks, output.Network{Name: "bar"})
			n := output.FSMMessage{}
			n.Networks.Items = nNetworks
			i, err := MapNetworksToCreate(&o, n)
			Convey("Then we should receive a list of new networks", func() {
				So(err, ShouldBeNil)
				So(len(i), ShouldEqual, 1)
				So(i[0].Name, ShouldEqual, "bar")
			})
		})

		Convey("When we're adding a new network 'group' with the name of an existing one", func() {
			nNetworks = append(nNetworks, output.Network{Name: "bar"})
			nNetworks = append(nNetworks, output.Network{Name: "bar"})
			n := output.FSMMessage{}
			n.Networks.Items = nNetworks
			_, err := MapNetworksToCreate(&o, n)
			Convey("Then we should receive a list of new networks", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Input includes duplicate network names")
			})
		})

		Convey("When we're changing the name of a network", nil)
		Convey("When we're decreasing the number of networks", nil)
	})
}
