/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package builder

import (
	"errors"
	"net"

	"github.com/r3labs/definition-mapper/input"
)

// MapServiceIP : validates and returns the ernest ip
func MapServiceIP(payload input.Payload) (string, error) {
	ip := payload.Service.ServiceIP
	if ip == "" {
		return ip, nil
	}
	if ok := net.ParseIP(ip); ok == nil {
		return "", errors.New("ServiceIP is not a valid IP")
	}

	return ip, nil
}
