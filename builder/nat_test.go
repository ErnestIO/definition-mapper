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

func TestNATSMapping(t *testing.T) {
	Convey("Given I'm mapping the input NATS", t, func() {
		networks := []input.Network{}
		networks = append(networks, input.Network{
			Name:   "foo",
			Subnet: "10.1.0.11/24",
		})
		ports := []input.PortForwarding{}
		ports = append(ports, input.PortForwarding{
			FromPort:    "80",
			Destination: "10.1.0.12",
			ToPort:      "80",
		})
		iRouters := []input.Router{}
		iRouters = append(iRouters, input.Router{
			Name:           "router",
			Networks:       networks,
			PortForwarding: ports,
		})
		p := input.Payload{
			Datacenter: input.Datacenter{
				ExternalNetwork: "foo",
			},
			Service: input.ServiceDefinition{
				Bootstrapping: "salt",
				Routers:       iRouters,
			},
		}

		Convey("When using a salt bootstrappig type", func() {
			r, err := MapNATS(p)

			Convey("Then nats to be created should > than provided ones", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
				So(len(r[0].Rules), ShouldEqual, 5)

				Convey("And salt nat rules are mapped", func() {
					So(r[0].Rules[0].Type, ShouldEqual, "dnat")
					So(r[0].Rules[0].OriginIP, ShouldEqual, "")
					So(r[0].Rules[0].OriginPort, ShouldEqual, "8000")
					So(r[0].Rules[0].TranslationIP, ShouldEqual, "10.254.254.100")
					So(r[0].Rules[0].TranslationPort, ShouldEqual, "8000")
					So(r[0].Rules[0].Protocol, ShouldEqual, "tcp")
					So(r[0].Rules[0].Network, ShouldEqual, "foo")

					So(r[0].Rules[1].Type, ShouldEqual, "dnat")
					So(r[0].Rules[1].OriginIP, ShouldEqual, "")
					So(r[0].Rules[1].OriginPort, ShouldEqual, "22")
					So(r[0].Rules[1].TranslationIP, ShouldEqual, "10.254.254.100")
					So(r[0].Rules[1].TranslationPort, ShouldEqual, "22")
					So(r[0].Rules[1].Protocol, ShouldEqual, "tcp")
					So(r[0].Rules[1].Network, ShouldEqual, "foo")

					So(r[0].Rules[2].Type, ShouldEqual, "snat")
					So(r[0].Rules[2].OriginIP, ShouldEqual, "10.254.254.0/24")
					So(r[0].Rules[2].OriginPort, ShouldEqual, "any")
					So(r[0].Rules[2].TranslationIP, ShouldEqual, "")
					So(r[0].Rules[2].TranslationPort, ShouldEqual, "any")
					So(r[0].Rules[2].Protocol, ShouldEqual, "any")
					So(r[0].Rules[2].Network, ShouldEqual, "foo")
				})

				Convey("And input nat rules are mapped", func() {
					So(r[0].Rules[3].Type, ShouldEqual, "snat")
					So(r[0].Rules[3].OriginIP, ShouldEqual, "10.1.0.11/24")
					So(r[0].Rules[3].OriginPort, ShouldEqual, "any")
					So(r[0].Rules[3].TranslationIP, ShouldEqual, "")
					So(r[0].Rules[3].TranslationPort, ShouldEqual, "any")
					So(r[0].Rules[3].Protocol, ShouldEqual, "any")
					So(r[0].Rules[3].Network, ShouldEqual, "foo")
				})

				Convey("And port forwarding nat generated rules are mapped", func() {
					So(r[0].Rules[4].Type, ShouldEqual, "dnat")
					So(r[0].Rules[4].OriginIP, ShouldEqual, "")
					So(r[0].Rules[4].OriginPort, ShouldEqual, "80")
					So(r[0].Rules[4].TranslationIP, ShouldEqual, "10.1.0.12")
					So(r[0].Rules[4].TranslationPort, ShouldEqual, "80")
					So(r[0].Rules[4].Protocol, ShouldEqual, "tcp")
					So(r[0].Rules[4].Network, ShouldEqual, "foo")
				})
			})
		})

		Convey("When using a default bootstrappig type", func() {
			p.Service.Bootstrapping = ""
			r, err := MapNATS(p)

			Convey("Then nats to be created should > than provided ones", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
				So(len(r[0].Rules), ShouldEqual, 2)

				Convey("And input nat rules are mapped", func() {
					So(r[0].Rules[0].Type, ShouldEqual, "snat")
					So(r[0].Rules[0].OriginIP, ShouldEqual, "10.1.0.11/24")
					So(r[0].Rules[0].OriginPort, ShouldEqual, "any")
					So(r[0].Rules[0].TranslationIP, ShouldEqual, "")
					So(r[0].Rules[0].TranslationPort, ShouldEqual, "any")
					So(r[0].Rules[0].Protocol, ShouldEqual, "any")
					So(r[0].Rules[0].Network, ShouldEqual, "foo")
				})

				Convey("And port forwarding nat generated rules are mapped", func() {
					So(r[0].Rules[1].Type, ShouldEqual, "dnat")
					So(r[0].Rules[1].OriginIP, ShouldEqual, "")
					So(r[0].Rules[1].OriginPort, ShouldEqual, "80")
					So(r[0].Rules[1].TranslationIP, ShouldEqual, "10.1.0.12")
					So(r[0].Rules[1].TranslationPort, ShouldEqual, "80")
					So(r[0].Rules[1].Protocol, ShouldEqual, "tcp")
					So(r[0].Rules[1].Network, ShouldEqual, "foo")
				})
			})
		})
	})
}

func TestHasChangedNats(t *testing.T) {
	Convey("Given I'm mapping the input NATS", t, func() {
		o := []output.Nat{}
		n := []output.Nat{}
		oRules := []output.NatRule{}
		nRules := []output.NatRule{}
		o = append(o, output.Nat{
			Name:  "foo",
			Rules: oRules,
		})
		n = append(n, output.Nat{
			Name:  "foo",
			Rules: nRules,
		})

		Convey("When nothing has changed", func() {
			r := HasChangedNats(o, n)
			Convey("Then should return false", func() {
				So(r, ShouldEqual, false)
			})
		})

		Convey("When we change a rule", func() {
			o[0].Rules = append(o[0].Rules, output.NatRule{
				Network:  "10.1.0.11",
				OriginIP: "1.1.1.1",
			})
			n[0].Rules = append(n[0].Rules, output.NatRule{
				Network:  "10.1.0.11",
				OriginIP: "1.1.1.2",
			})
			r := HasChangedNats(o, n)
			Convey("Then should return false", func() {
				So(r, ShouldEqual, true)
			})
		})
	})
}
