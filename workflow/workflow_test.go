/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package workflow

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsValidWorkflow(t *testing.T) {
	Convey("Given a new workflow", t, func() {
		w := New("workflows/service_create.json")
		Convey("When I inspect the first arc", func() {
			firstArc := w.Arcs[0]
			Convey("Then it should return the correct values", func() {
				So(firstArc.From, ShouldEqual, "created")
				So(firstArc.To, ShouldEqual, "started")
				So(firstArc.Event, ShouldEqual, "service.create")
			})
		})
		Convey("When I add a new valid set of workflow arcs", func() {
			arc, err := LoadArcs("workflows/nats_update.json")
			Convey("There should be no error loading the arc", func() {
				So(err, ShouldBeNil)
			})
			w.Add(arc)
			secondArc := w.Arcs[1]
			Convey("Then the second arc should have been updated", func() {
				So(secondArc.To, ShouldEqual, "updating_nats")
				So(secondArc.Event, ShouldEqual, "nats.update")
			})
			thirdArc := w.Arcs[2]
			Convey("Then the third arc should have been added", func() {
				So(thirdArc.From, ShouldEqual, "updating_nats")
				So(thirdArc.To, ShouldEqual, "nats_updated")
				So(thirdArc.Event, ShouldEqual, "nats.update.done")
			})
		})
		Convey("When I finish the workflow", func() {
			arc, err := LoadArcs("workflows/nats_update.json")
			w.Add(arc)
			w.Finish("workflows/service_create_done.json")
			Convey("There should be no errors", func() {
				So(err, ShouldBeNil)
			})
			thirdArc := w.Arcs[len(w.Arcs)-3]
			Convey("Then the third from last arc should be loaded", func() {
				So(thirdArc.From, ShouldEqual, "nats_updated")
				So(thirdArc.To, ShouldEqual, "done")
				So(thirdArc.Event, ShouldEqual, "service.create.done")
			})
			secondArc := w.Arcs[len(w.Arcs)-2]
			Convey("Then the second from last arc should be loaded", func() {
				So(secondArc.From, ShouldEqual, "pre-failed")
				So(secondArc.To, ShouldEqual, "failed")
				So(secondArc.Event, ShouldEqual, "to_error")
			})
			finalArc := w.Arcs[len(w.Arcs)-1]
			Convey("Then the final arc should be loaded", func() {
				So(finalArc.From, ShouldEqual, "failed")
				So(finalArc.To, ShouldEqual, "errored")
				So(finalArc.Event, ShouldEqual, "service.create.error")
			})
			Convey("None of the values should be nil", func() {
				for _, arc := range w.Arcs {
					So(arc.From, ShouldNotBeNil)
					So(arc.To, ShouldNotBeNil)
					So(arc.Event, ShouldNotBeNil)
				}
			})
		})
	})
}

func TestIsInvalidWorkflow(t *testing.T) {
	Convey("Given an invalid workflow", t, func() {
		w := New("workflows/service_create.json")
		Convey("When I load a non-existent arc file", func() {
			_, err := LoadArcs("doesnt-exist")
			Convey("Then I should receive an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When I add an empty Arc", func() {
			arc, _ := LoadArcs("fixtures/invalid_workflow_empty.json")
			w.Add(arc)
			err := w.Finish("workflows/service_create_done.json")
			Convey("Then finishing the workflow should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
		Convey("When I add an invalid sequence Arc", func() {
			arc, _ := LoadArcs("fixtures/invalid_workflow_sequence.json")
			w.Add(arc)
			err := w.Finish("workflows/service_create_done.json")
			Convey("Then finishing the workflow should return an error", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
