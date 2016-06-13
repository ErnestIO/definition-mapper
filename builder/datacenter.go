/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"github.com/ernestio/definition-mapper/input"
	"github.com/ernestio/definition-mapper/output"
)

// MapDatacenter : Maps input datacenters to ernest understandable datacanters
func MapDatacenters(payload input.Payload) (datacenters []output.Datacenter, err error) {
	d := output.Datacenter{}
	if valid, err := payload.Datacenter.IsValid(); valid == false {
		return datacenters, err
	}
	d.Name = payload.Datacenter.Name
	d.Username = payload.Datacenter.Username
	d.Password = payload.Datacenter.Password
	d.Region = payload.Datacenter.Region
	d.Type = payload.Datacenter.Type
	d.ExternalNetwork = payload.Datacenter.ExternalNetwork
	d.VCloudURL = payload.Datacenter.VCloudURL
	d.VseURL = payload.Datacenter.VseURL

	return append(datacenters, d), err
}
