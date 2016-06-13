/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"testing"

	"github.com/ernestio/definition-mapper/input"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRouterMapping(t *testing.T) {
	Convey("Given a valid name and datacenter type", t, func() {
		routers := []input.Router{}
		routers = append(routers, input.Router{
			Name: "foo",
		})
		p := input.Payload{
			Service: input.ServiceDefinition{
				Name:    "foo",
				Routers: routers,
			},
			Datacenter: input.Datacenter{Name: "bar", Type: "bar"},
		}
		r, e := MapRouters(p, nil, "")
		Convey("Should successfully map an output router", func() {
			So(r[0].Name, ShouldEqual, "foo")
			So(r[0].Type, ShouldEqual, p.Datacenter.Type)
			So(e, ShouldBeNil)
		})
	})
}
