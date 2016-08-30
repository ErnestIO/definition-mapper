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
	json.Unmarshal(body, &service)

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

func SubscribeCreateService(msg *nats.Msg) {
	res := route(msg, "creation")
	n.Publish(msg.Reply, res)
}

func SubscribeDeleteService(msg *nats.Msg) {
	res := route(msg, "deletion")
	n.Publish(msg.Reply, res)
}

func Subscribe() {
	n.Subscribe("definition.map.creation", SubscribeCreateService)
	n.Subscribe("definition.map.deletion", SubscribeDeleteService)
}

func setup() {
	n = ecc.NewConfig(os.Getenv("NATS_URI")).Nats()
}

func main() {
	setup()
	Subscribe()
	runtime.Goexit()
}
