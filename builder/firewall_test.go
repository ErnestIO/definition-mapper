/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ernestio/definition-mapper/output"
)

func TestHasChangedFirewall(t *testing.T) {
	Convey("Given the firewall rules has not changed", t, func() {
		n := []output.FirewallRule{}
		n = append(n, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		o := []output.FirewallRule{}
		o = append(o, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		Convey("When I check if rules has been changed", func() {
			fo := []output.Firewall{}
			fo = append(fo, output.Firewall{Rules: o})
			fn := []output.Firewall{}
			fn = append(fn, output.Firewall{Rules: n})
			result := HasChangedFirewalls(fo, fn)
			Convey("Then I should receive false", func() {
				So(result, ShouldBeFalse)
			})
		})
	})

	Convey("Given the firewall rules has been changed", t, func() {
		n := []output.FirewallRule{}
		n = append(n, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		o := []output.FirewallRule{}
		o = append(o, output.FirewallRule{DestinationIP: "london", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		Convey("When I check if rules has been changed", func() {
			result := HasChangedFirewallRules(o, n)
			Convey("Then I should receive true", func() {
				So(result, ShouldBeTrue)
			})
		})
	})

	Convey("Given we removed existing rules", t, func() {
		n := []output.FirewallRule{}
		n = append(n, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		n = append(n, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		o := []output.FirewallRule{}
		o = append(o, output.FirewallRule{DestinationIP: "rome", DestinationPort: "222", Protocol: "A", SourceIP: "A", SourcePort: "222"})
		Convey("When I check if rules has been changed", func() {
			result := HasChangedFirewallRules(o, n)
			Convey("Then I should receive true", func() {
				So(result, ShouldBeTrue)
			})
		})
	})

}
