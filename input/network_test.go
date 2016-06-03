/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsValidCIDR(t *testing.T) {
	n := Network{Name: "foo", Router: "10.11.1.1", Subnet: "10.11.1.1/11"}
	Convey("Given a network with valid subnet", t, func() {
		Convey("When I try to validate this network", func() {
			valid, err := n.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, true)
				So(err, ShouldBeNil)
			})
		})
	})
	Convey("Given a network with an invalid subnet", t, func() {
		n.Subnet = "10.11.1.11"
		Convey("When I try to validate this network", func() {
			valid, err := n.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
			})
		})
	})
	Convey("Given a network with an empty subnet", t, func() {
		n.Subnet = ""
		Convey("When I try to validate this network", func() {
			valid, err := n.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestIsValidName(t *testing.T) {
	n := Network{Name: "foo", Router: "10.11.1.1", Subnet: "10.11.1.1/11"}
	Convey("Given a network with an invalid name", t, func() {
		n.Name = ""
		Convey("When I try to validate this network", func() {
			valid, err := n.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
			})
		})
	})
	Convey("Given a network with a name > 50 chars", t, func() {
		n.Name = "aksjhdlkashdliuhliusncldiudnalsundlaiunsdliausndliuansdlksbdlas"
		Convey("When I try to validate this network", func() {
			valid, err := n.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
