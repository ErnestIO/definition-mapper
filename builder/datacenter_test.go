/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"testing"

	"github.com/ernestio/definition-mapper/input"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDatacenterMapping(t *testing.T) {
	Convey("Given I try to map a datacenter", t, func() {
		p := input.Payload{
			Datacenter: input.Datacenter{
				Name:     "name",
				Username: "username",
				Password: "pwd",
				Region:   "region",
				Type:     "type",
			},
		}
		Convey("When the datacenter is not valid", func() {
			Convey("And the datacenter name is empty", func() {
				p.Datacenter.Name = ""
				_, err := MapDatacenters(p)
				Convey("Then we will receive an invalid name error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Datacenter name should not be null")
				})
			})
			Convey("And the datacenter name > 50", func() {
				p.Datacenter.Name = "namebiggerthan50characters_xxxxxxxxxxxxxxxxxxxxxxxx"
				_, err := MapDatacenters(p)
				Convey("Then we will receive a > 50 error", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, "Datacenter name can't be greater than 50 characters")
				})
			})
		})

		Convey("When the datacenter is valid", func() {
			d, err := MapDatacenters(p)
			Convey("Then we will receive a valid mapped datacenter", func() {
				So(err, ShouldBeNil)
				So(d[0].Name, ShouldEqual, p.Datacenter.Name)
				So(d[0].Username, ShouldEqual, p.Datacenter.Username)
				So(d[0].Password, ShouldEqual, p.Datacenter.Password)
				So(d[0].Region, ShouldEqual, p.Datacenter.Region)
				So(d[0].Type, ShouldEqual, p.Datacenter.Type)
			})
		})

	})
}
