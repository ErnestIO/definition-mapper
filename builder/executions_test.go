/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
)

func TestExecutionsMapping(t *testing.T) {
	Convey("Given I'm trying to map executions", t, func() {
		execs := []input.Exec{}
		execs = append(execs, input.Exec{Commands: []string{"foo", "bar"}})
		execs = append(execs, input.Exec{Commands: []string{"bar"}})
		instances := []input.Instance{}
		instances = append(instances, input.Instance{
			Name:        "foo",
			Count:       1,
			Provisioner: execs,
		})

		p := input.Payload{
			Service: input.ServiceDefinition{
				Name:          "service",
				Bootstrapping: "salt",
				Instances:     instances,
			},
			Datacenter: input.Datacenter{Name: "datacenter"},
		}

		m := output.FSMMessage{}

		Convey("When the input doesn't have defined executions", func() {
			p.Service.Instances[0].Provisioner = []input.Exec{}
			r, err := GenerateExecutions(p, nil, m)

			Convey("Then I should get a list of executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 0)
			})
		})

		Convey("When the input has defined executions", func() {
			r, err := GenerateExecutions(p, nil, m)

			Convey("Then I should get a list of executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 2)

				Convey("And should successfully map these executions", func() {
					So(r[0].Name, ShouldEqual, "Execution foo 1")
					So(r[0].Type, ShouldEqual, "salt")
					So(r[0].Target, ShouldEqual, "list:datacenter-service-foo-1")
					So(r[0].Payload, ShouldEqual, "foo; bar")
				})
			})
		})
	})
}

func TestBuildChangedExecutions(t *testing.T) {
	Convey("Given I want to know if executions has changed", t, func() {
		// Previous Instances, Executions
		pi := []output.Instance{{Name: "test"}}
		pe := []output.Execution{{Name: "Execution test 1", Type: "salt", Payload: "echo hello"}}

		// Current Instances
		ci := []input.Instance{}
		instance := input.Instance{Name: "test"}
		instance.Provisioner = []input.Exec{}
		instance.Provisioner = append(instance.Provisioner, input.Exec{Commands: []string{"echo hello"}})
		ci = append(ci, instance)

		sr := output.FSMMessage{}
		sr.Executions.Items = pe
		sr.Instances.Items = pi

		p := input.Payload{
			Service: input.ServiceDefinition{
				Name:          "service",
				Bootstrapping: "salt",
				Instances:     ci,
			},
			Datacenter: input.Datacenter{Name: "datacenter"},
		}

		m := output.FSMMessage{}

		Convey("When executions haven't changed", func() {
			r, err := GenerateExecutions(p, &sr, m)
			Convey("Then I should receive an empty list of executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 0)
			})
		})

		Convey("When I update a provisioners payload", func() {
			p.Service.Instances[0].Provisioner[0] = input.Exec{Commands: []string{"echo bye"}}
			r, err := GenerateExecutions(p, &sr, m)
			Convey("Then I should receive a list of changed executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
			})
		})

		Convey("When I update the instance count without changing the execution payload", func() {
			p.Service.Instances[0].Count = 2
			m.InstancesToCreate.Items = []output.Instance{}
			m.InstancesToCreate.Items = append(m.InstancesToCreate.Items, output.Instance{Name: "datacenter-service-test-2"})
			r, err := GenerateExecutions(p, &sr, m)
			Convey("Then I should receive a list of changed executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
				So(r[0].Target, ShouldEqual, "list:datacenter-service-test-2")
			})
		})

		Convey("When I update the instance count and change the execution payload", func() {
			p.Service.Instances[0].Provisioner[0] = input.Exec{Commands: []string{"echo bye"}}
			p.Service.Instances[0].Count = 2
			m.InstancesToCreate.Items = []output.Instance{}
			m.InstancesToCreate.Items = append(m.InstancesToCreate.Items, output.Instance{Name: "datacenter-service-test-2"})
			r, err := GenerateExecutions(p, &sr, m)
			Convey("Then I should receive a list of changed executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
				So(r[0].Target, ShouldEqual, "list:datacenter-service-test-1,datacenter-service-test-2")
			})
		})

		Convey("When I add a new instance type without changing the execution payload for an existing instance", func() {
			execs := []input.Exec{input.Exec{Commands: []string{"date"}}}
			instance := input.Instance{Name: "new", Provisioner: execs, Count: 1}
			p.Service.Instances = append(p.Service.Instances, instance)
			m.InstancesToCreate.Items = []output.Instance{}
			m.InstancesToCreate.Items = append(m.InstancesToCreate.Items, output.Instance{Name: "datacenter-service-new-1"})
			r, err := GenerateExecutions(p, &sr, m)
			Convey("Then I should receive a list of changed executions", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 1)
				So(r[0].Target, ShouldEqual, "list:datacenter-service-new-1")
			})
		})

	})

}
