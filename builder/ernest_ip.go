/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"errors"

	"github.com/r3labs/definition-mapper/input"
)

// MapErnestIP : validates and returns the ernest ip
func MapErnestIP(payload input.Payload) ([]string, error) {
	if len(payload.Service.ErnestIP) == 0 && payload.Service.Bootstrapping == "salt" {
		return []string{}, errors.New("You've defined a salt bootstrapping method, but no ernest ips are mapped, please add them to ernest_ip field")
	}
	return payload.Service.ErnestIP, nil
}
