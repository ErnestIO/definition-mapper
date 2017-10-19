/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ernestio/definition-mapper/build"
	"github.com/ernestio/definition-mapper/handlers"
	"github.com/ernestio/definition-mapper/request"
	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
	"github.com/r3labs/akira"
	"github.com/r3labs/graph"
)

var n akira.Connector

// StartMappingHandlers : start the primary mapping handlers
func StartMappingHandlers() {
	_, _ = n.Subscribe("mapping.get.*", func(msg *nats.Msg) {
		var r request.Request
		var g *graph.Graph
		var data []byte
		var err error

		defer response(msg.Reply, &data, &err)

		err = json.Unmarshal(msg.Data, &r)
		if err != nil {
			return
		}

		parts := strings.Split(msg.Subject, ".")
		switch parts[2] {
		case "create":
			g, err = handlers.Create(&r)
		case "update":
			g, err = handlers.Update(&r)
		case "delete":
			g, err = handlers.Delete(&r)
		case "import":
			g, err = handlers.Import(&r)
		case "diff":
			g, err = handlers.Diff(&r)
		}

		if err != nil {
			return
		}

		data, err = g.ToJSON()
	})
}

// StartSecondaryHandlers : start secondary handlers
func StartSecondaryHandlers() {
	_, _ = n.Subscribe("build.import.done", func(msg *nats.Msg) {
		var ig map[string]interface{}
		var b *build.Build
		var data []byte
		var err error

		defer response(msg.Reply, &data, &err)

		err = json.Unmarshal(msg.Data, &ig)
		if err != nil {
			return
		}

		b, err = handlers.ImportComplete(ig)
		if err != nil {
			return
		}

		data, err = json.Marshal(b)
		if err != nil {
			return
		}

		_, err = n.Request("build.set.mapping", data, time.Second*5)
		if err != nil {
			return
		}

		_, err = n.Request("build.set.definition", data, time.Second*5)
		if err != nil {
			return
		}

		data = []byte(`{"status": "success"}`)
	})
}

func setup() {
	n = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func main() {
	setup()
	StartMappingHandlers()
	StartSecondaryHandlers()
	runtime.Goexit()
}
