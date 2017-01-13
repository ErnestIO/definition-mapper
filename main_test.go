/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"errors"
	"testing"
	"time"

	"github.com/nats-io/nats"
	. "github.com/smartystreets/goconvey/convey"
)

func waitMsg(ch chan *nats.Msg) (*nats.Msg, error) {
	select {
	case msg := <-ch:
		return msg, nil
	case <-time.After(time.Millisecond * 5000):
	}
	return nil, errors.New("timeout")
}
func dummySubscribe(subject string) {
	_, _ = n.Subscribe(subject, func(m *nats.Msg) {
		_ = n.Publish(m.Reply, []byte(subject))
	})
}

func setupTest() {
	dummySubscribe("definition.map.creation.vcloud")
	dummySubscribe("definition.map.creation.aws")
	dummySubscribe("definition.map.deletion.vcloud")
	dummySubscribe("definition.map.deletion.aws")
}

func TestVSE(t *testing.T) {
	setup()
	Subscribe()
	setupTest()

	Convey("Given I have a running definition mapper", t, func() {
		Convey("When I send a service without datacenter type", func() {
			body := []byte(`{"foo":"bar"}`)
			res, err := n.Request("definition.map.creation", body, time.Second)
			So(err, ShouldBeNil)
			So(string(res.Data), ShouldEqual, `{"error":"Invalid type"}`)
		})

		Convey("When I send a service to create with a vcloud datacenter type", func() {
			body := []byte(`{"datacenter":{"type":"vcloud"},"foo":"bar"}`)
			res, err := n.Request("definition.map.creation", body, time.Second)
			Convey("Then I should successfully create a valid service", func() {
				So(err, ShouldBeNil)
				So(string(res.Data), ShouldEqual, "definition.map.creation.vcloud")
			})
		})

		Convey("When I send a service to delete with a vcloud datacenter type", func() {
			body := []byte(`{"datacenter":{"type":"vcloud"},"foo":"bar"}`)
			res, err := n.Request("definition.map.deletion", body, time.Second)
			Convey("Then I should successfully delete a valid service", func() {
				So(err, ShouldBeNil)
				So(string(res.Data), ShouldEqual, "definition.map.deletion.vcloud")
			})
		})

		Convey("When I send a service to create with a aws datacenter type", func() {
			body := []byte(`{"datacenter":{"type":"aws"},"foo":"bar"}`)
			res, err := n.Request("definition.map.creation", body, time.Second)
			Convey("Then I should successfully create a valid service", func() {
				So(err, ShouldBeNil)
				So(string(res.Data), ShouldEqual, "definition.map.creation.aws")
			})
		})

		Convey("When I send a service to delete with a aws datacenter type", func() {
			body := []byte(`{"datacenter":{"type":"aws"},"foo":"bar"}`)
			res, err := n.Request("definition.map.deletion", body, time.Second)
			Convey("Then I should successfully delete a valid service", func() {
				So(err, ShouldBeNil)
				So(string(res.Data), ShouldEqual, "definition.map.deletion.aws")
			})
		})

	})
}
