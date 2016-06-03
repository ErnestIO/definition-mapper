/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"io/ioutil"
	"net"
	"strings"
	"testing"

	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"

	. "github.com/smartystreets/goconvey/convey"
)

func getInputPayload(name string) (payload input.Payload) {
	body, err := ioutil.ReadFile("fixtures/mapping/" + name + ".in")
	if err != nil {
		panic(err)
	}
	input := []byte(body)
	if err := json.Unmarshal(input, &payload); err != nil {
		panic(err)
	}

	return payload
}

// This tests are useful if you're trying to accomplish a refactor, and
// you want to test as a black box the input & output of this service.
// They are skipped by default, so you should change the SkipConve by a
// single Convey if you want to use them. Happy testing :-)
func TestInput2OutputMapping(t *testing.T) {
	tests := []string{"basic", "instances_update", "instances_add"}
	SkipConvey("Given I'm mapping different entire definition", t, func() {
		for _, test := range tests {
			Convey("When I map a "+test+" definition", func() {
				payload := getInputPayload(test)
				m, err := BuildFSMMessage(payload, nil)
				So(err, ShouldBeNil)
				msg, err := json.Marshal(m)
				output := string(msg)
				So(err, ShouldBeNil)
				msg, err = ioutil.ReadFile("fixtures/mapping/" + test + ".out")
				expected := strings.TrimSpace(string(msg))
				So(err, ShouldBeNil)
				Convey("Then output should be expected", func() {
					So(output, ShouldResemble, expected)
				})
			})
		}
	})
}

func getBasicFSMMessage() *output.FSMMessage {
	msg := output.FSMMessage{
		ID:            "previous",
		ServiceName:   "foo",
		Endpoint:      "1.1.1.1",
		Bootstrapping: "salt",
		Service:       "previous",
	}
	msg.Routers.Items = append(msg.Routers.Items, output.Router{
		Name: "default",
	})
	msg.Networks.Items = append(msg.Networks.Items, output.Network{
		Name:   "foo-foo-default",
		Subnet: "10.1.0.0/24",
	})
	msg.Networks.Items = append(msg.Networks.Items, output.Network{
		Name:   "foo-foo-salt",
		Subnet: "10.2.0.0/24",
	})
	disks := []output.InstanceDisk{}
	disks = append(disks, output.InstanceDisk{ID: 1, Size: 10})
	msg.Instances.Items = append(msg.Instances.Items, output.Instance{
		Name:        "new",
		Catalog:     "cat",
		Image:       "img",
		Cpus:        1,
		Memory:      1024,
		Disks:       disks,
		NetworkName: "default",
		IP:          net.IP("1.1.1.1"),
	})
	return &msg
}

func getBasicPayload() *input.Payload {
	datacenter := input.Datacenter{Name: "foo"}

	networks := make([]input.Network, 0)
	networks = append(networks, input.Network{
		Name:   "default",
		Router: "default",
		Subnet: "10.1.0.0/24",
	})

	instances := make([]input.Instance, 0)
	instances = append(instances, input.Instance{
		Count:  1,
		Cpus:   1,
		Image:  "cat/img",
		Memory: "8GB",
		Disks:  []string{"10GB", "10GB"},
		Name:   "new",
		Networks: input.InstanceNetworks{
			Name:    "default",
			StartIP: net.ParseIP("10.1.0.11"),
		},
	})

	iRouters := make([]input.Router, 0)
	iRouters = append(iRouters, input.Router{
		Name:     "router",
		Networks: networks,
	})

	service := input.ServiceDefinition{
		Name:          "foo",
		Datacenter:    "foo",
		Instances:     instances,
		Routers:       iRouters,
		Bootstrapping: "salt",
		ErnestIP:      []string{"10.1.0.11"},
	}

	payload := input.Payload{
		ServiceID:  "new",
		PrevID:     "previous",
		Service:    service,
		Datacenter: datacenter,
	}

	return &payload
}

func TestAddInstancesMapping(t *testing.T) {
	Convey("Given I have a previously existing service", t, func() {
		prev := getBasicFSMMessage()
		Convey("When I submit a new definition with instances increase", func() {
			payload := *getBasicPayload()
			m, err := BuildFSMMessage(payload, prev)
			Convey("Then I should get a valid output", func() {
				So(err, ShouldEqual, nil)
				So(m.ID, ShouldEqual, payload.ServiceID)
				So(len(m.NetworksToCreate.Items), ShouldEqual, 0)
				So(len(m.Workflow.Arcs), ShouldEqual, 14)
				So(len(m.InstancesToCreate.Items), ShouldEqual, 2)
				So(len(m.InstancesToCreate.Items[1].Disks), ShouldEqual, 2)
				So(m.Instances.Items[1].Disks[0].Size, ShouldEqual, 10240)
			})
		})
	})
}
