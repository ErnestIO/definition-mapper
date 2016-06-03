/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

func TestInstancesMapping(t *testing.T) {
	Convey("Given I'm mapping the input instances", t, func() {
		datacenter := input.Datacenter{Name: "datacenter"}
		networks := []input.Network{}
		networks = append(networks, input.Network{Name: "datacenter-service-bar"})
		disks := []string{}
		disks = append(disks, "10GB")
		instances := []input.Instance{}
		instances = append(instances, input.Instance{
			Name:   "foo",
			Image:  "catalog/image",
			Cpus:   1,
			Memory: "2GB",
			Count:  1,
			Networks: input.InstanceNetworks{
				Name:    "bar",
				StartIP: net.ParseIP("10.1.0.1"),
			},
			Disks: disks,
		})
		iRouters := []input.Router{}
		iRouters = append(iRouters, input.Router{
			Name:     "router",
			Networks: networks,
		})
		p := input.Payload{
			Service: input.ServiceDefinition{
				Bootstrapping: "salt",
				Instances:     instances,
				Routers:       iRouters,
			},
			Datacenter: datacenter,
		}
		oNetworks := []output.Network{}
		oNetworks = append(oNetworks, output.Network{
			Name:   "datacenter-service-bar",
			Subnet: "10.1.0.1/24",
		})
		m := output.FSMMessage{
			ServiceName: "service",
		}
		m.Networks.Items = oNetworks

		Convey("When my service is not valid", func() {

			Convey("And an instance name is empty", func() {
				p.Service.Instances[0].Name = ""
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance name should not be null")
			})

			Convey("And an instance name is > 50", func() {
				p.Service.Instances[0].Name = "very_large_name_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance name can't be greater than 50 characters")
			})

			Convey("And an instance image is empty", func() {
				p.Service.Instances[0].Image = ""
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image should not be null")
			})

			Convey("And an instance image has an invalid format", func() {
				p.Service.Instances[0].Image = "bar"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image invalid, use format <catalog>/<image>")
			})

			Convey("And an instance image-catalog is empty", func() {
				p.Service.Instances[0].Image = "/bar"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image catalog should not be null, use format <catalog>/<image>")
			})

			Convey("And an instance image-image is empty", func() {
				p.Service.Instances[0].Image = "foo/"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image image should not be null, use format <catalog>/<image>")
			})

			Convey("And an instance cpus are < 1", func() {
				p.Service.Instances[0].Cpus = 0
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance cpus should not be < 1")
			})

			Convey("And an instance memory is empty", func() {
				p.Service.Instances[0].Memory = ""
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance memory should not be null")
			})

			Convey("And an instance count is < 1", func() {
				p.Service.Instances[0].Count = 0
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance count should not be < 1")
			})

			Convey("And an instance network name is null", func() {
				p.Service.Instances[0].Networks.Name = ""
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance network name should not be null")
			})

			Convey("And an instance network start_ip is null", func() {
				p.Service.Instances[0].Networks.StartIP = net.ParseIP("")
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance network start_ip should not be null")
			})

			Convey("And instances don't specify a valid network", func() {
				p.Service.Instances[0].Networks.Name = "wrong"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance must specify a valid network defined on the spec!")
			})

			// TODO: This probably should be validated on networks not here
			Convey("And network does not have a valid subnet", func() {
				m.Networks.Items[0].Subnet = "wrong"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Could not parse Network Subnet")
			})

			Convey("And instance ip is not valid inside network range", func() {
				m.Networks.Items[0].Subnet = "12.0.0.17/10"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Start IP invalid. IP must be a valid IP in the same range as it's network")
			})

			Convey("And memory format is not valid", func() {
				p.Service.Instances[0].Memory = "2Gb"
				_, err := MapInstances(p, m)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Invalid memory format")
			})
		})

		Convey("When my service is valid", func() {
			Convey("And is salt bootstrapped", func() {
				i, err := MapInstances(p, m)

				Convey("Then an extra instance should be mapped", func() {
					So(err, ShouldBeNil)
					So(len(i), ShouldEqual, len(p.Service.Instances)+1)
					So(i[0].Name, ShouldEqual, "datacenter-service-salt-master")
					So(i[0].Catalog, ShouldEqual, "r3")
					So(i[0].Image, ShouldEqual, "r3-salt-master")
					So(i[0].Cpus, ShouldEqual, 1)
					So(i[0].Memory, ShouldEqual, 2048)
					So(i[0].Disks, ShouldResemble, []output.InstanceDisk{})
					So(i[0].NetworkName, ShouldEqual, "datacenter-service-salt")
					So(i[0].IP, ShouldResemble, net.ParseIP("10.254.254.100"))
				})

				Convey("Then an defined instances should be mapped", func() {
					So(i[1].Name, ShouldEqual, "datacenter-service-foo-1")
					So(i[1].Catalog, ShouldEqual, "catalog")
					So(i[1].Image, ShouldEqual, "image")
					So(i[1].Cpus, ShouldEqual, 1)
					So(i[1].Memory, ShouldEqual, 2048)
					So(i[1].Disks[0].Size, ShouldEqual, 10240)
					So(i[1].NetworkName, ShouldEqual, "datacenter-service-bar")
					So(i[1].IP, ShouldResemble, net.ParseIP("10.1.0.1"))
				})
			})

			Convey("And is not salt bootstrapped", func() {
				p.Service.Bootstrapping = "foo"
				i, err := MapInstances(p, m)
				Convey("Then defined instances should be mapped", func() {
					So(err, ShouldBeNil)
					So(len(i), ShouldEqual, len(p.Service.Instances))
					So(i[0].Name, ShouldEqual, "datacenter-service-foo-1")
					So(i[0].Catalog, ShouldEqual, "catalog")
					So(i[0].Image, ShouldEqual, "image")
					So(i[0].Cpus, ShouldEqual, 1)
					So(i[0].Memory, ShouldEqual, 2048)
					So(i[0].Disks[0].Size, ShouldEqual, 10240)
					So(i[0].NetworkName, ShouldEqual, "datacenter-service-bar")
					So(i[0].IP, ShouldResemble, net.ParseIP("10.1.0.1"))
				})
			})
		})
	})
}

