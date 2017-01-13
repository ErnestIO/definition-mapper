/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"os"
	"runtime"
	"time"

	ecc "github.com/ernestio/ernest-config-client"
	"github.com/nats-io/nats"
)

var n *nats.Conn
var err error

func getType(body []byte) string {
	var service struct {
		Datacenter struct {
			Type string `json:"type"`
		} `json:"datacenter"`
	}

	if err := json.Unmarshal(body, &service); err != nil {
		log.Panic(err)
	}

	return service.Datacenter.Type
}

func route(msg *nats.Msg, action string) []byte {
	mapper := ""
	switch getType(msg.Data) {
	case "vcloud", "fake", "vcloud-fake":
		mapper = "vcloud"
	case "aws", "aws-fake":
		mapper = "aws"
	default:
		return []byte(`{"error":"Invalid type"}`)
	}

	subject := "definition.map." + action + "." + mapper
	msg, err := n.Request(subject, msg.Data, time.Second)
	if err != nil {
		log.Println("Error processing " + subject)
		log.Println(err.Error())
	}

	return msg.Data
}

// SubscribeCreateService : definition.map.creation subscriber
func SubscribeCreateService(msg *nats.Msg) {
	res := route(msg, "creation")
	if err := n.Publish(msg.Reply, res); err != nil {
		log.Panic(err)
	}
}

// SubscribeImportService : definition.map.import subscriber
func SubscribeImportService(msg *nats.Msg) {
	res := route(msg, "import")
	if err := n.Publish(msg.Reply, res); err != nil {
		log.Panic(err)
	}
}

// SubscribeDeleteService : definition.map.deletion subscriber
func SubscribeDeleteService(msg *nats.Msg) {
	res := route(msg, "deletion")
	if err := n.Publish(msg.Reply, res); err != nil {
		log.Panic(err)
	}
}

// Subscribe : Manages all subscriptions
func Subscribe() {
	if _, err := n.Subscribe("definition.map.creation", SubscribeCreateService); err != nil {
		log.Panic(err)
	}
	if _, err := n.Subscribe("definition.map.import", SubscribeImportService); err != nil {
		log.Panic(err)
	}
	if _, err := n.Subscribe("definition.map.deletion", SubscribeDeleteService); err != nil {
		log.Panic(err)
	}
}

func setup() {
	n = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func main() {
	setup()
	Subscribe()
	runtime.Goexit()
}
