/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package input

import (
	"net"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsValidInstanceName(t *testing.T) {
	i := Instance{Name: "foo"}
	Convey("Given an instance with an invalid name", t, func() {
		i.Name = ""
		Convey("When I try to validate this isntance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance name should not be null")
			})
		})
	})
	Convey("Given an instance with a name > 50 chars", t, func() {
		i.Name = "aksjhdlkashdliuhliusncldiudnalsundlaiunsdliausndliuansdlksbdlas"
		Convey("When I try to validate this intance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance name can't be greater than 50 characters")
			})
		})
	})
}

func TestIsValidInstanceImage(t *testing.T) {
	i := Instance{Name: "foo", Image: "image"}
	Convey("Given an instance with an invalid name", t, func() {
		i.Image = ""
		Convey("When I try to validate this isntance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image should not be null")
			})
		})
	})
	Convey("Given an instance with an invalid image format", t, func() {
		i.Image = "aksjhdlkashdliuhliusncldiud"
		Convey("When I try to validate this intance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image invalid, use format <catalog>/<image>")
			})
		})
	})
	Convey("Given an instance with an empty image catalog", t, func() {
		i.Image = "/image"
		Convey("When I try to validate this intance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image catalog should not be null, use format <catalog>/<image>")
			})
		})
	})
	Convey("Given an instance with an empty image", t, func() {
		i.Image = "catalog/"
		Convey("When I try to validate this intance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance image image should not be null, use format <catalog>/<image>")
			})
		})
	})
}

func TestIsValidInstanceCpu(t *testing.T) {
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2}
	Convey("Given an instance with a cpu field less than one", t, func() {
		i.Cpus = 0
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance cpus should not be < 1")
			})
		})
	})
}

func TestIsValidInstanceMemory(t *testing.T) {
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2, Memory: "2GB"}
	Convey("Given an instance with an empty memory field", t, func() {
		i.Memory = ""
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance memory should not be null")
			})
		})
	})
}

func TestIsValidInstanceCount(t *testing.T) {
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2, Memory: "2GB", Count: 1}
	Convey("Given an instance with a cpu field less than one", t, func() {
		i.Count = 0
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance count should not be < 1")
			})
		})
	})
}

func TestIsValidInstanceNetworksName(t *testing.T) {
	n := InstanceNetworks{Name: "foo", StartIP: net.ParseIP("127.0.0.1")}
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2, Memory: "2GB", Count: 1, Networks: n}
	Convey("Given an instance with an invalid networks name", t, func() {
		i.Networks.Name = ""
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance network name should not be null")
			})
		})
	})
}

func TestIsValidInstanceNetworkStartIP(t *testing.T) {
	n := InstanceNetworks{Name: "foo", StartIP: net.ParseIP("127.0.0.1")}
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2, Memory: "2GB", Count: 1, Networks: n}
	Convey("Given an instance with an invalid networks name", t, func() {
		i.Networks.StartIP = nil
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, false)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "Instance network start_ip should not be null")
			})
		})
	})
}

func TestIsValidInstanceHappyPath(t *testing.T) {
	n := InstanceNetworks{Name: "foo", StartIP: net.ParseIP("127.0.0.1")}
	i := Instance{Name: "foo", Image: "catalog/image", Cpus: 2, Memory: "2GB", Count: 1, Networks: n}
	Convey("Given a valid instance", t, func() {
		Convey("When I try to validate this instance", func() {
			valid, err := i.IsValid()
			Convey("Then should return true", func() {
				So(valid, ShouldEqual, true)
				So(err, ShouldBeNil)
			})
		})
	})
}
