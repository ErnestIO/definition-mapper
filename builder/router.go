/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"github.com/r3labs/definition-mapper/input"
	"github.com/r3labs/definition-mapper/output"
)

// MapRouters : Maps input router to an ernest formatted router
func MapRouters(payload input.Payload, prev *output.FSMMessage, serviceIP string) (routers []output.Router, err error) {
	if prev == nil {
		for _, v := range payload.Service.Routers {
			r := output.Router{}
			r.Name = v.Name
			r.Type = payload.Datacenter.Type
			r.IP = serviceIP
			routers = append(routers, r)
		}
	} else {
		routers = prev.Routers.Items
	}
	return routers, nil
}