func TestInstancesToUpdateMapping(t *testing.T) {
	Convey("Given I'm mapping the input instances to be updated", t, func() {
		inInstances := []output.Instance{}
		inInstances = append(inInstances, output.Instance{
			Name:   "foo",
			Cpus:   1,
			Memory: 2048,
			Disks:  []output.InstanceDisk{},
		})
		o := output.FSMMessage{}
		o.Instances.Items = inInstances
		outInstances := []output.Instance{}
		outInstances = append(outInstances, output.Instance{
			Name:   "foo",
			Cpus:   1,
			Memory: 2048,
			Disks:  []output.InstanceDisk{},
		})
		n := output.FSMMessage{}
		n.Instances.Items = outInstances

		Convey("When there are no updates", func() {
			r, err := MapInstancesToUpdate(&o, n)

			Convey("Then I should get an empty list of instances", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 0)
			})
		})

		Convey("When there are updates on non allowed fields", func() {
			o.Instances.Items[0].Catalog = "something"
			r, err := MapInstancesToUpdate(&o, n)

			Convey("Then I should get an empty list of instances", func() {
				So(err, ShouldBeNil)
				So(len(r), ShouldEqual, 0)
			})
		})

		Convey("When there are updates on allowed fields", func() {

			Convey("And cpu resource has changed", func() {
				n.Instances.Items[0].Cpus = 2
				r, err := MapInstancesToUpdate(&o, n)
				Convey("Then I should get a list of instances to update with updated cpu", func() {
					So(err, ShouldBeNil)
					So(len(r), ShouldEqual, 1)
					So(r[0].Cpus, ShouldEqual, 2)
				})
			})

			Convey("And memory resource has changed", func() {
				n.Instances.Items[0].Memory = 1024
				r, err := MapInstancesToUpdate(&o, n)
				Convey("Then I should get a list of instances to update with updated memory", func() {
					So(err, ShouldBeNil)
					So(len(r), ShouldEqual, 1)
					So(r[0].Memory, ShouldEqual, 1024)
				})
			})

			Convey("And disks resource has changed", func() {
				n.Instances.Items[0].Disks = append(n.Instances.Items[0].Disks, output.InstanceDisk{
					Size: 100,
				})
				r, err := MapInstancesToUpdate(&o, n)
				Convey("Then I should get a list of instances to update with updated disks", func() {
					So(err, ShouldBeNil)
					So(len(r), ShouldEqual, 1)
					So(r[0].Disks[0].Size, ShouldEqual, 100)
				})
			})
		})
	})
}

