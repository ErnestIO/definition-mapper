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

func TestBootstrapsMapping(t *testing.T) {
	Convey("Given I'm mapping the input bootstraps", t, func() {
		commands := []string{}
		commands = append(commands, "echo 'hola'")
		commands = append(commands, "echo 'hola'")
		execs := []input.Exec{}
		execs = append(execs, input.Exec{Commands: commands})
		instances := []input.Instance{}
		instances = append(instances, input.Instance{
			Name:        "foo",
			Provisioner: execs,
		})

		p := input.Payload{Service: input.ServiceDefinition{
			Bootstrapping: "salt",
			Instances:     instances,
		}}

		mInstances := []output.Instance{}
		mInstances = append(mInstances, output.Instance{Name: "foo"})
		m := output.FSMMessage{}
		m.Instances.Items = mInstances

		Convey("When I provide a provisioner with some commands", func() {
			r := GenerateBootstraps(p, nil, m)
			Convey("Then I should get a controlled error", func() {
				So(len(r), ShouldEqual, 1)
				So(r[0].Name, ShouldEqual, "Bootstrap foo")
				So(r[0].Type, ShouldEqual, "salt")
				So(r[0].Target, ShouldEqual, "list:salt-master.localdomain")
				So(r[0].Payload, ShouldEqual, "/usr/bin/bootstrap -master 10.254.254.100 -host <nil> -username ernest -password 'b00tStr4pp3rR' -max-retries 20 -minion-name foo")

			})
		})
	})
}
