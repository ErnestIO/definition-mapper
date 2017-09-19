/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
)

// Error : default error message
type Error struct {
	Error string `json:"error"`
}

func response(reply string, data *[]byte, err *error) {
	var rdata []byte
	if data != nil {
		rdata = *data
	}

	if *err != nil {
		rdata, _ = json.Marshal(Error{Error: (*err).Error()})
	}

	if reply != "" {
		n.Publish(reply, rdata)
	} else if *err != nil {
		log.Println(*err)
	}
}