func TestInstancesToDeleteMapping(t *testing.T) {
	Convey("Given I'm mapping the input instances to be deleted", t, func() {
		outputInstances := []output.Instance{}
		inputInstances := []output.Instance{}
		outputInstances = append(outputInstances, output.Instance{Name: "supu"})
		outputInstances = append(outputInstances, output.Instance{Name: "supu"})
		inputInstances = append(inputInstances, output.Instance{Name: "supu"})
		inputInstances = append(inputInstances, output.Instance{Name: "supu"})
		Convey("When I have an equal input and output", func() {
			o := output.FSMMessage{}
			o.Instances.Items = inputInstances
			n := output.FSMMessage{}
			n.Instances.Items = outputInstances
			result, err := MapInstancesToDelete(&o, n)
			Convey("Then I should receive no instances to delete", func() {
				So(len(result), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})

			Convey("When the input provides less instances than the existing ones", func() {
				inputInstances = append(inputInstances, output.Instance{Name: "supu"})
				o := output.FSMMessage{}
				o.Instances.Items = inputInstances
				n := output.FSMMessage{}
				n.Instances.Items = outputInstances
				result, err := MapInstancesToDelete(&o, n)
				Convey("Then I should receive one instance to be deleted", func() {
					So(len(result), ShouldEqual, 1)
					So(err, ShouldBeNil)
				})
			})

			Convey("When the input provides more instances than the existing ones", func() {
				outputInstances = append(outputInstances, output.Instance{Name: "supu"})
				o := output.FSMMessage{}
				o.Instances.Items = inputInstances
				n := output.FSMMessage{}
				n.Instances.Items = outputInstances
				Convey("When get the instances to be deleted", func() {
					result, err := MapInstancesToDelete(&o, n)
					Convey("Then I should receive no instances to delete", func() {
						So(len(result), ShouldEqual, 0)
						So(err, ShouldBeNil)
					})
				})
			})
		})

	})
}

func TestInstancesToCreateMapping(t *testing.T) {
	Convey("Given I'm mapping the input instances to be created", t, func() {
		outputInstances := []output.Instance{}
		inputInstances := []output.Instance{}
		outputInstances = append(outputInstances, output.Instance{Name: "supu"})
		outputInstances = append(outputInstances, output.Instance{Name: "supu"})
		inputInstances = append(inputInstances, output.Instance{Name: "supu"})
		inputInstances = append(inputInstances, output.Instance{Name: "supu"})
		Convey("When I have an equal input and output", func() {
			o := output.FSMMessage{}
			o.Instances.Items = inputInstances
			n := output.FSMMessage{}
			n.Instances.Items = outputInstances
			result, err := MapInstancesToCreate(&o, n)
			Convey("Then I should receive no instances to be created", func() {
				So(len(result), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
		})

		Convey("When input is providing more instances than existing ones", func() {
			outputInstances = append(outputInstances, output.Instance{Name: "supu"})
			o := output.FSMMessage{}
			o.Instances.Items = inputInstances
			n := output.FSMMessage{}
			n.Instances.Items = outputInstances
			result, err := MapInstancesToCreate(&o, n)
			Convey("Then I should receive one instance to be created", func() {
				So(len(result), ShouldEqual, 1)
				So(err, ShouldBeNil)
			})
		})

		Convey("When input provides less instances than existing ones", func() {
			inputInstances = append(inputInstances, output.Instance{Name: "supu"})
			o := output.FSMMessage{}
			o.Instances.Items = inputInstances
			n := output.FSMMessage{}
			n.Instances.Items = outputInstances
			result, err := MapInstancesToCreate(&o, n)
			Convey("Then I should receive one instance", func() {
				So(len(result), ShouldEqual, 0)
				So(err, ShouldBeNil)
			})
		})
	})
}
