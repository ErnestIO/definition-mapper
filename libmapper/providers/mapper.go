/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package providers

import (
	"github.com/ernestio/definition-mapper/libmapper"
	aws "github.com/ernestio/definition-mapper/libmapper/providers/aws/mapper"
	azure "github.com/ernestio/definition-mapper/libmapper/providers/azure/mapper"
	vcloud "github.com/ernestio/definition-mapper/libmapper/providers/vcloud/mapper"
)

// NewMapper : Get a new mapper based on a specified type
func NewMapper(t string) (m libmapper.Mapper) {
	switch t {
	case "aws", "aws-fake":
		m = aws.New()
	case "vcloud", "vcloud-fake":
		m = vcloud.New()
	case "azure", "azure-fake":
		m = azure.New()
	}

	return m
}
