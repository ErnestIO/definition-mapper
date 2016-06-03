/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPortForwardingIsValidDestination(t *testing.T) {
	r := PortForwarding{Destination: "127.0.0.1", FromPort: "any", ToPort: "any"}
	Convey("Given an existing valid nat", t, func() {
		natRules := make([]NatRule, 1)
		natRules[0] = NatRule{OriginPort: "foo", TranslationPort: "tPort", TranslationIP: "tIP"}
		nat := Nat{Rules: natRules}

		Convey("When I try to validate a rule with non ip destination", func() {
			r.Destination = "foo"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding must be a valid IP")
			})
		})

		Convey("When I try to validate a rule with a valid destination", func() {
			r.Destination = "127.0.0.1"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, true)
				So(err, ShouldBeNil)
			})
		})

	})
}

func TestPortForwardingIsValidFromPort(t *testing.T) {
	r := PortForwarding{Destination: "127.0.0.1", FromPort: "any", ToPort: "any"}
	Convey("Given an existing valid network", t, func() {
		natRules := make([]NatRule, 1)
		natRules[0] = NatRule{OriginPort: "foo", TranslationPort: "tPort", TranslationIP: "tIP"}
		nat := Nat{Rules: natRules}

		Convey("When From Port is any", func() {
			r.FromPort = "any"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, true)
				So(err, ShouldBeNil)
			})
		})

		Convey("When From Port is not any and not numeric", func() {
			r.FromPort = "foo"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port is not valid")
			})
		})

		Convey("When From Port is not any and not in range", func() {
			r.FromPort = "0"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port is out of range [1 - 65535]")
			})
		})

		Convey("When From Port is not any and great than range", func() {
			r.FromPort = "65536"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding From Port is out of range [1 - 65535]")
			})
		})
	})
}

func TestPortForwardingIsValidToPort(t *testing.T) {
	r := PortForwarding{Destination: "127.0.0.1", FromPort: "any", ToPort: "any"}
	Convey("Given an existing valid network", t, func() {
		natRules := make([]NatRule, 1)
		natRules[0] = NatRule{OriginPort: "foo", TranslationPort: "tPort", TranslationIP: "tIP"}
		nat := Nat{Rules: natRules}

		Convey("When To Port is any", func() {
			r.ToPort = "any"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, true)
				So(err, ShouldBeNil)
			})
		})

		Convey("When From Port is not any and not numeric", func() {
			r.ToPort = "foo"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding To Port is not valid")
			})
		})

		Convey("When To Port is not any and not in range", func() {
			r.ToPort = "0"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding To Port is out of range [1 - 65535]")
			})
		})

		Convey("When To Port is not any and great than range", func() {
			r.ToPort = "65536"
			valid, err := r.IsValid(nat)
			Convey("Then should return an error", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Port Forwarding To Port is out of range [1 - 65535]")
			})
		})
	})
}
